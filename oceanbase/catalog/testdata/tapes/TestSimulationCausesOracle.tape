-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestSimulationCausesOracle default#1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation
< ok 0 0
> exec
| CREATE DATABASE omni_t_simulation DEFAULT CHARSET=utf8mb4; USE omni_t_simulation;
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by key()
| (partition `p0`)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation
< ok 2 0
> exec
| CREATE DATABASE omni_t_simulation DEFAULT CHARSET=utf8mb4; USE omni_t_simulation;
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()(PARTITION a)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by key()
| (partition `a`)
> exec
| DROP DATABASE IF EXISTS omni_t_simulation
< ok 2 0
> exec
| CREATE DATABASE omni_t_simulation DEFAULT CHARSET=utf8mb4; USE omni_t_simulation;
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE()
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by range() size ('2048MB')
> exec
| DROP DATABASE IF EXISTS omni_t_simulation
< ok 2 0
> exec
| CREATE DATABASE omni_t_simulation DEFAULT CHARSET=utf8mb4; USE omni_t_simulation;
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE COLUMNS()
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
|  partition by range() size ('2048MB')
> exec
| DROP DATABASE IF EXISTS omni_t_simulation
< ok 2 0
> exec
| CREATE DATABASE omni_t_simulation DEFAULT CHARSET=utf8mb4; USE omni_t_simulation;
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()(PARTITION a, PARTITION a)
< error
| 1517 HY000 Duplicate partition name 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation
< ok 2 0
> exec
| CREATE DATABASE omni_t_simulation DEFAULT CHARSET=utf8mb4; USE omni_t_simulation;
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE(a)
< error
| 1235 0A000 mismatching between primary key prefix and partition key is not supported
> exec
| DROP DATABASE IF EXISTS omni_t_simulation
< ok 2 0
> exec
| CREATE DATABASE omni_t_simulation DEFAULT CHARSET=utf8mb4; USE omni_t_simulation;
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CHARSET a
< error
| 1115 42000 Unknown character set: 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation
< ok 2 0
> exec
| CREATE DATABASE omni_t_simulation DEFAULT CHARSET=utf8mb4; USE omni_t_simulation;
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CONVERT TO CHARACTER SET a
< error
| 1115 42000 Unknown character set: 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation
< ok 2 0
> exec
| CREATE DATABASE omni_t_simulation DEFAULT CHARSET=utf8mb4; USE omni_t_simulation;
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a COLLATE a
< error
| 1273 HY000 Unknown collation: 'a'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation
< ok 2 0
> exec
| CREATE DATABASE omni_t_simulation DEFAULT CHARSET=utf8mb4; USE omni_t_simulation;
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SERIAL
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   `a` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `a` (`a`) BLOCK_SIZE 16384 LOCAL,
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX AUTO_INCREMENT = 1 AUTO_INCREMENT_MODE = 'ORDER' DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation
< ok 3 0
> exec
| CREATE DATABASE omni_t_simulation DEFAULT CHARSET=utf8mb4; USE omni_t_simulation;
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a) STORING(a)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL,
|   KEY `a` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS omni_t_simulation
< ok 3 0
