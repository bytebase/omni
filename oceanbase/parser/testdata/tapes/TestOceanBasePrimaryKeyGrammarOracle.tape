-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOceanBasePrimaryKeyGrammarOracle default#1
> exec
| CREATE DATABASE IF NOT EXISTS obpk
< ok 1 0
> exec
| USE obpk
< ok 0 0
> exec
| DROP TABLE IF EXISTS obpk_t
< ok 0 0
> exec
| CREATE TABLE obpk_t (id INT NOT NULL, a INT)
< ok 0 0
> exec
| PREPARE obpk_probe FROM 'CREATE TABLE t (id INT, PRIMARY KEY (id) KEY_BLOCK_SIZE 8)'
< error
| 1064 42000 …near 'KEY_BLOCK_SIZE 8)' at line 1
> exec
| PREPARE obpk_probe FROM 'CREATE TABLE t (id INT, PRIMARY KEY (id) KEY_BLOCK_SIZE = 8)'
< error
| 1064 42000 …near 'KEY_BLOCK_SIZE = 8)' at line 1
> exec
| PREPARE obpk_probe FROM 'CREATE TABLE t (id INT, PRIMARY KEY (id) VISIBLE)'
< error
| 1064 42000 …near 'VISIBLE)' at line 1
> exec
| PREPARE obpk_probe FROM 'CREATE TABLE t (id INT, PRIMARY KEY (id) INVISIBLE)'
< error
| 1064 42000 …near 'INVISIBLE)' at line 1
> exec
| PREPARE obpk_probe FROM 'CREATE TABLE t (id INT, PRIMARY KEY (id) WITH PARSER ngram)'
< error
| 1064 42000 …near 'WITH PARSER ngram)' at line 1
> exec
| PREPARE obpk_probe FROM 'CREATE TABLE t (id INT, PRIMARY KEY (id) COMMENT ''x'' USING BTREE)'
< error
| 1064 42000 …near 'USING BTREE)' at line 1
> exec
| PREPARE obpk_probe FROM 'CREATE TABLE t (id INT, PRIMARY KEY (id) USING BTREE USING HASH)'
< error
| 1064 42000 …near 'USING HASH)' at line 1
> exec
| PREPARE obpk_probe FROM 'CREATE TABLE t (id INT, PRIMARY KEY (id) COMMENT ''x'' COMMENT ''y'')'
< error
| 1064 42000 …near 'COMMENT 'y')' at line 1
> exec
| PREPARE obpk_probe FROM 'CREATE TABLE t (id INT, PRIMARY KEY (id) COMMENT = ''x'')'
< error
| 1064 42000 …near '= 'x')' at line 1
> exec
| PREPARE obpk_probe FROM 'CREATE TABLE t (id INT, PRIMARY KEY (id) ENGINE_ATTRIBUTE=''{}'')'
< error
| 1064 42000 …near 'ENGINE_ATTRIBUTE='{}')' at line 1
> exec
| PREPARE obpk_probe FROM 'CREATE TABLE t (id INT, PRIMARY KEY (id) SECONDARY_ENGINE_ATTRIBUTE=''{}'')'
< error
| 1064 42000 …near 'SECONDARY_ENGINE_ATTRIBUTE='{}')' at line 1
> exec
| PREPARE obpk_probe FROM 'CREATE TABLE t (id INT, PRIMARY KEY (id) VISIBLE COMMENT ''x'')'
< error
| 1064 42000 …near 'VISIBLE COMMENT 'x')' at line 1
> exec
| PREPARE obpk_probe FROM 'CREATE TABLE t (id INT, CONSTRAINT c PRIMARY KEY (id) INVISIBLE)'
< error
| 1064 42000 …near 'INVISIBLE)' at line 1
> exec
| PREPARE obpk_probe FROM 'CREATE TABLE t (id INT, CONSTRAINT c PRIMARY KEY (id) KEY_BLOCK_SIZE 8)'
< error
| 1064 42000 …near 'KEY_BLOCK_SIZE 8)' at line 1
> exec
| PREPARE obpk_probe FROM 'CREATE TABLE t (id INT, CONSTRAINT c PRIMARY KEY (id) COMMENT ''x'' USING HASH)'
< error
| 1064 42000 …near 'USING HASH)' at line 1
> exec
| PREPARE obpk_probe FROM 'CREATE TABLE t (id INT, PRIMARY KEY (id) USING BTREE COMMENT ''x'')'
< ok 0 0
> exec
| PREPARE obpk_probe FROM 'CREATE TABLE t (id INT, PRIMARY KEY USING BTREE (id) USING HASH)'
< ok 0 0
> exec
| PREPARE obpk_probe FROM 'CREATE TABLE t (id INT, PRIMARY KEY USING BTREE (id) USING HASH COMMENT ''x'')'
< ok 0 0
> exec
| PREPARE obpk_probe FROM 'CREATE TABLE t (id INT, PRIMARY KEY USING HASH (id) COMMENT ''x'')'
< ok 0 0
> exec
| PREPARE obpk_probe FROM 'CREATE TABLE t (id INT, CONSTRAINT `c` PRIMARY KEY (id) USING BTREE COMMENT ''x'')'
< ok 0 0
> exec
| PREPARE obpk_probe FROM 'CREATE TABLE t (id INT, CONSTRAINT c PRIMARY KEY pk (id))'
< error
| 1064 42000 …near 'pk (id))' at line 1
> exec
| PREPARE obpk_probe FROM 'CREATE TABLE t (id INT, PRIMARY (id))'
< error
| 1064 42000 …near '(id))' at line 1
> exec
| PREPARE obpk_probe FROM 'CREATE TABLE t (id INT, PRIMARY KEY pk (id))'
< ok 0 0
> exec
| PREPARE obpk_probe FROM 'CREATE TABLE t (id INT, PRIMARY KEY `pk` (id))'
< ok 0 0
> exec
| PREPARE obpk_probe FROM 'CREATE TABLE t (id INT, PRIMARY KEY pk USING BTREE (id))'
< ok 0 0
> exec
| PREPARE obpk_probe FROM 'CREATE TABLE t (id INT, CONSTRAINT c PRIMARY KEY USING BTREE (id))'
< ok 0 0
> exec
| PREPARE obpk_probe FROM 'CREATE TABLE t (id INT, CONSTRAINT PRIMARY KEY (id))'
< ok 0 0
> exec
| PREPARE obpk_probe FROM 'CREATE TABLE t (id INT, PRIMARY KEY (id(3)))'
< error
| 1064 42000 …near '(3)))' at line 1
> exec
| PREPARE obpk_probe FROM 'CREATE TABLE t (id INT, PRIMARY KEY (id DESC))'
< error
| 1064 42000 …near 'DESC))' at line 1
> exec
| PREPARE obpk_probe FROM 'CREATE TABLE t (id INT, PRIMARY KEY (id ASC))'
< error
| 1064 42000 …near 'ASC))' at line 1
> exec
| PREPARE obpk_probe FROM 'CREATE TABLE t (id INT, PRIMARY KEY ((id + 1)))'
< error
| 1064 42000 …near '(id + 1)))' at line 1
> exec
| PREPARE obpk_probe FROM 'CREATE TABLE t (id INT, a INT, PRIMARY KEY (id, a))'
< ok 0 0
> exec
| PREPARE obpk_probe FROM 'CREATE TABLE t (id INT, PRIMARY KEY (`id`))'
< ok 0 0
> exec
| PREPARE obpk_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY COMMENT ''x'')'
< ok 0 0
> exec
| PREPARE obpk_probe FROM 'ALTER TABLE obpk_t ADD PRIMARY KEY USING BTREE (id)'
< error
| 1064 42000 …near 'USING BTREE (id)' at line 1
> exec
| PREPARE obpk_probe FROM 'ALTER TABLE obpk_t ADD CONSTRAINT c PRIMARY KEY USING BTREE (id)'
< error
| 1064 42000 …near 'USING BTREE (id)' at line 1
> exec
| PREPARE obpk_probe FROM 'ALTER TABLE obpk_t ADD CONSTRAINT c PRIMARY KEY pk (id)'
< error
| 1064 42000 …near 'pk (id)' at line 1
> exec
| PREPARE obpk_probe FROM 'ALTER TABLE obpk_t ADD PRIMARY (id)'
< error
| 1064 42000 …near '(id)' at line 1
> exec
| PREPARE obpk_probe FROM 'ALTER TABLE obpk_t ADD PRIMARY KEY (id(3))'
< error
| 1064 42000 …near '(3))' at line 1
> exec
| PREPARE obpk_probe FROM 'ALTER TABLE obpk_t ADD PRIMARY KEY (id DESC)'
< error
| 1064 42000 …near 'DESC)' at line 1
> exec
| PREPARE obpk_probe FROM 'ALTER TABLE obpk_t ADD PRIMARY KEY (id ASC)'
< error
| 1064 42000 …near 'ASC)' at line 1
> exec
| PREPARE obpk_probe FROM 'ALTER TABLE obpk_t ADD PRIMARY KEY ((id + 1))'
< error
| 1064 42000 …near '(id + 1))' at line 1
> exec
| PREPARE obpk_probe FROM 'ALTER TABLE obpk_t ADD PRIMARY KEY pk (id)'
< ok 0 0
> exec
| PREPARE obpk_probe FROM 'ALTER TABLE obpk_t ADD PRIMARY KEY (id) KEY_BLOCK_SIZE 8'
< ok 0 0
> exec
| PREPARE obpk_probe FROM 'ALTER TABLE obpk_t ADD PRIMARY KEY (id) VISIBLE'
< ok 0 0
> exec
| PREPARE obpk_probe FROM 'ALTER TABLE obpk_t ADD PRIMARY KEY (id) INVISIBLE'
< ok 0 0
> exec
| PREPARE obpk_probe FROM 'ALTER TABLE obpk_t ADD PRIMARY KEY (id) COMMENT ''x'' USING BTREE'
< ok 0 0
> exec
| PREPARE obpk_probe FROM 'ALTER TABLE obpk_t ADD PRIMARY KEY (id) USING BTREE USING HASH'
< ok 0 0
> exec
| PREPARE obpk_probe FROM 'ALTER TABLE obpk_t ADD CONSTRAINT c PRIMARY KEY (id) INVISIBLE'
< ok 0 0
> exec
| DROP DATABASE IF EXISTS obpk
< ok 1 0
