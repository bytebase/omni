-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestContainer_Section_1_13_RegularIndexes default#1
> exec
| DROP TABLE IF EXISTS t_idx1
< ok 0 0
> exec
| CREATE TABLE t_idx1 (a INT, KEY `idx_a` (a))
< ok 0 0
> query
| SHOW CREATE TABLE t_idx1
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_idx1
: bytes
| CREATE TABLE `t_idx1` (
|   `a` int(11) DEFAULT NULL,
|   KEY `idx_a` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_idx2
< ok 0 0
> exec
| CREATE TABLE t_idx2 (a INT, KEY (a))
< ok 0 0
> query
| SHOW CREATE TABLE t_idx2
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_idx2
: bytes
| CREATE TABLE `t_idx2` (
|   `a` int(11) DEFAULT NULL,
|   KEY `a` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_idx3
< ok 0 0
> exec
| CREATE TABLE t_idx3 (a INT, b INT, c INT, KEY `idx_abc` (a, b, c))
< ok 0 0
> query
| SHOW CREATE TABLE t_idx3
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_idx3
: bytes
| CREATE TABLE `t_idx3` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL,
|   KEY `idx_abc` (`a`,`b`,`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_idx4
< ok 0 0
> exec
| CREATE TABLE t_idx4 (a VARCHAR(255), KEY `idx_a` (a(10)))
< ok 0 0
> query
| SHOW CREATE TABLE t_idx4
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_idx4
: bytes
| CREATE TABLE `t_idx4` (
|   `a` varchar(255) DEFAULT NULL,
|   KEY `idx_a` (`a`(10))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_idx5
< ok 0 0
> exec
| CREATE TABLE t_idx5 (a INT, KEY `idx_a` (a DESC))
< ok 0 0
> query
| SHOW CREATE TABLE t_idx5
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_idx5
: bytes
| CREATE TABLE `t_idx5` (
|   `a` int(11) DEFAULT NULL,
|   KEY `idx_a` (`a` DESC)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_idx6
< ok 0 0
> exec
| CREATE TABLE t_idx6 (a INT, b INT, KEY `idx_ab` (a ASC, b DESC))
< ok 0 0
> query
| SHOW CREATE TABLE t_idx6
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_idx6
: bytes
| CREATE TABLE `t_idx6` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   KEY `idx_ab` (`a`,`b` DESC)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_idx7
< ok 0 0
> exec
| CREATE TABLE t_idx7 (a INT, KEY `idx_a` (a) USING HASH) ENGINE=MEMORY
< ok 0 0
> query
| SHOW CREATE TABLE t_idx7
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_idx7
: bytes
| CREATE TABLE `t_idx7` (
|   `a` int(11) DEFAULT NULL,
|   KEY `idx_a` (`a`) USING HASH
| ) ENGINE=MEMORY DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_idx8
< ok 0 0
> exec
| CREATE TABLE t_idx8 (a INT, KEY `idx_a` (a) USING BTREE)
< ok 0 0
> query
| SHOW CREATE TABLE t_idx8
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_idx8
: bytes
| CREATE TABLE `t_idx8` (
|   `a` int(11) DEFAULT NULL,
|   KEY `idx_a` (`a`) USING BTREE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
