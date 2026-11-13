-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestTableOptionValueOracle default#1
> exec
| CREATE TABLE ta (a INT)
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) KEY_BLOCK_SIZE'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) KEY_BLOCK_SIZE ='
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) ENGINE'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) ENGINE ='
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) ROW_FORMAT'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) COMPRESSION'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) STATS_PERSISTENT ='
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) TABLESPACE'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) DEFAULT CHARSET ='
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) COMMENT'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) COMMENT ='
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) COMMENT ENGINE = InnoDB'
< error
| 1064 42000 …near 'ENGINE = InnoDB' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) CONNECTION'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) CONNECTION ='
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) PASSWORD'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) PASSWORD ='
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta ENGINE ='
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta KEY_BLOCK_SIZE'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta COMMENT'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta COMMENT ='
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta ALGORITHM ='
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta LOCK'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta CONVERT TO CHARACTER SET'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta CONVERT TO CHARACTER SET utf8mb4 COLLATE'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) CHARACTER'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) CHARACTER ENGINE = InnoDB'
< error
| 1064 42000 …near 'ENGINE = InnoDB' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta COMMENT ''x'' CHARACTER'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta COMMENT ''x'' DISABLE KEYS'
< error
| 1064 42000 …near 'DISABLE KEYS' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta COMMENT ''x'' ADD COLUMN b3 INT'
< error
| 1064 42000 …near 'ADD COLUMN b3 INT' at line 1
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
