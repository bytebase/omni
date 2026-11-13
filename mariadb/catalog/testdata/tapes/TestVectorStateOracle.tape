-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestVectorStateOracle default#1
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M=6 M=DEFAULT)
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
|   `v` vector(3) NOT NULL,
|   VECTOR KEY `i` (`v`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M=6 M=7)
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
|   `v` vector(3) NOT NULL,
|   VECTOR KEY `i` (`v`) `M`=7
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M=6 DISTANCE=cosine M=7)
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
|   `v` vector(3) NOT NULL,
|   VECTOR KEY `i` (`v`) `DISTANCE`=cosine `M`=7
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) m=6 M=DEFAULT DISTANCE=cosine)
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
|   `v` vector(3) NOT NULL,
|   VECTOR KEY `i` (`v`) `DISTANCE`=cosine
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, a INT, VECTOR INDEX i (v) M=6, KEY k (a))
< ok 0 0
> exec
| CREATE OR REPLACE VECTOR INDEX i ON t (v) M=8
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
|   `v` vector(3) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   KEY `k` (`a`),
|   VECTOR KEY `i` (`v`) `M`=8
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M=6)
< ok 0 0
> exec
| CREATE OR REPLACE VECTOR INDEX j ON t (v) M=9
< error
| 1235 42000 This version of MariaDB doesn't yet support 'multiple VECTOR indexes'
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
|   `v` vector(3) NOT NULL,
|   VECTOR KEY `i` (`v`) `M`=6
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, a INT, VECTOR INDEX i (v))
< ok 0 0
> exec
| CREATE OR REPLACE INDEX i ON t (a)
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
|   `v` vector(3) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   KEY `i` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, a INT, KEY k (a))
< ok 0 0
> exec
| CREATE OR REPLACE VECTOR INDEX k ON t (v)
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
|   `v` vector(3) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   VECTOR KEY `k` (`v`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, KEY a (a), KEY b (b))
< ok 0 0
> exec
| CREATE OR REPLACE INDEX a ON t (b)
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
|   `b` int(11) DEFAULT NULL,
|   KEY `b` (`b`),
|   KEY `a` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT, UNIQUE KEY u (a))
< ok 0 0
> exec
| CREATE OR REPLACE INDEX u ON t (a)
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
|   KEY `u` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY k (a) COMMENT 'c')
< ok 0 0
> exec
| CREATE OR REPLACE INDEX k ON t (nosuch)
< error
| 1072 42000 Key column 'nosuch' doesn't exist in table
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
|   KEY `k` (`a`) COMMENT 'c'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, KEY k (a) COMMENT 'c')
< ok 0 0
> exec
| CREATE OR REPLACE VECTOR INDEX k ON t (a)
< error
| 1210 HY000 Incorrect arguments to VECTOR INDEX
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
|   KEY `k` (`a`) COMMENT 'c'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (id INT)
< ok 0 0
> exec
| ALTER TABLE t ADD VECTOR i (v), ADD COLUMN v VECTOR(3) NOT NULL
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
|   `id` int(11) DEFAULT NULL,
|   `v` vector(3) NOT NULL,
|   VECTOR KEY `i` (`v`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (id INT, v VECTOR(3))
< ok 0 0
> exec
| ALTER TABLE t ADD VECTOR i (v), MODIFY v VECTOR(3) NOT NULL
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
|   `id` int(11) DEFAULT NULL,
|   `v` vector(3) NOT NULL,
|   VECTOR KEY `i` (`v`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (id INT, v VECTOR(3) NOT NULL)
< ok 0 0
> exec
| ALTER TABLE t ADD VECTOR i (v), DROP COLUMN v
< error
| 1072 42000 Key column 'v' doesn't exist in table
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (id INT, w VECTOR(3) NOT NULL)
< ok 0 0
> exec
| ALTER TABLE t ADD VECTOR i (w), CHANGE w z VECTOR(3) NOT NULL
< error
| 1072 42000 Key column 'w' doesn't exist in table
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (id INT)
< ok 0 0
> exec
| ALTER TABLE t ADD VECTOR i (w), ADD COLUMN w INT
< error
| 1210 HY000 Incorrect arguments to VECTOR INDEX
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (id INT)
< ok 0 0
> exec
| ALTER TABLE t ADD VECTOR i (w), ADD COLUMN w VECTOR(3) NULL
< error
| 1252 42000 All parts of a VECTOR index must be NOT NULL
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v))
< ok 0 0
> exec
| ALTER TABLE t ADD VECTOR j (v), DROP INDEX i
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
|   `v` vector(3) NOT NULL,
|   VECTOR KEY `j` (`v`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, w VECTOR(3) NOT NULL)
< ok 0 0
> exec
| ALTER TABLE t ADD VECTOR k (v), ADD VECTOR l (w)
< error
| 1235 42000 This version of MariaDB doesn't yet support 'multiple VECTOR indexes'
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, a INT)
< ok 0 0
> exec
| ALTER TABLE t ADD VECTOR i (v), ADD KEY z (a)
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
|   `v` vector(3) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   KEY `z` (`a`),
|   VECTOR KEY `i` (`v`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M=6)
< ok 0 0
> exec
| ALTER TABLE t ADD VECTOR IF NOT EXISTS i (v), ADD VECTOR IF NOT EXISTS i (v)
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
|   `v` vector(3) NOT NULL,
|   VECTOR KEY `i` (`v`) `M`=6
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL)
< ok 0 0
> exec
| ALTER TABLE t ADD VECTOR IF NOT EXISTS i (v), ADD VECTOR IF NOT EXISTS i (v)
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
|   `v` vector(3) NOT NULL,
|   VECTOR KEY `i` (`v`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL)
< ok 0 0
> exec
| ALTER TABLE t ADD VECTOR i (v), ADD VECTOR IF NOT EXISTS i (v)
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
|   `v` vector(3) NOT NULL,
|   VECTOR KEY `i` (`v`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL)
< ok 0 0
> exec
| ALTER TABLE t ADD VECTOR IF NOT EXISTS i (v), ADD VECTOR i (v)
< error
| 1061 42000 Duplicate key name 'i'
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, a INT)
< ok 0 0
> exec
| ALTER TABLE t ADD VECTOR IF NOT EXISTS i (v), ADD INDEX i (a)
< error
| 1061 42000 Duplicate key name 'i'
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, a INT)
< ok 0 0
> exec
| ALTER TABLE t ADD INDEX i (a), ADD VECTOR IF NOT EXISTS i (v)
< error
| 1061 42000 Duplicate key name 'i'
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (id INT)
< ok 0 0
> exec
| ALTER TABLE t ADD VECTOR IF NOT EXISTS i (v), ADD COLUMN v VECTOR(3) NOT NULL
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
|   `id` int(11) DEFAULT NULL,
|   `v` vector(3) NOT NULL,
|   VECTOR KEY `i` (`v`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
