-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestContainerIgnoredIndex default#1
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c TEXT, g GEOMETRY NOT NULL, KEY i (a) IGNORED, UNIQUE KEY u (b) COMMENT 'u note' IGNORED, KEY n (a, b) NOT IGNORED, KEY v (b) VISIBLE, KEY iv (a) IGNORED VISIBLE, KEY ni (b) NOT IGNORED IGNORED, FULLTEXT KEY f (c) IGNORED, SPATIAL KEY s (g) COMMENT 'sp' IGNORED, KEY my (a) /*!80000 INVISIBLE */)
< ok 0 0
> query
| SELECT DISTINCT INDEX_NAME, IGNORED FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 't'
< result
< column "INDEX_NAME" "VARCHAR"
< column "IGNORED" "VARCHAR"
< row
: bytes
| u
: bytes
| YES
< row
: bytes
| i
: bytes
| YES
< row
: bytes
| n
: bytes
| NO
< row
: bytes
| v
: bytes
| NO
< row
: bytes
| iv
: bytes
| YES
< row
: bytes
| ni
: bytes
| YES
< row
: bytes
| s
: bytes
| YES
< row
: bytes
| my
: bytes
| NO
< row
: bytes
| f
: bytes
| YES
> query
| SHOW CREATE TABLE t
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
|   `c` text DEFAULT NULL,
|   `g` geometry NOT NULL,
|   UNIQUE KEY `u` (`b`) COMMENT 'u note' IGNORED,
|   KEY `i` (`a`) IGNORED,
|   KEY `n` (`a`,`b`),
|   KEY `v` (`b`),
|   KEY `iv` (`a`) IGNORED,
|   KEY `ni` (`b`) IGNORED,
|   SPATIAL KEY `s` (`g`) COMMENT 'sp' IGNORED,
|   KEY `my` (`a`),
|   FULLTEXT KEY `f` (`c`) IGNORED
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE INDEX ci ON t (a) IGNORED
< ok 0 0
> exec
| CREATE UNIQUE INDEX cu ON t (b) COMMENT 'c' IGNORED
< ok 0 0
> exec
| ALTER TABLE t ADD INDEX ai (b) IGNORED, ADD KEY an (a) NOT IGNORED
< ok 0 0
> exec
| ALTER TABLE t ALTER INDEX ci NOT IGNORED
< ok 0 0
> exec
| ALTER TABLE t ALTER KEY an IGNORED
< ok 0 0
> exec
| ALTER TABLE t ALTER INDEX IF EXISTS nosuch IGNORED
< ok 0 0
> exec
| ALTER TABLE t RENAME INDEX ai TO ar
< ok 0 0
> query
| SELECT DISTINCT INDEX_NAME, IGNORED FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 't'
< result
< column "INDEX_NAME" "VARCHAR"
< column "IGNORED" "VARCHAR"
< row
: bytes
| cu
: bytes
| YES
< row
: bytes
| ci
: bytes
| NO
< row
: bytes
| ar
: bytes
| YES
< row
: bytes
| an
: bytes
| YES
> query
| SHOW CREATE TABLE t
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
|   UNIQUE KEY `cu` (`b`) COMMENT 'c' IGNORED,
|   KEY `ci` (`a`),
|   KEY `ar` (`b`) IGNORED,
|   KEY `an` (`a`) IGNORED
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, KEY k (a))
< ok 0 0
> exec
| ALTER TABLE t ADD INDEX n (b), ALTER INDEX IF EXISTS n IGNORED, ALTER INDEX k IGNORED
< ok 0 0
> exec
| ALTER TABLE t RENAME INDEX n TO n2, ALTER INDEX IF EXISTS n2 IGNORED
< ok 0 0
> query
| SELECT DISTINCT INDEX_NAME, IGNORED FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 't'
< result
< column "INDEX_NAME" "VARCHAR"
< column "IGNORED" "VARCHAR"
< row
: bytes
| k
: bytes
| YES
< row
: bytes
| n2
: bytes
| NO
> query
| SHOW CREATE TABLE t
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
|   KEY `k` (`a`) IGNORED,
|   KEY `n2` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY k (a))
< ok 0 0
> exec
| ALTER TABLE t ALTER INDEX k NOT IGNORED, ALTER INDEX k IGNORED
< ok 0 0
> query
| SELECT DISTINCT INDEX_NAME, IGNORED FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 't'
< result
< column "INDEX_NAME" "VARCHAR"
< column "IGNORED" "VARCHAR"
< row
: bytes
| k
: bytes
| YES
> query
| SHOW CREATE TABLE t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   KEY `k` (`a`) IGNORED
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE src (a INT, KEY i (a) COMMENT 'c' IGNORED)
< ok 0 0
> exec
| CREATE TABLE t LIKE src
< ok 0 0
> query
| SELECT DISTINCT INDEX_NAME, IGNORED FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 't'
< result
< column "INDEX_NAME" "VARCHAR"
< column "IGNORED" "VARCHAR"
< row
: bytes
| i
: bytes
| YES
> query
| SHOW CREATE TABLE t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   KEY `i` (`a`) COMMENT 'c' IGNORED
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT, PRIMARY KEY (a) NOT IGNORED VISIBLE, KEY k (b) IGNORED)
< ok 0 0
> query
| SELECT DISTINCT INDEX_NAME, IGNORED FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 't'
< result
< column "INDEX_NAME" "VARCHAR"
< column "IGNORED" "VARCHAR"
< row
: bytes
| PRIMARY
: bytes
| NO
< row
: bytes
| k
: bytes
| YES
> query
| SHOW CREATE TABLE t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`a`),
|   KEY `k` (`b`) IGNORED
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, c INT, d VARCHAR(10) NOT NULL, UNIQUE KEY v (b), UNIQUE KEY u (a) IGNORED, UNIQUE KEY n (a, c) IGNORED, UNIQUE KEY p (d(5)) IGNORED, UNIQUE KEY h (a) USING HASH IGNORED)
< ok 0 0
> query
| SELECT DISTINCT INDEX_NAME, IGNORED FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 't'
< result
< column "INDEX_NAME" "VARCHAR"
< column "IGNORED" "VARCHAR"
< row
: bytes
| v
: bytes
| NO
< row
: bytes
| u
: bytes
| YES
< row
: bytes
| p
: bytes
| YES
< row
: bytes
| n
: bytes
| YES
< row
: bytes
| h
: bytes
| YES
> query
| SHOW CREATE TABLE t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL,
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` varchar(10) NOT NULL,
|   UNIQUE KEY `v` (`b`),
|   UNIQUE KEY `u` (`a`) IGNORED,
|   UNIQUE KEY `p` (`d`(5)) IGNORED,
|   UNIQUE KEY `n` (`a`,`c`) IGNORED,
|   UNIQUE KEY `h` (`a`) USING HASH IGNORED
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL UNIQUE, UNIQUE KEY u (a) IGNORED)
< ok 0 0
> query
| SELECT DISTINCT INDEX_NAME, IGNORED FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 't'
< result
< column "INDEX_NAME" "VARCHAR"
< column "IGNORED" "VARCHAR"
< row
: bytes
| b
: bytes
| NO
< row
: bytes
| u
: bytes
| YES
> query
| SHOW CREATE TABLE t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL,
|   `b` int(11) NOT NULL,
|   UNIQUE KEY `b` (`b`),
|   UNIQUE KEY `u` (`a`) IGNORED
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT NOT NULL, UNIQUE KEY u (a) IGNORED, UNIQUE KEY w (b))
< ok 0 0
> exec
| ALTER TABLE t MODIFY a INT NOT NULL
< ok 0 0
> query
| SELECT DISTINCT INDEX_NAME, IGNORED FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 't'
< result
< column "INDEX_NAME" "VARCHAR"
< column "IGNORED" "VARCHAR"
< row
: bytes
| w
: bytes
| NO
< row
: bytes
| u
: bytes
| YES
> query
| SHOW CREATE TABLE t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL,
|   `b` int(11) NOT NULL,
|   UNIQUE KEY `w` (`b`),
|   UNIQUE KEY `u` (`a`) IGNORED
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY v (b), UNIQUE KEY w (a))
< ok 0 0
> exec
| ALTER TABLE t ALTER INDEX w IGNORED
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE KEY x (a, b) IGNORED
< ok 0 0
> query
| SELECT DISTINCT INDEX_NAME, IGNORED FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 't'
< result
< column "INDEX_NAME" "VARCHAR"
< column "IGNORED" "VARCHAR"
< row
: bytes
| v
: bytes
| NO
< row
: bytes
| w
: bytes
| YES
< row
: bytes
| x
: bytes
| YES
> query
| SHOW CREATE TABLE t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) NOT NULL,
|   `b` int(11) NOT NULL,
|   UNIQUE KEY `v` (`b`),
|   UNIQUE KEY `w` (`a`) IGNORED,
|   UNIQUE KEY `x` (`a`,`b`) IGNORED
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
