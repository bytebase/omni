-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestSequenceSimulation default#1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a
< error
| 4314 HY000 no options specified for ALTER SEQUENCE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP SEQUENCE a
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE IF NOT EXISTS a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a NOMAXVALUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a NOMAXVALUE NOMAXVALUE
< error
| 4301 HY000 duplicate MAXVALUE/NOMAXVALUE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a NOMINVALUE NOMAXVALUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 1
< error
| 4320 HY000 MINVALUE must be less than MAXVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a NOMINVALUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a CYCLE
< error
| 4330 HY000 ascending sequences that CYCLE must specify MAXVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a NOCYCLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 1
< error
| 4326 HY000 the number of values to CACHE must be greater than 1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a NOCACHE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 1
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a ORDER
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| Y
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a NOORDER
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a RESTART
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your OceanBase version for the right syntax to use
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE + 1
< error
| 4320 HY000 MINVALUE must be less than MAXVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE - 1
< error
| 4320 HY000 MINVALUE must be less than MAXVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE + 1.5
< error
| 4317 HY000 sequence parameter MAXVALUE must be an integer
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE - 1.5
< error
| 4317 HY000 sequence parameter MAXVALUE must be an integer
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 1.5
< error
| 4317 HY000 sequence parameter MAXVALUE must be an integer
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH 5 INCREMENT BY 2
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 5
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 2
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY -1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| -1
: bytes
| -9223372036854775808
: bytes
| -1
: bytes
| -1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE -10 INCREMENT BY -2
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| -1
: bytes
| -10
: bytes
| -1
: bytes
| -2
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY -3 START WITH 4
< error
| 4324 HY000 START WITH cannot be more than MAXVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH 5 INCREMENT BY 2 MINVALUE 1 MAXVALUE 100 CACHE 10 CYCLE ORDER
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 5
: bytes
| 1
: bytes
| 100
: bytes
| 2
: bytes
| 10
: bytes
| Y
: bytes
| Y
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 100 CYCLE CACHE 10
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 100
: bytes
| 1
: bytes
| 10
: bytes
| N
: bytes
| Y
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 2 CYCLE MAXVALUE 3
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 3
: bytes
| 1
: bytes
| 2
: bytes
| N
: bytes
| Y
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 100 CYCLE CACHE 200
< error
| 4328 HY000 number to CACHE must be less than one cycle
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 2 CYCLE MAXVALUE 2
< error
| 4328 HY000 number to CACHE must be less than one cycle
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 5 CYCLE
< error
| 4328 HY000 number to CACHE must be less than one cycle
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY -1 CYCLE
< error
| 4329 HY000 descending sequences that CYCLE must specify MINVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY -1 MINVALUE -100 CYCLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| -1
: bytes
| -100
: bytes
| -1
: bytes
| -1
: bytes
| 20
: bytes
| N
: bytes
| Y
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 0
< error
| 4326 HY000 the number of values to CACHE must be greater than 1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 2
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 2
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY 0
< error
| 4318 HY000 INCREMENT must be a nonzero integer
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY 0 CACHE 1
< error
| 4318 HY000 INCREMENT must be a nonzero integer
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 1 MAXVALUE 1
< error
| 4326 HY000 the number of values to CACHE must be greater than 1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH 0
< error
| 4322 HY000 START WITH cannot be less than MINVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH -5
< error
| 4322 HY000 START WITH cannot be less than MINVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH 100 MAXVALUE 10
< error
| 4324 HY000 START WITH cannot be more than MAXVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE 10 MAXVALUE 5
< error
| 4320 HY000 MINVALUE must be less than MAXVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE 9223372036854775807
< error
| 4320 HY000 MINVALUE must be less than MAXVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE 1 MAXVALUE 3 INCREMENT BY 5
< error
| 4321 HY000 INCREMENT must be less than MAXVALUE minus MINVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 9223372036854775808
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 99999999999999999999999999999
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE -99999999999999999999999999999 INCREMENT BY -1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| -1
: bytes
| -9223372036854775808
: bytes
| -1
: bytes
| -1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 10.0
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 10
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 99999999999999999999.5
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 1.5 CYCLE
< error
| 4317 HY000 sequence parameter MAXVALUE must be an integer
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE 1.5
< error
| 4317 HY000 sequence parameter MINVALUE must be an integer
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY 1.5
< error
| 4317 HY000 sequence parameter INCREMENT must be an integer
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH 1.5
< error
| 4317 HY000 sequence parameter START WITH must be an integer
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 1.5
< error
| 4317 HY000 sequence parameter CACHE must be an integer
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 1e5
< error
| 5114 42000 Invalid numeric
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a NOMAXVALUE INCREMENT BY -1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| -1
: bytes
| -9223372036854775808
: bytes
| -1
: bytes
| -1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY 1 INCREMENT BY 2
< error
| 4312 HY000 duplicate INCREMENT BY specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH 1 START WITH 2
< error
| 4313 HY000 duplicate START WITH specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 5 MAXVALUE 6
< error
| 4301 HY000 duplicate MAXVALUE/NOMAXVALUE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE 1 MINVALUE 2
< error
| 4302 HY000 duplicate MINVALUE/NOMINVALUE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a NOMINVALUE NOMINVALUE
< error
| 4302 HY000 duplicate MINVALUE/NOMINVALUE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 2 CACHE 3
< error
| 4304 HY000 duplicate CACHE/NOCACHE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a NOCACHE NOCACHE
< error
| 4304 HY000 duplicate CACHE/NOCACHE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a CYCLE CYCLE
< error
| 4303 HY000 duplicate CYCLE/NOCYCLE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a ORDER ORDER
< error
| 4305 HY000 duplicate ORDER/NOORDER specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a RESTART RESTART
< error
| 4382 HY000 duplicate RESTART specified
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 5 NOMAXVALUE
< error
| 4306 HY000 conflicting MAXVALUE/NOMAXVALUE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE 5 NOMINVALUE
< error
| 4307 HY000 conflicting MINVALUE/NOMINVALUE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 5 NOCACHE
< error
| 4309 HY000 conflicting CACHE/NOCACHE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a ORDER NOORDER
< error
| 4310 HY000 conflicting ORDER/NOORDER specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a NOCYCLE CYCLE
< error
| 4308 HY000 conflicting CYCLE/NOCYCLE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a CYCLE NOCYCLE ORDER NOORDER
< error
| 4310 HY000 conflicting ORDER/NOORDER specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE IF NOT EXISTS a CYCLE
< error
| 4330 HY000 ascending sequences that CYCLE must specify MAXVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE IF NOT EXISTS a CYCLE NOCYCLE
< error
| 4308 HY000 conflicting CYCLE/NOCYCLE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE nope.a
< error
| 1049 42000 Unknown database 'nope'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE nope.a CYCLE
< error
| 1049 42000 Unknown database 'nope'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE nope.a CYCLE NOCYCLE
< error
| 4308 HY000 conflicting CYCLE/NOCYCLE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE A
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE b
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| b
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
< error
| 1059 42000 Identifier name 'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb' is too long
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a NOMAXVALUE
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a NOMINVALUE
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a MAXVALUE 1
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a MAXVALUE 1000
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a MINVALUE 1
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a MINVALUE -50
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a START WITH 1
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a RESTART
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a RESTART START WITH 5
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a RESTART START WITH 1000
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a INCREMENT BY -1
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a INCREMENT BY 1
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a INCREMENT BY 0
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a INCREMENT BY -1 RESTART
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a CYCLE
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a NOCYCLE
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a CACHE 1
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a CACHE 2
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a CACHE 1000
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a NOCACHE
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a ORDER
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a NOORDER
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a MAXVALUE 1.5
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a NOMAXVALUE CYCLE
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a CYCLE NOCYCLE
< error
| 4308 HY000 conflicting CYCLE/NOCYCLE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a RESTART RESTART
< error
| 4382 HY000 duplicate RESTART specified
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE A NOCACHE
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE b NOCACHE
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE b
< error
| 4314 HY000 no options specified for ALTER SEQUENCE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE nope.a NOCACHE
< error
| 1049 42000 Unknown database 'nope'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
< error
| 1059 42000 Identifier name 'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb' is too long
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP SEQUENCE A
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP SEQUENCE b
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS b
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP SEQUENCE nope.a
< error
| 1049 42000 Unknown database 'nope'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS nope.a
< error
| 1049 42000 Unknown database 'nope'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
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
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a
< error
| 4314 HY000 no options specified for ALTER SEQUENCE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP SEQUENCE a
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE IF NOT EXISTS a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
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
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a NOMAXVALUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
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
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a NOMAXVALUE NOMAXVALUE
< error
| 4301 HY000 duplicate MAXVALUE/NOMAXVALUE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a NOMINVALUE NOMAXVALUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
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
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
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
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
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
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 1
< error
| 4320 HY000 MINVALUE must be less than MAXVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
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
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a NOMINVALUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
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
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a CYCLE
< error
| 4330 HY000 ascending sequences that CYCLE must specify MAXVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a NOCYCLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
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
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 1
< error
| 4326 HY000 the number of values to CACHE must be greater than 1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a NOCACHE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
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
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 1
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a ORDER
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
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
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| Y
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a NOORDER
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
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
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a RESTART
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your OceanBase version for the right syntax to use
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE + 1
< error
| 4320 HY000 MINVALUE must be less than MAXVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE - 1
< error
| 4320 HY000 MINVALUE must be less than MAXVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE + 1.5
< error
| 4317 HY000 sequence parameter MAXVALUE must be an integer
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE - 1.5
< error
| 4317 HY000 sequence parameter MAXVALUE must be an integer
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 1.5
< error
| 4317 HY000 sequence parameter MAXVALUE must be an integer
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
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
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH 5 INCREMENT BY 2
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
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
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 5
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 2
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY -1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
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
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| -1
: bytes
| -9223372036854775808
: bytes
| -1
: bytes
| -1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE -10 INCREMENT BY -2
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
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
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| -1
: bytes
| -10
: bytes
| -1
: bytes
| -2
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY -3 START WITH 4
< error
| 4324 HY000 START WITH cannot be more than MAXVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH 5 INCREMENT BY 2 MINVALUE 1 MAXVALUE 100 CACHE 10 CYCLE ORDER
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
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
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 5
: bytes
| 1
: bytes
| 100
: bytes
| 2
: bytes
| 10
: bytes
| Y
: bytes
| Y
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 100 CYCLE CACHE 10
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
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
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 100
: bytes
| 1
: bytes
| 10
: bytes
| N
: bytes
| Y
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 2 CYCLE MAXVALUE 3
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
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
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 3
: bytes
| 1
: bytes
| 2
: bytes
| N
: bytes
| Y
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 100 CYCLE CACHE 200
< error
| 4328 HY000 number to CACHE must be less than one cycle
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 2 CYCLE MAXVALUE 2
< error
| 4328 HY000 number to CACHE must be less than one cycle
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 5 CYCLE
< error
| 4328 HY000 number to CACHE must be less than one cycle
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY -1 CYCLE
< error
| 4329 HY000 descending sequences that CYCLE must specify MINVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY -1 MINVALUE -100 CYCLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
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
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| -1
: bytes
| -100
: bytes
| -1
: bytes
| -1
: bytes
| 20
: bytes
| N
: bytes
| Y
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 0
< error
| 4326 HY000 the number of values to CACHE must be greater than 1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 2
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
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
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 2
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY 0
< error
| 4318 HY000 INCREMENT must be a nonzero integer
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY 0 CACHE 1
< error
| 4318 HY000 INCREMENT must be a nonzero integer
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 1 MAXVALUE 1
< error
| 4326 HY000 the number of values to CACHE must be greater than 1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH 0
< error
| 4322 HY000 START WITH cannot be less than MINVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH -5
< error
| 4322 HY000 START WITH cannot be less than MINVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH 100 MAXVALUE 10
< error
| 4324 HY000 START WITH cannot be more than MAXVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE 10 MAXVALUE 5
< error
| 4320 HY000 MINVALUE must be less than MAXVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE 9223372036854775807
< error
| 4320 HY000 MINVALUE must be less than MAXVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE 1 MAXVALUE 3 INCREMENT BY 5
< error
| 4321 HY000 INCREMENT must be less than MAXVALUE minus MINVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 9223372036854775808
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
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
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 99999999999999999999999999999
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
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
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE -99999999999999999999999999999 INCREMENT BY -1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
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
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| -1
: bytes
| -9223372036854775808
: bytes
| -1
: bytes
| -1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 10.0
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
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
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 10
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 99999999999999999999.5
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
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
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 1.5 CYCLE
< error
| 4317 HY000 sequence parameter MAXVALUE must be an integer
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE 1.5
< error
| 4317 HY000 sequence parameter MINVALUE must be an integer
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY 1.5
< error
| 4317 HY000 sequence parameter INCREMENT must be an integer
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH 1.5
< error
| 4317 HY000 sequence parameter START WITH must be an integer
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 1.5
< error
| 4317 HY000 sequence parameter CACHE must be an integer
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 1e5
< error
| 5114 42000 Invalid numeric
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a NOMAXVALUE INCREMENT BY -1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
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
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| -1
: bytes
| -9223372036854775808
: bytes
| -1
: bytes
| -1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY 1 INCREMENT BY 2
< error
| 4312 HY000 duplicate INCREMENT BY specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH 1 START WITH 2
< error
| 4313 HY000 duplicate START WITH specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 5 MAXVALUE 6
< error
| 4301 HY000 duplicate MAXVALUE/NOMAXVALUE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE 1 MINVALUE 2
< error
| 4302 HY000 duplicate MINVALUE/NOMINVALUE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a NOMINVALUE NOMINVALUE
< error
| 4302 HY000 duplicate MINVALUE/NOMINVALUE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 2 CACHE 3
< error
| 4304 HY000 duplicate CACHE/NOCACHE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a NOCACHE NOCACHE
< error
| 4304 HY000 duplicate CACHE/NOCACHE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a CYCLE CYCLE
< error
| 4303 HY000 duplicate CYCLE/NOCYCLE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a ORDER ORDER
< error
| 4305 HY000 duplicate ORDER/NOORDER specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a RESTART RESTART
< error
| 4382 HY000 duplicate RESTART specified
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 5 NOMAXVALUE
< error
| 4306 HY000 conflicting MAXVALUE/NOMAXVALUE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE 5 NOMINVALUE
< error
| 4307 HY000 conflicting MINVALUE/NOMINVALUE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 5 NOCACHE
< error
| 4309 HY000 conflicting CACHE/NOCACHE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a ORDER NOORDER
< error
| 4310 HY000 conflicting ORDER/NOORDER specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a NOCYCLE CYCLE
< error
| 4308 HY000 conflicting CYCLE/NOCYCLE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE a CYCLE NOCYCLE ORDER NOORDER
< error
| 4310 HY000 conflicting ORDER/NOORDER specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE IF NOT EXISTS a CYCLE
< error
| 4330 HY000 ascending sequences that CYCLE must specify MAXVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE IF NOT EXISTS a CYCLE NOCYCLE
< error
| 4308 HY000 conflicting CYCLE/NOCYCLE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE nope.a
< error
| 1049 42000 Unknown database 'nope'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE nope.a CYCLE
< error
| 1049 42000 Unknown database 'nope'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE nope.a CYCLE NOCYCLE
< error
| 4308 HY000 conflicting CYCLE/NOCYCLE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE A
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
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
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE b
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
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
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| b
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
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
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SEQUENCE bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
< error
| 1059 42000 Identifier name 'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb' is too long
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a NOMAXVALUE
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a NOMINVALUE
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a MAXVALUE 1
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a MAXVALUE 1000
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a MINVALUE 1
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a MINVALUE -50
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a START WITH 1
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a RESTART
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a RESTART START WITH 5
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a RESTART START WITH 1000
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a INCREMENT BY -1
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a INCREMENT BY 1
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a INCREMENT BY 0
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a INCREMENT BY -1 RESTART
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a CYCLE
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a NOCYCLE
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a CACHE 1
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a CACHE 2
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a CACHE 1000
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a NOCACHE
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a ORDER
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a NOORDER
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a MAXVALUE 1.5
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a NOMAXVALUE CYCLE
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a CYCLE NOCYCLE
< error
| 4308 HY000 conflicting CYCLE/NOCYCLE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE a RESTART RESTART
< error
| 4382 HY000 duplicate RESTART specified
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE A NOCACHE
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE b NOCACHE
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE b
< error
| 4314 HY000 no options specified for ALTER SEQUENCE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE nope.a NOCACHE
< error
| 1049 42000 Unknown database 'nope'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER SEQUENCE bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
< error
| 1059 42000 Identifier name 'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb' is too long
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP SEQUENCE A
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP SEQUENCE b
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS b
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `a`
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
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP SEQUENCE nope.a
< error
| 1049 42000 Unknown database 'nope'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS nope.a
< error
| 1049 42000 Unknown database 'nope'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `a`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE VIEW `a` AS select `omni_t_simulation`.`t`.`b` AS `b`,`omni_t_simulation`.`t`.`a` AS `a`,`omni_t_simulation`.`t`.`d` AS `d` from `omni_t_simulation`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE a
< error
| 4314 HY000 no options specified for ALTER SEQUENCE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| DROP SEQUENCE a
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE IF NOT EXISTS a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `a`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE VIEW `a` AS select `omni_t_simulation`.`t`.`b` AS `b`,`omni_t_simulation`.`t`.`a` AS `a`,`omni_t_simulation`.`t`.`d` AS `d` from `omni_t_simulation`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a NOMAXVALUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `a`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE VIEW `a` AS select `omni_t_simulation`.`t`.`b` AS `b`,`omni_t_simulation`.`t`.`a` AS `a`,`omni_t_simulation`.`t`.`d` AS `d` from `omni_t_simulation`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a NOMAXVALUE NOMAXVALUE
< error
| 4301 HY000 duplicate MAXVALUE/NOMAXVALUE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a NOMINVALUE NOMAXVALUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `a`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE VIEW `a` AS select `omni_t_simulation`.`t`.`b` AS `b`,`omni_t_simulation`.`t`.`a` AS `a`,`omni_t_simulation`.`t`.`d` AS `d` from `omni_t_simulation`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `a`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE VIEW `a` AS select `omni_t_simulation`.`t`.`b` AS `b`,`omni_t_simulation`.`t`.`a` AS `a`,`omni_t_simulation`.`t`.`d` AS `d` from `omni_t_simulation`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `a`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE VIEW `a` AS select `omni_t_simulation`.`t`.`b` AS `b`,`omni_t_simulation`.`t`.`a` AS `a`,`omni_t_simulation`.`t`.`d` AS `d` from `omni_t_simulation`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 1
< error
| 4320 HY000 MINVALUE must be less than MAXVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `a`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE VIEW `a` AS select `omni_t_simulation`.`t`.`b` AS `b`,`omni_t_simulation`.`t`.`a` AS `a`,`omni_t_simulation`.`t`.`d` AS `d` from `omni_t_simulation`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a NOMINVALUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `a`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE VIEW `a` AS select `omni_t_simulation`.`t`.`b` AS `b`,`omni_t_simulation`.`t`.`a` AS `a`,`omni_t_simulation`.`t`.`d` AS `d` from `omni_t_simulation`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a CYCLE
< error
| 4330 HY000 ascending sequences that CYCLE must specify MAXVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a NOCYCLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `a`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE VIEW `a` AS select `omni_t_simulation`.`t`.`b` AS `b`,`omni_t_simulation`.`t`.`a` AS `a`,`omni_t_simulation`.`t`.`d` AS `d` from `omni_t_simulation`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 1
< error
| 4326 HY000 the number of values to CACHE must be greater than 1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a NOCACHE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `a`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE VIEW `a` AS select `omni_t_simulation`.`t`.`b` AS `b`,`omni_t_simulation`.`t`.`a` AS `a`,`omni_t_simulation`.`t`.`d` AS `d` from `omni_t_simulation`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 1
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a ORDER
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `a`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE VIEW `a` AS select `omni_t_simulation`.`t`.`b` AS `b`,`omni_t_simulation`.`t`.`a` AS `a`,`omni_t_simulation`.`t`.`d` AS `d` from `omni_t_simulation`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| Y
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a NOORDER
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `a`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE VIEW `a` AS select `omni_t_simulation`.`t`.`b` AS `b`,`omni_t_simulation`.`t`.`a` AS `a`,`omni_t_simulation`.`t`.`d` AS `d` from `omni_t_simulation`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a RESTART
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your OceanBase version for the right syntax to use
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE + 1
< error
| 4320 HY000 MINVALUE must be less than MAXVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE - 1
< error
| 4320 HY000 MINVALUE must be less than MAXVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE + 1.5
< error
| 4317 HY000 sequence parameter MAXVALUE must be an integer
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE - 1.5
< error
| 4317 HY000 sequence parameter MAXVALUE must be an integer
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 1.5
< error
| 4317 HY000 sequence parameter MAXVALUE must be an integer
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `a`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE VIEW `a` AS select `omni_t_simulation`.`t`.`b` AS `b`,`omni_t_simulation`.`t`.`a` AS `a`,`omni_t_simulation`.`t`.`d` AS `d` from `omni_t_simulation`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH 5 INCREMENT BY 2
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `a`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE VIEW `a` AS select `omni_t_simulation`.`t`.`b` AS `b`,`omni_t_simulation`.`t`.`a` AS `a`,`omni_t_simulation`.`t`.`d` AS `d` from `omni_t_simulation`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 5
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 2
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY -1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `a`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE VIEW `a` AS select `omni_t_simulation`.`t`.`b` AS `b`,`omni_t_simulation`.`t`.`a` AS `a`,`omni_t_simulation`.`t`.`d` AS `d` from `omni_t_simulation`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| -1
: bytes
| -9223372036854775808
: bytes
| -1
: bytes
| -1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE -10 INCREMENT BY -2
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `a`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE VIEW `a` AS select `omni_t_simulation`.`t`.`b` AS `b`,`omni_t_simulation`.`t`.`a` AS `a`,`omni_t_simulation`.`t`.`d` AS `d` from `omni_t_simulation`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| -1
: bytes
| -10
: bytes
| -1
: bytes
| -2
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY -3 START WITH 4
< error
| 4324 HY000 START WITH cannot be more than MAXVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH 5 INCREMENT BY 2 MINVALUE 1 MAXVALUE 100 CACHE 10 CYCLE ORDER
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `a`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE VIEW `a` AS select `omni_t_simulation`.`t`.`b` AS `b`,`omni_t_simulation`.`t`.`a` AS `a`,`omni_t_simulation`.`t`.`d` AS `d` from `omni_t_simulation`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 5
: bytes
| 1
: bytes
| 100
: bytes
| 2
: bytes
| 10
: bytes
| Y
: bytes
| Y
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 100 CYCLE CACHE 10
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `a`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE VIEW `a` AS select `omni_t_simulation`.`t`.`b` AS `b`,`omni_t_simulation`.`t`.`a` AS `a`,`omni_t_simulation`.`t`.`d` AS `d` from `omni_t_simulation`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 100
: bytes
| 1
: bytes
| 10
: bytes
| N
: bytes
| Y
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 2 CYCLE MAXVALUE 3
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `a`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE VIEW `a` AS select `omni_t_simulation`.`t`.`b` AS `b`,`omni_t_simulation`.`t`.`a` AS `a`,`omni_t_simulation`.`t`.`d` AS `d` from `omni_t_simulation`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 3
: bytes
| 1
: bytes
| 2
: bytes
| N
: bytes
| Y
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 100 CYCLE CACHE 200
< error
| 4328 HY000 number to CACHE must be less than one cycle
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 2 CYCLE MAXVALUE 2
< error
| 4328 HY000 number to CACHE must be less than one cycle
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 5 CYCLE
< error
| 4328 HY000 number to CACHE must be less than one cycle
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY -1 CYCLE
< error
| 4329 HY000 descending sequences that CYCLE must specify MINVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY -1 MINVALUE -100 CYCLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `a`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE VIEW `a` AS select `omni_t_simulation`.`t`.`b` AS `b`,`omni_t_simulation`.`t`.`a` AS `a`,`omni_t_simulation`.`t`.`d` AS `d` from `omni_t_simulation`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| -1
: bytes
| -100
: bytes
| -1
: bytes
| -1
: bytes
| 20
: bytes
| N
: bytes
| Y
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 0
< error
| 4326 HY000 the number of values to CACHE must be greater than 1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 2
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `a`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE VIEW `a` AS select `omni_t_simulation`.`t`.`b` AS `b`,`omni_t_simulation`.`t`.`a` AS `a`,`omni_t_simulation`.`t`.`d` AS `d` from `omni_t_simulation`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 2
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY 0
< error
| 4318 HY000 INCREMENT must be a nonzero integer
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY 0 CACHE 1
< error
| 4318 HY000 INCREMENT must be a nonzero integer
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 1 MAXVALUE 1
< error
| 4326 HY000 the number of values to CACHE must be greater than 1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH 0
< error
| 4322 HY000 START WITH cannot be less than MINVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH -5
< error
| 4322 HY000 START WITH cannot be less than MINVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH 100 MAXVALUE 10
< error
| 4324 HY000 START WITH cannot be more than MAXVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE 10 MAXVALUE 5
< error
| 4320 HY000 MINVALUE must be less than MAXVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE 9223372036854775807
< error
| 4320 HY000 MINVALUE must be less than MAXVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE 1 MAXVALUE 3 INCREMENT BY 5
< error
| 4321 HY000 INCREMENT must be less than MAXVALUE minus MINVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 9223372036854775808
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `a`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE VIEW `a` AS select `omni_t_simulation`.`t`.`b` AS `b`,`omni_t_simulation`.`t`.`a` AS `a`,`omni_t_simulation`.`t`.`d` AS `d` from `omni_t_simulation`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 99999999999999999999999999999
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `a`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE VIEW `a` AS select `omni_t_simulation`.`t`.`b` AS `b`,`omni_t_simulation`.`t`.`a` AS `a`,`omni_t_simulation`.`t`.`d` AS `d` from `omni_t_simulation`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE -99999999999999999999999999999 INCREMENT BY -1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `a`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE VIEW `a` AS select `omni_t_simulation`.`t`.`b` AS `b`,`omni_t_simulation`.`t`.`a` AS `a`,`omni_t_simulation`.`t`.`d` AS `d` from `omni_t_simulation`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| -1
: bytes
| -9223372036854775808
: bytes
| -1
: bytes
| -1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 10.0
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `a`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE VIEW `a` AS select `omni_t_simulation`.`t`.`b` AS `b`,`omni_t_simulation`.`t`.`a` AS `a`,`omni_t_simulation`.`t`.`d` AS `d` from `omni_t_simulation`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 10
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 99999999999999999999.5
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `a`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE VIEW `a` AS select `omni_t_simulation`.`t`.`b` AS `b`,`omni_t_simulation`.`t`.`a` AS `a`,`omni_t_simulation`.`t`.`d` AS `d` from `omni_t_simulation`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 1.5 CYCLE
< error
| 4317 HY000 sequence parameter MAXVALUE must be an integer
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE 1.5
< error
| 4317 HY000 sequence parameter MINVALUE must be an integer
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY 1.5
< error
| 4317 HY000 sequence parameter INCREMENT must be an integer
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH 1.5
< error
| 4317 HY000 sequence parameter START WITH must be an integer
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 1.5
< error
| 4317 HY000 sequence parameter CACHE must be an integer
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 1e5
< error
| 5114 42000 Invalid numeric
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a NOMAXVALUE INCREMENT BY -1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `a`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE VIEW `a` AS select `omni_t_simulation`.`t`.`b` AS `b`,`omni_t_simulation`.`t`.`a` AS `a`,`omni_t_simulation`.`t`.`d` AS `d` from `omni_t_simulation`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| -1
: bytes
| -9223372036854775808
: bytes
| -1
: bytes
| -1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY 1 INCREMENT BY 2
< error
| 4312 HY000 duplicate INCREMENT BY specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH 1 START WITH 2
< error
| 4313 HY000 duplicate START WITH specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 5 MAXVALUE 6
< error
| 4301 HY000 duplicate MAXVALUE/NOMAXVALUE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE 1 MINVALUE 2
< error
| 4302 HY000 duplicate MINVALUE/NOMINVALUE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a NOMINVALUE NOMINVALUE
< error
| 4302 HY000 duplicate MINVALUE/NOMINVALUE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 2 CACHE 3
< error
| 4304 HY000 duplicate CACHE/NOCACHE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a NOCACHE NOCACHE
< error
| 4304 HY000 duplicate CACHE/NOCACHE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a CYCLE CYCLE
< error
| 4303 HY000 duplicate CYCLE/NOCYCLE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a ORDER ORDER
< error
| 4305 HY000 duplicate ORDER/NOORDER specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a RESTART RESTART
< error
| 4382 HY000 duplicate RESTART specified
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 5 NOMAXVALUE
< error
| 4306 HY000 conflicting MAXVALUE/NOMAXVALUE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE 5 NOMINVALUE
< error
| 4307 HY000 conflicting MINVALUE/NOMINVALUE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 5 NOCACHE
< error
| 4309 HY000 conflicting CACHE/NOCACHE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a ORDER NOORDER
< error
| 4310 HY000 conflicting ORDER/NOORDER specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a NOCYCLE CYCLE
< error
| 4308 HY000 conflicting CYCLE/NOCYCLE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE a CYCLE NOCYCLE ORDER NOORDER
< error
| 4310 HY000 conflicting ORDER/NOORDER specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE IF NOT EXISTS a CYCLE
< error
| 4330 HY000 ascending sequences that CYCLE must specify MAXVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE IF NOT EXISTS a CYCLE NOCYCLE
< error
| 4308 HY000 conflicting CYCLE/NOCYCLE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE nope.a
< error
| 1049 42000 Unknown database 'nope'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE nope.a CYCLE
< error
| 1049 42000 Unknown database 'nope'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE nope.a CYCLE NOCYCLE
< error
| 4308 HY000 conflicting CYCLE/NOCYCLE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE A
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `a`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE VIEW `a` AS select `omni_t_simulation`.`t`.`b` AS `b`,`omni_t_simulation`.`t`.`a` AS `a`,`omni_t_simulation`.`t`.`d` AS `d` from `omni_t_simulation`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE b
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `a`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE VIEW `a` AS select `omni_t_simulation`.`t`.`b` AS `b`,`omni_t_simulation`.`t`.`a` AS `a`,`omni_t_simulation`.`t`.`d` AS `d` from `omni_t_simulation`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| b
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `a`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE VIEW `a` AS select `omni_t_simulation`.`t`.`b` AS `b`,`omni_t_simulation`.`t`.`a` AS `a`,`omni_t_simulation`.`t`.`d` AS `d` from `omni_t_simulation`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| CREATE SEQUENCE bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
< error
| 1059 42000 Identifier name 'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb' is too long
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE a NOMAXVALUE
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE a NOMINVALUE
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE a MAXVALUE 1
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE a MAXVALUE 1000
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE a MINVALUE 1
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE a MINVALUE -50
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE a START WITH 1
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE a RESTART
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE a RESTART START WITH 5
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE a RESTART START WITH 1000
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE a INCREMENT BY -1
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE a INCREMENT BY 1
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE a INCREMENT BY 0
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE a INCREMENT BY -1 RESTART
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE a CYCLE
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE a NOCYCLE
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE a CACHE 1
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE a CACHE 2
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE a CACHE 1000
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE a NOCACHE
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE a ORDER
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE a NOORDER
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE a MAXVALUE 1.5
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE a NOMAXVALUE CYCLE
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE a CYCLE NOCYCLE
< error
| 4308 HY000 conflicting CYCLE/NOCYCLE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE a RESTART RESTART
< error
| 4382 HY000 duplicate RESTART specified
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE A NOCACHE
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE b NOCACHE
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE b
< error
| 4314 HY000 no options specified for ALTER SEQUENCE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE nope.a NOCACHE
< error
| 1049 42000 Unknown database 'nope'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| ALTER SEQUENCE bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
< error
| 1059 42000 Identifier name 'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb' is too long
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| DROP SEQUENCE A
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| DROP SEQUENCE b
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS b
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE VIEW `a`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "TEXT" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE VIEW `a` AS select `omni_t_simulation`.`t`.`b` AS `b`,`omni_t_simulation`.`t`.`a` AS `a`,`omni_t_simulation`.`t`.`d` AS `d` from `omni_t_simulation`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "TEXT" "string"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| DROP SEQUENCE nope.a
< error
| 1049 42000 Unknown database 'nope'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE TABLE t (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d));
| CREATE VIEW a AS SELECT b, a, d FROM t
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS nope.a
< error
| 1049 42000 Unknown database 'nope'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| ALTER SEQUENCE a
< error
| 4314 HY000 no options specified for ALTER SEQUENCE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| DROP SEQUENCE a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE IF NOT EXISTS a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a NOMAXVALUE
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a NOMAXVALUE NOMAXVALUE
< error
| 4301 HY000 duplicate MAXVALUE/NOMAXVALUE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a NOMINVALUE NOMAXVALUE
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY 1
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH 1
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 1
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE 1
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a NOMINVALUE
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a CYCLE
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a NOCYCLE
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 1
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a NOCACHE
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a ORDER
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a NOORDER
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a RESTART
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your OceanBase version for the right syntax to use
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE + 1
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE - 1
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE + 1.5
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE - 1.5
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 1.5
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH 5 INCREMENT BY 2
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY -1
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE -10 INCREMENT BY -2
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY -3 START WITH 4
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH 5 INCREMENT BY 2 MINVALUE 1 MAXVALUE 100 CACHE 10 CYCLE ORDER
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 100 CYCLE CACHE 10
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 2 CYCLE MAXVALUE 3
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 100 CYCLE CACHE 200
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 2 CYCLE MAXVALUE 2
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 5 CYCLE
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY -1 CYCLE
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY -1 MINVALUE -100 CYCLE
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 0
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 2
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY 0
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY 0 CACHE 1
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 1 MAXVALUE 1
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH 0
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH -5
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH 100 MAXVALUE 10
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE 10 MAXVALUE 5
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE 9223372036854775807
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE 1 MAXVALUE 3 INCREMENT BY 5
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 9223372036854775808
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 99999999999999999999999999999
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE -99999999999999999999999999999 INCREMENT BY -1
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 10.0
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 99999999999999999999.5
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 1.5 CYCLE
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE 1.5
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY 1.5
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH 1.5
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 1.5
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 1e5
< error
| 5114 42000 Invalid numeric
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a NOMAXVALUE INCREMENT BY -1
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY 1 INCREMENT BY 2
< error
| 4312 HY000 duplicate INCREMENT BY specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH 1 START WITH 2
< error
| 4313 HY000 duplicate START WITH specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 5 MAXVALUE 6
< error
| 4301 HY000 duplicate MAXVALUE/NOMAXVALUE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE 1 MINVALUE 2
< error
| 4302 HY000 duplicate MINVALUE/NOMINVALUE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a NOMINVALUE NOMINVALUE
< error
| 4302 HY000 duplicate MINVALUE/NOMINVALUE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 2 CACHE 3
< error
| 4304 HY000 duplicate CACHE/NOCACHE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a NOCACHE NOCACHE
< error
| 4304 HY000 duplicate CACHE/NOCACHE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a CYCLE CYCLE
< error
| 4303 HY000 duplicate CYCLE/NOCYCLE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a ORDER ORDER
< error
| 4305 HY000 duplicate ORDER/NOORDER specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a RESTART RESTART
< error
| 4382 HY000 duplicate RESTART specified
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 5 NOMAXVALUE
< error
| 4306 HY000 conflicting MAXVALUE/NOMAXVALUE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE 5 NOMINVALUE
< error
| 4307 HY000 conflicting MINVALUE/NOMINVALUE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 5 NOCACHE
< error
| 4309 HY000 conflicting CACHE/NOCACHE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a ORDER NOORDER
< error
| 4310 HY000 conflicting ORDER/NOORDER specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a NOCYCLE CYCLE
< error
| 4308 HY000 conflicting CYCLE/NOCYCLE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE a CYCLE NOCYCLE ORDER NOORDER
< error
| 4310 HY000 conflicting ORDER/NOORDER specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE IF NOT EXISTS a CYCLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE IF NOT EXISTS a CYCLE NOCYCLE
< error
| 4308 HY000 conflicting CYCLE/NOCYCLE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE nope.a
< error
| 1049 42000 Unknown database 'nope'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE nope.a CYCLE
< error
| 1049 42000 Unknown database 'nope'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE nope.a CYCLE NOCYCLE
< error
| 4308 HY000 conflicting CYCLE/NOCYCLE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE A
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE b
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
< row
: bytes
| b
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
< row
: bytes
| bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| CREATE SEQUENCE bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
< error
| 1059 42000 Identifier name 'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb' is too long
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| ALTER SEQUENCE a NOMAXVALUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775808
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| ALTER SEQUENCE a NOMINVALUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| -9223372036854775809
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| ALTER SEQUENCE a MAXVALUE 1
< error
| 4320 HY000 MINVALUE must be less than MAXVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| ALTER SEQUENCE a MAXVALUE 1000
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 1000
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| ALTER SEQUENCE a MINVALUE 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| ALTER SEQUENCE a MINVALUE -50
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| -50
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| ALTER SEQUENCE a START WITH 1
< error
| 4311 HY000 cannot alter starting sequence number
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| ALTER SEQUENCE a RESTART
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| ALTER SEQUENCE a RESTART START WITH 5
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 5
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| ALTER SEQUENCE a RESTART START WITH 1000
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1000
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| ALTER SEQUENCE a INCREMENT BY -1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| -1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| ALTER SEQUENCE a INCREMENT BY 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| ALTER SEQUENCE a INCREMENT BY 0
< error
| 4318 HY000 INCREMENT must be a nonzero integer
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| ALTER SEQUENCE a INCREMENT BY -1 RESTART
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| -1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| ALTER SEQUENCE a CYCLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| Y
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| ALTER SEQUENCE a NOCYCLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| ALTER SEQUENCE a CACHE 1
< error
| 4326 HY000 the number of values to CACHE must be greater than 1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| ALTER SEQUENCE a CACHE 2
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 2
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| ALTER SEQUENCE a CACHE 1000
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 1000
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| ALTER SEQUENCE a NOCACHE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 1
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| ALTER SEQUENCE a ORDER
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| Y
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| ALTER SEQUENCE a NOORDER
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| ALTER SEQUENCE a MAXVALUE 1.5
< error
| 4317 HY000 sequence parameter MAXVALUE must be an integer
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| ALTER SEQUENCE a NOMAXVALUE CYCLE
< error
| 4330 HY000 ascending sequences that CYCLE must specify MAXVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| ALTER SEQUENCE a CYCLE NOCYCLE
< error
| 4308 HY000 conflicting CYCLE/NOCYCLE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| ALTER SEQUENCE a RESTART RESTART
< error
| 4382 HY000 duplicate RESTART specified
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| ALTER SEQUENCE A NOCACHE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 1
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| ALTER SEQUENCE b NOCACHE
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| ALTER SEQUENCE b
< error
| 4314 HY000 no options specified for ALTER SEQUENCE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| ALTER SEQUENCE nope.a NOCACHE
< error
| 1049 42000 Unknown database 'nope'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| ALTER SEQUENCE bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
< error
| 1059 42000 Identifier name 'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb' is too long
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| DROP SEQUENCE A
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| DROP SEQUENCE b
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS b
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| DROP SEQUENCE nope.a
< error
| 1049 42000 Unknown database 'nope'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS nope.a
< error
| 1049 42000 Unknown database 'nope'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| ALTER SEQUENCE a
< error
| 4314 HY000 no options specified for ALTER SEQUENCE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| DROP SEQUENCE a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE IF NOT EXISTS a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 10
: bytes
| -100
: bytes
| 100
: bytes
| -2
: bytes
| 5
: bytes
| Y
: bytes
| Y
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a NOMAXVALUE
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a NOMAXVALUE NOMAXVALUE
< error
| 4301 HY000 duplicate MAXVALUE/NOMAXVALUE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a NOMINVALUE NOMAXVALUE
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY 1
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH 1
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 1
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE 1
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a NOMINVALUE
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a CYCLE
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a NOCYCLE
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 1
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a NOCACHE
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a ORDER
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a NOORDER
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a RESTART
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your OceanBase version for the right syntax to use
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE + 1
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE - 1
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE + 1.5
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE - 1.5
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 1.5
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS a
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH 5 INCREMENT BY 2
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY -1
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE -10 INCREMENT BY -2
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY -3 START WITH 4
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH 5 INCREMENT BY 2 MINVALUE 1 MAXVALUE 100 CACHE 10 CYCLE ORDER
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 100 CYCLE CACHE 10
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 2 CYCLE MAXVALUE 3
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 100 CYCLE CACHE 200
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 2 CYCLE MAXVALUE 2
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 5 CYCLE
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY -1 CYCLE
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY -1 MINVALUE -100 CYCLE
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 0
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 2
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY 0
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY 0 CACHE 1
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 1 MAXVALUE 1
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH 0
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH -5
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH 100 MAXVALUE 10
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE 10 MAXVALUE 5
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE 9223372036854775807
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE 1 MAXVALUE 3 INCREMENT BY 5
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 9223372036854775808
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 99999999999999999999999999999
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE -99999999999999999999999999999 INCREMENT BY -1
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 10.0
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 99999999999999999999.5
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 1.5 CYCLE
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE 1.5
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY 1.5
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH 1.5
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 1.5
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 1e5
< error
| 5114 42000 Invalid numeric
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a NOMAXVALUE INCREMENT BY -1
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a INCREMENT BY 1 INCREMENT BY 2
< error
| 4312 HY000 duplicate INCREMENT BY specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a START WITH 1 START WITH 2
< error
| 4313 HY000 duplicate START WITH specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 5 MAXVALUE 6
< error
| 4301 HY000 duplicate MAXVALUE/NOMAXVALUE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE 1 MINVALUE 2
< error
| 4302 HY000 duplicate MINVALUE/NOMINVALUE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a NOMINVALUE NOMINVALUE
< error
| 4302 HY000 duplicate MINVALUE/NOMINVALUE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 2 CACHE 3
< error
| 4304 HY000 duplicate CACHE/NOCACHE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a NOCACHE NOCACHE
< error
| 4304 HY000 duplicate CACHE/NOCACHE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a CYCLE CYCLE
< error
| 4303 HY000 duplicate CYCLE/NOCYCLE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a ORDER ORDER
< error
| 4305 HY000 duplicate ORDER/NOORDER specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a RESTART RESTART
< error
| 4382 HY000 duplicate RESTART specified
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a MAXVALUE 5 NOMAXVALUE
< error
| 4306 HY000 conflicting MAXVALUE/NOMAXVALUE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a MINVALUE 5 NOMINVALUE
< error
| 4307 HY000 conflicting MINVALUE/NOMINVALUE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a CACHE 5 NOCACHE
< error
| 4309 HY000 conflicting CACHE/NOCACHE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a ORDER NOORDER
< error
| 4310 HY000 conflicting ORDER/NOORDER specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a NOCYCLE CYCLE
< error
| 4308 HY000 conflicting CYCLE/NOCYCLE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE a CYCLE NOCYCLE ORDER NOORDER
< error
| 4310 HY000 conflicting ORDER/NOORDER specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE IF NOT EXISTS a CYCLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 10
: bytes
| -100
: bytes
| 100
: bytes
| -2
: bytes
| 5
: bytes
| Y
: bytes
| Y
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE IF NOT EXISTS a CYCLE NOCYCLE
< error
| 4308 HY000 conflicting CYCLE/NOCYCLE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE nope.a
< error
| 1049 42000 Unknown database 'nope'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE nope.a CYCLE
< error
| 1049 42000 Unknown database 'nope'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE nope.a CYCLE NOCYCLE
< error
| 4308 HY000 conflicting CYCLE/NOCYCLE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE A
< error
| 4297 HY000 name is already used by an existing object
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE b
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 10
: bytes
| -100
: bytes
| 100
: bytes
| -2
: bytes
| 5
: bytes
| Y
: bytes
| Y
< row
: bytes
| b
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 10
: bytes
| -100
: bytes
| 100
: bytes
| -2
: bytes
| 5
: bytes
| Y
: bytes
| Y
< row
: bytes
| bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
: bytes
| 1
: bytes
| 1
: bytes
| 9223372036854775807
: bytes
| 1
: bytes
| 20
: bytes
| N
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| CREATE SEQUENCE bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
< error
| 1059 42000 Identifier name 'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb' is too long
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| ALTER SEQUENCE a NOMAXVALUE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 10
: bytes
| -100
: bytes
| 9223372036854775808
: bytes
| -2
: bytes
| 5
: bytes
| Y
: bytes
| Y
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| ALTER SEQUENCE a NOMINVALUE
< error
| 4329 HY000 descending sequences that CYCLE must specify MINVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| ALTER SEQUENCE a MAXVALUE 1
< error
| 4324 HY000 START WITH cannot be more than MAXVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| ALTER SEQUENCE a MAXVALUE 1000
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 10
: bytes
| -100
: bytes
| 1000
: bytes
| -2
: bytes
| 5
: bytes
| Y
: bytes
| Y
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| ALTER SEQUENCE a MINVALUE 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 10
: bytes
| 1
: bytes
| 100
: bytes
| -2
: bytes
| 5
: bytes
| Y
: bytes
| Y
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| ALTER SEQUENCE a MINVALUE -50
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 10
: bytes
| -50
: bytes
| 100
: bytes
| -2
: bytes
| 5
: bytes
| Y
: bytes
| Y
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| ALTER SEQUENCE a START WITH 1
< error
| 4311 HY000 cannot alter starting sequence number
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| ALTER SEQUENCE a RESTART
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 100
: bytes
| -100
: bytes
| 100
: bytes
| -2
: bytes
| 5
: bytes
| Y
: bytes
| Y
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| ALTER SEQUENCE a RESTART START WITH 5
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 5
: bytes
| -100
: bytes
| 100
: bytes
| -2
: bytes
| 5
: bytes
| Y
: bytes
| Y
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| ALTER SEQUENCE a RESTART START WITH 1000
< error
| 4324 HY000 START WITH cannot be more than MAXVALUE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| ALTER SEQUENCE a INCREMENT BY -1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 10
: bytes
| -100
: bytes
| 100
: bytes
| -1
: bytes
| 5
: bytes
| Y
: bytes
| Y
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| ALTER SEQUENCE a INCREMENT BY 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 10
: bytes
| -100
: bytes
| 100
: bytes
| 1
: bytes
| 5
: bytes
| Y
: bytes
| Y
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| ALTER SEQUENCE a INCREMENT BY 0
< error
| 4318 HY000 INCREMENT must be a nonzero integer
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| ALTER SEQUENCE a INCREMENT BY -1 RESTART
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 100
: bytes
| -100
: bytes
| 100
: bytes
| -1
: bytes
| 5
: bytes
| Y
: bytes
| Y
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| ALTER SEQUENCE a CYCLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 10
: bytes
| -100
: bytes
| 100
: bytes
| -2
: bytes
| 5
: bytes
| Y
: bytes
| Y
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| ALTER SEQUENCE a NOCYCLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 10
: bytes
| -100
: bytes
| 100
: bytes
| -2
: bytes
| 5
: bytes
| Y
: bytes
| N
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| ALTER SEQUENCE a CACHE 1
< error
| 4326 HY000 the number of values to CACHE must be greater than 1
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| ALTER SEQUENCE a CACHE 2
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 10
: bytes
| -100
: bytes
| 100
: bytes
| -2
: bytes
| 2
: bytes
| Y
: bytes
| Y
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| ALTER SEQUENCE a CACHE 1000
< error
| 4328 HY000 number to CACHE must be less than one cycle
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| ALTER SEQUENCE a NOCACHE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 10
: bytes
| -100
: bytes
| 100
: bytes
| -2
: bytes
| 1
: bytes
| Y
: bytes
| Y
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| ALTER SEQUENCE a ORDER
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 10
: bytes
| -100
: bytes
| 100
: bytes
| -2
: bytes
| 5
: bytes
| Y
: bytes
| Y
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| ALTER SEQUENCE a NOORDER
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 10
: bytes
| -100
: bytes
| 100
: bytes
| -2
: bytes
| 5
: bytes
| N
: bytes
| Y
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| ALTER SEQUENCE a MAXVALUE 1.5
< error
| 4317 HY000 sequence parameter MAXVALUE must be an integer
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| ALTER SEQUENCE a NOMAXVALUE CYCLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 10
: bytes
| -100
: bytes
| 9223372036854775808
: bytes
| -2
: bytes
| 5
: bytes
| Y
: bytes
| Y
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| ALTER SEQUENCE a CYCLE NOCYCLE
< error
| 4308 HY000 conflicting CYCLE/NOCYCLE specifications
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| ALTER SEQUENCE a RESTART RESTART
< error
| 4382 HY000 duplicate RESTART specified
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| ALTER SEQUENCE A NOCACHE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 10
: bytes
| -100
: bytes
| 100
: bytes
| -2
: bytes
| 1
: bytes
| Y
: bytes
| Y
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| ALTER SEQUENCE b NOCACHE
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| ALTER SEQUENCE b
< error
| 4314 HY000 no options specified for ALTER SEQUENCE
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| ALTER SEQUENCE nope.a NOCACHE
< error
| 1049 42000 Unknown database 'nope'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| ALTER SEQUENCE bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
< error
| 1059 42000 Identifier name 'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb' is too long
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| DROP SEQUENCE A
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| DROP SEQUENCE b
< error
| 4298 HY000 sequence does not exist
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS b
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_omni_t_simulation" "VARCHAR" "string"
< column "Table_type" "VARCHAR" "string"
> query
| SELECT SEQUENCE_NAME, LAST_NUMBER, MIN_VALUE, MAX_VALUE, INCREMENT_BY, CACHE_SIZE, ORDER_FLAG, CYCLE_FLAG FROM oceanbase.DBA_SEQUENCES WHERE SEQUENCE_OWNER = DATABASE() ORDER BY SEQUENCE_NAME
< result
< column "SEQUENCE_NAME" "VARCHAR" "string"
< column "LAST_NUMBER" "DECIMAL" "string" precision 38 scale 0
< column "MIN_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "MAX_VALUE" "DECIMAL" "string" precision 28 scale 0
< column "INCREMENT_BY" "DECIMAL" "string" precision 28 scale 0
< column "CACHE_SIZE" "DECIMAL" "string" precision 28 scale 0
< column "ORDER_FLAG" "VARCHAR"
< column "CYCLE_FLAG" "VARCHAR"
< row
: bytes
| a
: bytes
| 10
: bytes
| -100
: bytes
| 100
: bytes
| -2
: bytes
| 5
: bytes
| Y
: bytes
| Y
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| DROP SEQUENCE nope.a
< error
| 1049 42000 Unknown database 'nope'
> exec
| DROP DATABASE IF EXISTS omni_t_simulation; CREATE DATABASE omni_t_simulation; USE omni_t_simulation; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS nope.a
< error
| 1049 42000 Unknown database 'nope'
