-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOceanBaseIndexScopeGrammarOracle default#1
> exec
| CREATE DATABASE IF NOT EXISTS obscope
< ok 1 0
> exec
| USE obscope
< ok 0 0
> exec
| DROP TABLE IF EXISTS obscope_t
< ok 0 0
> exec
| CREATE TABLE obscope_t (id INT NOT NULL, a INT, c VARCHAR(100), g GEOMETRY NOT NULL SRID 0)
< ok 0 0
> exec
| PREPARE obscope_probe FROM 'CREATE TABLE t (id INT, a INT, PRIMARY KEY (id) GLOBAL)'
< error
| 1064 42000 …near 'GLOBAL)' at line 1
> exec
| PREPARE obscope_probe FROM 'CREATE TABLE t (id INT, a INT, PRIMARY KEY (id) LOCAL)'
< error
| 1064 42000 …near 'LOCAL)' at line 1
> exec
| PREPARE obscope_probe FROM 'CREATE TABLE t (id INT, a INT, PRIMARY KEY (a) GLOBAL) PARTITION BY HASH(id) PARTITIONS 2'
< error
| 1064 42000 …near 'GLOBAL) PARTITION BY HASH(id) PARTITIONS 2' at line 1
> exec
| PREPARE obscope_probe FROM 'CREATE TABLE t (id INT, a INT, PRIMARY KEY (a) LOCAL) PARTITION BY HASH(id) PARTITIONS 2'
< error
| 1064 42000 …near 'LOCAL) PARTITION BY HASH(id) PARTITIONS 2' at line 1
> exec
| PREPARE obscope_probe FROM 'CREATE TABLE t (id INT, a INT, PRIMARY KEY (id, a) GLOBAL) PARTITION BY HASH(id) PARTITIONS 2'
< error
| 1064 42000 …near 'GLOBAL) PARTITION BY HASH(id) PARTITIONS 2' at line 1
> exec
| PREPARE obscope_probe FROM 'CREATE TABLE t (id INT, a INT, CONSTRAINT pkc PRIMARY KEY (id) GLOBAL)'
< error
| 1064 42000 …near 'GLOBAL)' at line 1
> exec
| PREPARE obscope_probe FROM 'CREATE TABLE t (id INT, a INT, CONSTRAINT pkc PRIMARY KEY (id) LOCAL)'
< error
| 1064 42000 …near 'LOCAL)' at line 1
> exec
| PREPARE obscope_probe FROM 'CREATE TABLE t (id INT, a INT, PRIMARY KEY (id) COMMENT ''x'' GLOBAL)'
< error
| 1064 42000 …near 'GLOBAL)' at line 1
> exec
| PREPARE obscope_probe FROM 'CREATE TABLE t (id INT, a INT, PRIMARY KEY (id) BLOCK_SIZE 16384)'
< error
| 1064 42000 …near 'BLOCK_SIZE 16384)' at line 1
> exec
| PREPARE obscope_probe FROM 'CREATE TABLE t (id INT, a INT, PRIMARY KEY (id) USING BTREE)'
< ok 0 0
> exec
| PREPARE obscope_probe FROM 'CREATE TABLE t (id INT, a INT, PRIMARY KEY USING BTREE (id))'
< ok 0 0
> exec
| PREPARE obscope_probe FROM 'CREATE TABLE t (id INT, a INT, PRIMARY KEY (id) COMMENT ''x'')'
< ok 0 0
> exec
| PREPARE obscope_probe FROM 'CREATE TABLE t (id INT, a INT, CONSTRAINT pkc PRIMARY KEY (id) COMMENT ''x'')'
< ok 0 0
> exec
| PREPARE obscope_probe FROM 'CREATE TABLE t (id INT, a INT, PRIMARY KEY (a)) PARTITION BY HASH(id) PARTITIONS 2'
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| PREPARE obscope_probe FROM 'ALTER TABLE obscope_t ADD PRIMARY KEY (id) GLOBAL'
< ok 0 0
> exec
| PREPARE obscope_probe FROM 'ALTER TABLE obscope_t ADD PRIMARY KEY (id) LOCAL'
< ok 0 0
> exec
| PREPARE obscope_probe FROM 'ALTER TABLE obscope_t ADD PRIMARY KEY (id) BLOCK_SIZE 16384'
< ok 0 0
> exec
| PREPARE obscope_probe FROM 'ALTER TABLE obscope_t ADD CONSTRAINT pkc PRIMARY KEY (id) GLOBAL'
< ok 0 0
> exec
| PREPARE obscope_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), FULLTEXT KEY fk (c) LOCAL)'
< ok 0 0
> exec
| PREPARE obscope_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), FULLTEXT KEY fk (c) GLOBAL)'
< error
| 1235 0A000 global fulltext index is not supported
> exec
| PREPARE obscope_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(100), FULLTEXT KEY fk (c) WITH PARSER ngram LOCAL)'
< ok 0 0
> exec
| PREPARE obscope_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0, SPATIAL KEY sk (g) LOCAL)'
< ok 0 0
> exec
| PREPARE obscope_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0, SPATIAL KEY sk (g) GLOBAL)'
< ok 0 0
> exec
| PREPARE obscope_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0, SPATIAL KEY sk (g) GLOBAL) PARTITION BY HASH(id) PARTITIONS 2'
< error
| 1235 0A000 spatial global index not supported
> exec
| PREPARE obscope_probe FROM 'CREATE FULLTEXT INDEX fk ON obscope_t (c) LOCAL'
< ok 0 0
> exec
| PREPARE obscope_probe FROM 'CREATE FULLTEXT INDEX fk ON obscope_t (c) GLOBAL'
< error
| 1235 0A000 global fulltext index is not supported
> exec
| PREPARE obscope_probe FROM 'CREATE SPATIAL INDEX sk ON obscope_t (g) GLOBAL'
< ok 0 0
> exec
| PREPARE obscope_probe FROM 'ALTER TABLE obscope_t ADD FULLTEXT KEY fk (c) GLOBAL'
< error
| 1235 0A000 global fulltext index is not supported
> exec
| PREPARE obscope_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a) LOCAL GLOBAL)'
< ok 0 0
> exec
| PREPARE obscope_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a) GLOBAL LOCAL)'
< ok 0 0
> exec
| PREPARE obscope_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a) GLOBAL GLOBAL)'
< ok 0 0
> exec
| PREPARE obscope_probe FROM 'CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY k (a) LOCAL COMMENT ''x'' GLOBAL)'
< ok 0 0
> exec
| PREPARE obscope_probe FROM 'CREATE INDEX k ON obscope_t (a) LOCAL GLOBAL'
< ok 0 0
> exec
| PREPARE obscope_probe FROM 'ALTER TABLE obscope_t ADD UNIQUE KEY uk (a) GLOBAL LOCAL'
< ok 0 0
> exec
| DROP DATABASE IF EXISTS obscope
< ok 3 0
