-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestRangeColumnsMaxvalue_Oracle default#1
> exec
| CREATE DATABASE IF NOT EXISTS rc_test
< ok 0 0
> exec
| USE rc_test
< ok 0 0
> exec
| DROP TABLE IF EXISTS rc1
< ok 0 0
> exec
| CREATE TABLE rc1 (a INT, b INT) PARTITION BY RANGE COLUMNS(a, b) (PARTITION p0 VALUES LESS THAN (5, MAXVALUE), PARTITION p1 VALUES LESS THAN (MAXVALUE, MAXVALUE))
< ok 0 0
> query
| SHOW CREATE TABLE rc1
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| rc1
: bytes
| CREATE TABLE `rc1` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
| PARTITION BY RANGE COLUMNS(`a`,`b`)
| (PARTITION `p0` VALUES LESS THAN (5,MAXVALUE),
|  PARTITION `p1` VALUES LESS THAN (MAXVALUE,MAXVALUE))
> exec
| DROP TABLE IF EXISTS rc2
< ok 0 0
> exec
| CREATE TABLE rc2 (a INT, b VARCHAR(10), c DATE) PARTITION BY RANGE COLUMNS(a, b, c) (PARTITION p0 VALUES LESS THAN (1, 'm', MAXVALUE), PARTITION p1 VALUES LESS THAN (MAXVALUE, MAXVALUE, MAXVALUE))
< ok 0 0
> query
| SHOW CREATE TABLE rc2
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| rc2
: bytes
| CREATE TABLE `rc2` (
|   `a` int DEFAULT NULL,
|   `b` varchar(10) DEFAULT NULL,
|   `c` date DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
| PARTITION BY RANGE COLUMNS(`a`,`b`,`c`)
| (PARTITION `p0` VALUES LESS THAN (1,'m',MAXVALUE),
|  PARTITION `p1` VALUES LESS THAN (MAXVALUE,MAXVALUE,MAXVALUE))
> exec
| DROP TABLE IF EXISTS rc3
< ok 0 0
> exec
| CREATE TABLE rc3 (a INT, b INT) PARTITION BY RANGE COLUMNS(a, b) (PARTITION p0 VALUES LESS THAN (5, 10))
< ok 0 0
> exec
| ALTER TABLE rc3 ADD PARTITION (PARTITION p1 VALUES LESS THAN (10, MAXVALUE))
< ok 0 0
> query
| SHOW CREATE TABLE rc3
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| rc3
: bytes
| CREATE TABLE `rc3` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
| PARTITION BY RANGE COLUMNS(`a`,`b`)
| (PARTITION `p0` VALUES LESS THAN (5,10),
|  PARTITION `p1` VALUES LESS THAN (10,MAXVALUE))
> exec
| DROP TABLE IF EXISTS rc4
< ok 0 0
> exec
| CREATE TABLE rc4 (a INT, b INT) PARTITION BY RANGE COLUMNS(a, b) (PARTITION p0 VALUES LESS THAN (5, MAXVALUE), PARTITION p1 VALUES LESS THAN (6, 1))
< ok 0 0
> query
| SHOW CREATE TABLE rc4
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| rc4
: bytes
| CREATE TABLE `rc4` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
| PARTITION BY RANGE COLUMNS(`a`,`b`)
| (PARTITION `p0` VALUES LESS THAN (5,MAXVALUE),
|  PARTITION `p1` VALUES LESS THAN (6,1))
> exec
| DROP TABLE IF EXISTS rc5
< ok 0 0
> exec
| CREATE TABLE rc5 (a VARCHAR(10), b INT) PARTITION BY RANGE COLUMNS(a, b) (PARTITION p0 VALUES LESS THAN ('a', MAXVALUE), PARTITION p1 VALUES LESS THAN ('b', 1))
< ok 0 0
> query
| SHOW CREATE TABLE rc5
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| rc5
: bytes
| CREATE TABLE `rc5` (
|   `a` varchar(10) DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
| PARTITION BY RANGE COLUMNS(`a`,`b`)
| (PARTITION `p0` VALUES LESS THAN ('a',MAXVALUE),
|  PARTITION `p1` VALUES LESS THAN ('b',1))
> exec
| DROP DATABASE IF EXISTS rc_test
< ok 0 0
