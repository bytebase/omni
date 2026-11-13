-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_ConstraintOfIndexSurvivesOtherNames default#1
> exec
| DROP DATABASE IF EXISTS zz_coi
< ok 0 0
> exec
| CREATE DATABASE zz_coi
< ok 1 0
> exec
| USE zz_coi
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY i (a), CONSTRAINT i CHECK (a > 0))
< ok 0 0
> exec
| ALTER TABLE t DROP INDEX i
< ok 0 0
> query
| SHOW CREATE TABLE t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   CONSTRAINT `i` CHECK ((`a` > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS zz_coi
< ok 1 0
> exec
| CREATE DATABASE zz_coi
< ok 1 0
> exec
| USE zz_coi
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY i (a), CONSTRAINT i CHECK (a > 0))
< ok 0 0
> exec
| DROP INDEX i ON t
< ok 0 0
> query
| SHOW CREATE TABLE t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   CONSTRAINT `i` CHECK ((`a` > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS zz_coi
< ok 1 0
> exec
| CREATE DATABASE zz_coi
< ok 1 0
> exec
| USE zz_coi
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, KEY j (a), KEY i (b), CONSTRAINT i FOREIGN KEY (a) REFERENCES p (id))
< ok 0 0
> exec
| DROP INDEX i ON t
< ok 0 0
> query
| SHOW CREATE TABLE t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   CONSTRAINT `i` FOREIGN KEY (`a`) REFERENCES `zz_coi`.`p`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT ,
|   KEY `j` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS zz_coi
< ok 3 0
> exec
| CREATE DATABASE zz_coi
< ok 1 0
> exec
| USE zz_coi
< ok 0 0
> exec
| CREATE TABLE t (a INT, CONSTRAINT i CHECK (a > 0), UNIQUE KEY i (a))
< ok 0 0
> exec
| DROP INDEX i ON t
< ok 0 0
> query
| SHOW CREATE TABLE t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   CONSTRAINT `i` CHECK ((`a` > 0))
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
