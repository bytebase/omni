-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOceanBaseAlterTableSetGrammarOracle default#1
> exec
| CREATE DATABASE IF NOT EXISTS obset
< ok 1 0
> exec
| USE obset
< ok 0 0
> exec
| DROP TABLE IF EXISTS obset_t
< ok 0 0
> exec
| DROP TABLE IF EXISTS obset_p
< ok 0 0
> exec
| CREATE TABLE obset_t (a INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE obset_p (a DATETIME NOT NULL, PRIMARY KEY (a)) DYNAMIC_PARTITION_POLICY = (ENABLE = TRUE, TIME_UNIT = 'DAY') PARTITION BY RANGE COLUMNS(a) (PARTITION p0 VALUES LESS THAN ('2026-01-01'))
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET COMMENT ''x'''
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET COMMENT = ''y'''
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET PCTFREE = 5'
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET PCTFREE = 6 COMMENT ''z'''
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET PCTFREE = 7, COMMENT ''w'''
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET PCTFREE = 8, SET COMMENT ''v'''
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET AUTO_INCREMENT = 100'
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET ENGINE = InnoDB'
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET CHARSET = utf8mb4'
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET DEFAULT CHARSET utf8mb4'
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET CHARACTER SET utf8mb4'
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET DEFAULT CHARACTER SET = gbk COLLATE gbk_bin'
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET ROW_FORMAT = COMPRESSED'
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET COLLATE utf8mb4_bin'
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET BLOCK_SIZE = 32768 REPLICA_NUM = 1'
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET USE_BLOOM_FILTER = TRUE'
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET KEY_BLOCK_SIZE = 8'
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET AUTO_INCREMENT_MODE = ''ORDER'''
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET COMMENT ''x'' ENGINE = InnoDB ROW_FORMAT = DYNAMIC'
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET COMMENT ''x'', ADD COLUMN b INT'
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t ADD COLUMN b INT, SET COMMENT ''x'''
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET ORGANIZATION HEAP'
< error
| 1235 0A000 organization in the alter table statement not supported
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_p SET DYNAMIC_PARTITION_POLICY (ENABLE = FALSE)'
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_p SET DYNAMIC_PARTITION_POLICY = (ENABLE = FALSE)'
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_p DYNAMIC_PARTITION_POLICY (ENABLE = FALSE)'
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET COMMENT ''x'' COMMENT ''y'''
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET COMMENT ''x'', ENGINE = InnoDB, SET ROW_FORMAT = DYNAMIC'
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET COMMENT ''x'', DROP COLUMN a'
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET COMMENT ''x'', ALGORITHM = INPLACE'
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t ALGORITHM = INPLACE, SET COMMENT ''x'''
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET DEFAULT COLLATE = utf8mb4_bin'
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET STATS_PERSISTENT = 0'
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET MAX_ROWS = 10'
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET CHECKSUM = 1'
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET COMPRESSION = ''zstd_1.3.8'''
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET TABLET_SIZE = 134217728'
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET AUTO_INCREMENT 5'
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t set comment ''x'''
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET ORGANIZATION = INDEX'
< error
| 1235 0A000 organization in the alter table statement not supported
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t PCTFREE = 9 COMMENT ''u'''
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t ENGINE = InnoDB COMMENT ''x'''
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t COMMENT ''x'' PCTFREE 5'
< ok 0 0
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET'
< error
| 1064 42000 …near 'SET' at line 1
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET ALGORITHM = INPLACE'
< error
| 1064 42000 …near 'ALGORITHM = INPLACE' at line 1
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET LOCK = NONE'
< error
| 1064 42000 …near 'LOCK = NONE' at line 1
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET ADD COLUMN b INT'
< error
| 1064 42000 …near 'ADD COLUMN b INT' at line 1
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET SET COMMENT ''x'''
< error
| 1064 42000 …near 'SET COMMENT 'x'' at line 1
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET COMMENT ''x'' ADD COLUMN b INT'
< error
| 1064 42000 …near 'ADD COLUMN b INT' at line 1
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET, COMMENT ''x'''
< error
| 1064 42000 …near ', COMMENT 'x'' at line 1
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET CONVERT TO CHARACTER SET utf8mb4'
< error
| 1064 42000 …near 'CONVERT TO CHARACTER SET utf8mb4' at line 1
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET PARTITION BY HASH(a) PARTITIONS 2'
< error
| 1064 42000 …near 'PARTITION BY HASH(a) PARTITIONS 2' at line 1
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t `SET` COMMENT ''x'''
< error
| 1064 42000 …near '`SET` COMMENT 'x'' at line 1
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t COMMENT ''x'' CHARACTER'
< error
| 1064 42000 …near 'CHARACTER' at line 1
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET COMMENT ''x'' CHARACTER'
< error
| 1064 42000 …near 'CHARACTER' at line 1
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t COMMENT ''x'' CHARACTER, ADD COLUMN b INT'
< error
| 1064 42000 …near ', ADD COLUMN b INT' at line 1
> exec
| PREPARE obset_probe FROM 'ALTER TABLE obset_t SET COMMENT ''x'' DEFAULT CHARACTER'
< error
| 1064 42000 …near 'CHARACTER' at line 1
> exec
| PREPARE obset_probe FROM 'CREATE TABLE obset_u (a INT) CHARACTER'
< error
| 1064 42000 …near 'CHARACTER' at line 1
> exec
| PREPARE obset_probe FROM 'CREATE TABLE obset_u (a INT) COMMENT ''x'' CHARACTER'
< error
| 1064 42000 …near 'CHARACTER' at line 1
> exec
| DROP DATABASE IF EXISTS obset
< ok 2 0
