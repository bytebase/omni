-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestTableOptionValueTiDBOracle default#1
> exec
| DROP DATABASE IF EXISTS omni_option_value
< ok 0 0
> exec
| CREATE DATABASE omni_option_value
< ok 0 0
> exec
| USE omni_option_value
< ok 0 0
> exec
| CREATE TABLE ta (a INT)
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) KEY_BLOCK_SIZE'
< error
| 1064 42000 …line 1 column 37 near "" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) KEY_BLOCK_SIZE ='
< error
| 1064 42000 …line 1 column 39 near "" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) ENGINE'
< error
| 1064 42000 …line 1 column 29 near "" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) ENGINE ='
< error
| 1064 42000 …line 1 column 31 near "" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) ROW_FORMAT'
< error
| 1064 42000 …line 1 column 33 near "" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) COMPRESSION'
< error
| 1064 42000 …line 1 column 34 near "" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) STATS_PERSISTENT ='
< error
| 1064 42000 …line 1 column 41 near "" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) TABLESPACE'
< error
| 1064 42000 …line 1 column 33 near "" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) DEFAULT CHARSET ='
< error
| 1064 42000 …line 1 column 40 near "" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) COMMENT'
< error
| 1064 42000 …line 1 column 30 near "" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) COMMENT ='
< error
| 1064 42000 …line 1 column 32 near "" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) COMMENT ENGINE = InnoDB'
< error
| 1064 42000 …line 1 column 37 near "ENGINE = InnoDB" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) CONNECTION'
< error
| 1064 42000 …line 1 column 33 near "" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) CONNECTION ='
< error
| 1064 42000 …line 1 column 35 near "" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) PASSWORD'
< error
| 1064 42000 …line 1 column 31 near "" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) PASSWORD ='
< error
| 1064 42000 …line 1 column 33 near "" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta ENGINE ='
< error
| 1064 42000 …line 1 column 23 near "" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta KEY_BLOCK_SIZE'
< error
| 1064 42000 …line 1 column 29 near "" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta COMMENT'
< error
| 1064 42000 …line 1 column 22 near "" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta COMMENT ='
< error
| 1064 42000 …line 1 column 24 near "" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta ALGORITHM ='
< error
| 1064 42000 …line 1 column 26 near "" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta LOCK'
< error
| 1064 42000 …line 1 column 19 near "" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta CONVERT TO CHARACTER SET'
< error
| 1064 42000 …line 1 column 39 near "" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta CONVERT TO CHARACTER SET utf8mb4 COLLATE'
< error
| 1064 42000 …line 1 column 55 near "" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) CHARACTER'
< error
| 1064 42000 …line 1 column 32 near "" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) CHARACTER ENGINE = InnoDB'
< error
| 1064 42000 …line 1 column 39 near "ENGINE = InnoDB" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta COMMENT ''x'' CHARACTER'
< error
| 1064 42000 …line 1 column 36 near "" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta COMMENT ''x'' DISABLE KEYS'
< error
| 1064 42000 …line 1 column 34 near "DISABLE KEYS" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta COMMENT ''x'' ADD COLUMN b3 INT'
< error
| 1064 42000 …line 1 column 30 near "ADD COLUMN b3 INT" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) SHARD_ROW_ID_BITS'
< error
| 1064 42000 …line 1 column 40 near "" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) SHARD_ROW_ID_BITS ='
< error
| 1064 42000 …line 1 column 42 near "" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) PRE_SPLIT_REGIONS ='
< error
| 1064 42000 …line 1 column 42 near "" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) AUTO_ID_CACHE'
< error
| 1064 42000 …line 1 column 36 near "" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta AUTO_ID_CACHE DEFAULT'
< error
| 1064 42000 …line 1 column 36 near "DEFAULT" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta AUTO_ID_CACHE ''s'''
< error
| 1064 42000 …line 1 column 32 near "'s'" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) AUTO_RANDOM_BASE ='
< error
| 1064 42000 …line 1 column 41 near "" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) TTL_ENABLE ='
< error
| 1064 42000 …line 1 column 35 near "" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) PLACEMENT POLICY'
< error
| 1064 42000 …line 1 column 39 near "" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) PLACEMENT POLICY ='
< error
| 1064 42000 …line 1 column 41 near "" 
> exec
| PREPARE parse_probe FROM 'CREATE DATABASE d PLACEMENT POLICY'
< error
| 1064 42000 …line 1 column 34 near "" 
> exec
| PREPARE parse_probe FROM 'CREATE DATABASE d PLACEMENT POLICY ='
< error
| 1064 42000 …line 1 column 36 near "" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 8'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) ENGINE = InnoDB ROW_FORMAT = DYNAMIC'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) COMMENT ''x'''
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) COMMENT = "x" ENGINE = InnoDB'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) CONNECTION ''x'''
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) PASSWORD = ''x'' COMMENT = ''y'''
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta ENGINE = InnoDB'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta COMMENT ''x'''
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta ALGORITHM = INPLACE, LOCK = NONE'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_bin'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) CHARACTER SET utf8mb4'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta ENGINE = InnoDB COMMENT ''x'''
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta ENGINE = InnoDB ROW_FORMAT = DYNAMIC'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta COMMENT ''x'' ENGINE = InnoDB, ALGORITHM = INPLACE'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta ADD COLUMN b2 INT COMMENT ''c'', ENGINE = InnoDB COMMENT ''z'''
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_bin'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta AUTO_INCREMENT = 5 COMMENT ''x'' ENGINE = InnoDB'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) DEFAULT CHARACTER SET = utf8mb4 COLLATE utf8mb4_bin'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) SHARD_ROW_ID_BITS = 4'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) AUTO_ID_CACHE = 1'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta AUTO_ID_CACHE = 1000'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta AUTO_ID_CACHE 1000'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) PLACEMENT POLICY = DEFAULT'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE DATABASE d PLACEMENT POLICY DEFAULT'
< ok 0 0
> exec
| DROP DATABASE IF EXISTS omni_option_value
< ok 0 0
