-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestContainerIgnoredIndexErrors default#1
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, PRIMARY KEY (a) IGNORED)
< error
| 4174 HY000 A primary key cannot be marked as IGNORE
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, CONSTRAINT pk PRIMARY KEY (a) IGNORED)
< error
| 4174 HY000 A primary key cannot be marked as IGNORE
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a) IGNORED)
< error
| 4174 HY000 A primary key cannot be marked as IGNORE
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY u (a) IGNORED, UNIQUE KEY v (b))
< error
| 4174 HY000 A primary key cannot be marked as IGNORE
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT NOT NULL, UNIQUE KEY u (a), UNIQUE KEY w (b) IGNORED)
< error
| 4174 HY000 A primary key cannot be marked as IGNORE
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT AUTO_INCREMENT, UNIQUE KEY u (a) IGNORED)
< error
| 4174 HY000 A primary key cannot be marked as IGNORE
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a DESC) IGNORED)
< error
| 4174 HY000 A primary key cannot be marked as IGNORE
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY v (b), UNIQUE KEY u (a) IGNORED)
< ok 0 0
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT, UNIQUE KEY u (a, b) IGNORED)
< ok 0 0
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a VARCHAR(10) NOT NULL, UNIQUE KEY u (a(5)) IGNORED)
< ok 0 0
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a) USING HASH IGNORED)
< ok 0 0
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a) IGNORED, b INT NOT NULL, PRIMARY KEY (b))
< ok 0 0
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
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
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| ALTER TABLE t ADD PRIMARY KEY (a) IGNORED
< error
| 4174 HY000 A primary key cannot be marked as IGNORE
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
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL PRIMARY KEY)
< ok 0 0
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
|   PRIMARY KEY (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| ALTER TABLE t ALTER INDEX `PRIMARY` NOT IGNORED
< error
| 4174 HY000 A primary key cannot be marked as IGNORE
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
|   PRIMARY KEY (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL PRIMARY KEY)
< ok 0 0
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
|   PRIMARY KEY (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| ALTER TABLE t ALTER INDEX IF EXISTS `PRIMARY` IGNORED
< error
| 4174 HY000 A primary key cannot be marked as IGNORE
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
|   PRIMARY KEY (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a))
< ok 0 0
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
|   UNIQUE KEY `u` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| ALTER TABLE t ALTER INDEX u IGNORED
< error
| 4174 HY000 A primary key cannot be marked as IGNORE
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
|   UNIQUE KEY `u` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a))
< ok 0 0
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
|   UNIQUE KEY `u` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| ALTER TABLE t ALTER INDEX u NOT IGNORED
< error
| 4174 HY000 A primary key cannot be marked as IGNORE
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
|   UNIQUE KEY `u` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
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
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| ALTER TABLE t ALTER INDEX `PRIMARY` IGNORED
< error
| 1176 42000 Key 'PRIMARY' doesn't exist in table 't'
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
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY k (a))
< ok 0 0
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
|   KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| ALTER TABLE t ALTER INDEX nosuch IGNORED
< error
| 1176 42000 Key 'nosuch' doesn't exist in table 't'
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
|   KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY k (a))
< ok 0 0
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
|   KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| ALTER TABLE t ALTER INDEX IF EXISTS nosuch IGNORED
< ok 0 0
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY k (a))
< ok 0 0
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
|   KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| ALTER TABLE t ALTER INDEX k IGNORED, DROP INDEX k
< error
| 1176 42000 Key 'k' doesn't exist in table 't'
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
|   KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY k (a))
< ok 0 0
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
|   KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| ALTER TABLE t DROP INDEX k, ALTER INDEX k IGNORED
< error
| 1176 42000 Key 'k' doesn't exist in table 't'
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
|   KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY k (a))
< ok 0 0
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
|   KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| ALTER TABLE t ALTER INDEX k IGNORED, RENAME INDEX k TO k2
< error
| 1176 42000 Key 'k' doesn't exist in table 't'
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
|   KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY k (a))
< ok 0 0
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
|   KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| ALTER TABLE t RENAME INDEX k TO k2, ALTER INDEX k2 IGNORED
< error
| 1176 42000 Key 'k2' doesn't exist in table 't'
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
|   KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
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
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| ALTER TABLE t ADD INDEX n (a), ALTER INDEX n IGNORED
< error
| 1176 42000 Key 'n' doesn't exist in table 't'
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
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY k (a))
< ok 0 0
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
|   KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| ALTER TABLE t ALTER INDEX IF EXISTS k IGNORED, DROP INDEX k
< error
| 1176 42000 Key 'k' doesn't exist in table 't'
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
|   KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY k (a))
< ok 0 0
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
|   KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| ALTER TABLE t ALTER INDEX IF EXISTS k IGNORED, RENAME INDEX k TO k2
< error
| 1176 42000 Key 'k' doesn't exist in table 't'
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
|   KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY k (a))
< ok 0 0
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
|   KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| ALTER TABLE t DROP INDEX k, ADD INDEX k (a), ALTER INDEX k IGNORED
< error
| 1176 42000 Key 'k' doesn't exist in table 't'
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
|   KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY k (a))
< ok 0 0
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
|   KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| ALTER TABLE t DROP INDEX k, ADD INDEX k (a), ALTER INDEX IF EXISTS k IGNORED
< error
| 1176 42000 Key 'k' doesn't exist in table 't'
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
|   KEY `k` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (a), UNIQUE KEY u (b) IGNORED)
< ok 0 0
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
|   PRIMARY KEY (`a`),
|   UNIQUE KEY `u` (`b`) IGNORED
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| ALTER TABLE t DROP PRIMARY KEY
< error
| 4174 HY000 A primary key cannot be marked as IGNORE
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
|   PRIMARY KEY (`a`),
|   UNIQUE KEY `u` (`b`) IGNORED
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT, UNIQUE KEY u (a) IGNORED)
< ok 0 0
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
|   UNIQUE KEY `u` (`a`) IGNORED
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| ALTER TABLE t MODIFY a INT NOT NULL
< error
| 4174 HY000 A primary key cannot be marked as IGNORE
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
|   UNIQUE KEY `u` (`a`) IGNORED
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY v (b), UNIQUE KEY w (a) IGNORED)
< ok 0 0
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
|   UNIQUE KEY `w` (`a`) IGNORED
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| ALTER TABLE t DROP INDEX v
< error
| 4174 HY000 A primary key cannot be marked as IGNORE
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
|   UNIQUE KEY `w` (`a`) IGNORED
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY v (b), UNIQUE KEY w (a) IGNORED)
< ok 0 0
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
|   UNIQUE KEY `w` (`a`) IGNORED
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP INDEX v ON t
< error
| 4174 HY000 A primary key cannot be marked as IGNORE
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
|   UNIQUE KEY `w` (`a`) IGNORED
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL)
< ok 0 0
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
|   `a` int(11) NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| CREATE UNIQUE INDEX u ON t (a) IGNORED
< error
| 4174 HY000 A primary key cannot be marked as IGNORE
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
|   `a` int(11) NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL)
< ok 0 0
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
|   `a` int(11) NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| ALTER TABLE t ADD UNIQUE KEY u (a) IGNORED
< error
| 4174 HY000 A primary key cannot be marked as IGNORE
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
|   `a` int(11) NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t, src
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY v (b))
< ok 0 0
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
|   UNIQUE KEY `v` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| ALTER TABLE t ADD UNIQUE KEY u (a) IGNORED
< ok 0 0
