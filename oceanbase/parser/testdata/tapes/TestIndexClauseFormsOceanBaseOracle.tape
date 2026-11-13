-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestIndexClauseFormsOceanBaseOracle default#1
> exec
| DROP DATABASE IF EXISTS obidx
< ok 0 0
> exec
| DROP DATABASE IF EXISTS obidx_d
< ok 0 0
> exec
| CREATE DATABASE obidx
< ok 1 0
> exec
| CREATE DATABASE obidx_d
< ok 1 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE TABLE obidx_t (id INT PRIMARY KEY, a INT, b INT, KEY obidx_ia (a))
< ok 0 0
> exec
| CREATE TABLE obidx_d.obidx_t (a INT, b INT)
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) USING BTREE
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) USING HASH
< error
| 1061 42000 Duplicate key name 'obidx_i'
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) USING RTREE
< error
| 1064 42000 …near 'RTREE' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) TYPE BTREE
< error
| 1064 42000 …near 'TYPE BTREE' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) TYPE HASH
< error
| 1064 42000 …near 'TYPE HASH' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) TYPE RTREE
< error
| 1064 42000 …near 'TYPE RTREE' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i USING BTREE ON obidx_t (a)
< error
| 1061 42000 Duplicate key name 'obidx_i'
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i TYPE BTREE ON obidx_t (a)
< error
| 1064 42000 …near 'TYPE BTREE ON obidx_t (a)' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE UNIQUE INDEX obidx_i ON obidx_d.obidx_t (a, b ASC) TYPE BTREE
< error
| 1064 42000 …near 'TYPE BTREE' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE UNIQUE INDEX obidx_i ON obidx_d.obidx_t (a, b ASC) TYPE HASH
< error
| 1064 42000 …near 'TYPE HASH' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) TYPE BTREE COMMENT 'c'
< error
| 1064 42000 …near 'TYPE BTREE COMMENT 'c'' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) COMMENT 'c' TYPE BTREE
< error
| 1064 42000 …near 'TYPE BTREE' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) USING BTREE TYPE HASH
< error
| 1064 42000 …near 'TYPE HASH' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i TYPE BTREE ON obidx_t (a) TYPE HASH
< error
| 1064 42000 …near 'TYPE BTREE ON obidx_t (a) TYPE HASH' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE TABLE obidx_c11 (a INT, KEY obidx_k (a) USING BTREE)
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE TABLE obidx_c12 (a INT, KEY obidx_k USING BTREE (a))
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE TABLE obidx_c13 (a INT, PRIMARY KEY USING BTREE (a))
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| ALTER TABLE obidx_t ADD KEY obidx_k (a) USING HASH
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE TABLE obidx_c1 (a INT, KEY obidx_k (a) TYPE BTREE)
< error
| 1064 42000 …near 'TYPE BTREE)' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE TABLE obidx_c2 (a INT, KEY obidx_k TYPE BTREE (a))
< error
| 1064 42000 …near 'TYPE BTREE (a))' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE TABLE obidx_c3 (a INT, PRIMARY KEY obidx_pk TYPE BTREE (a))
< error
| 1064 42000 …near 'TYPE BTREE (a))' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE TABLE obidx_c4 (a INT, UNIQUE KEY obidx_u TYPE BTREE (a))
< error
| 1064 42000 …near 'TYPE BTREE (a))' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| ALTER TABLE obidx_t ADD KEY obidx_k (a) TYPE BTREE
< error
| 1064 42000 …near 'TYPE BTREE' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| ALTER TABLE obidx_t ADD KEY obidx_k TYPE BTREE (a)
< error
| 1064 42000 …near 'TYPE BTREE (a)' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| ALTER TABLE obidx_t ADD UNIQUE KEY obidx_u (a) TYPE HASH
< error
| 1064 42000 …near 'TYPE HASH' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE TABLE obidx_c5 (a INT, KEY TYPE BTREE (a))
< error
| 1064 42000 …near 'BTREE (a))' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE TABLE obidx_c6 (a INT, PRIMARY KEY TYPE BTREE (a))
< error
| 1064 42000 …near 'BTREE (a))' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE TABLE obidx_c7 (a INT, UNIQUE KEY TYPE BTREE (a))
< error
| 1064 42000 …near 'BTREE (a))' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| ALTER TABLE obidx_t ADD KEY TYPE BTREE (a)
< error
| 1064 42000 …near 'BTREE (a)' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| ALTER TABLE obidx_t ADD PRIMARY KEY TYPE BTREE (a)
< error
| 1064 42000 …near 'BTREE (a)' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) USING nosuchtype
< error
| 1064 42000 …near 'nosuchtype' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) TYPE nosuchtype
< error
| 1064 42000 …near 'TYPE nosuchtype' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) USING `BTREE`
< error
| 1064 42000 …near '`BTREE`' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i USING nosuchtype ON obidx_t (a)
< error
| 1064 42000 …near 'nosuchtype ON obidx_t (a)' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) USING
< error
| 1064 42000 …near 'USING' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) TYPE
< error
| 1064 42000 …near 'TYPE' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) TYPE 1
< error
| 1064 42000 …near 'TYPE 1' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) TYPE 'BTREE'
< error
| 1064 42000 …near 'TYPE 'BTREE'' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE TABLE obidx_c8 (a INT, KEY obidx_k (a) USING nosuchtype)
< error
| 1064 42000 …near 'nosuchtype)' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE TABLE obidx_c9 (a INT, KEY obidx_k USING nosuchtype (a))
< error
| 1064 42000 …near 'nosuchtype (a))' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| ALTER TABLE obidx_t ADD KEY obidx_k (a) USING nosuchtype
< error
| 1064 42000 …near 'nosuchtype' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) USING HNSW
< error
| 1064 42000 …near 'HNSW' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) USING hnsw
< error
| 1064 42000 …near 'hnsw' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) TYPE HNSW
< error
| 1064 42000 …near 'TYPE HNSW' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) USING `HNSW`
< error
| 1064 42000 …near '`HNSW`' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE TABLE obidx_c10 (a INT, KEY obidx_k (a) USING HNSW)
< error
| 1064 42000 …near 'HNSW)' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) ALGORITHM = DEFAULT
< error
| 1064 42000 …near 'ALGORITHM = DEFAULT' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) ALGORITHM DEFAULT
< error
| 1064 42000 …near 'ALGORITHM DEFAULT' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) LOCK = DEFAULT
< error
| 1064 42000 …near 'LOCK = DEFAULT' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) LOCK DEFAULT
< error
| 1064 42000 …near 'LOCK DEFAULT' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) ALGORITHM = INPLACE
< error
| 1064 42000 …near 'ALGORITHM = INPLACE' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) LOCK = NONE
< error
| 1064 42000 …near 'LOCK = NONE' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) LOCK = EXCLUSIVE
< error
| 1064 42000 …near 'LOCK = EXCLUSIVE' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) ALGORITHM = ident
< error
| 1064 42000 …near 'ALGORITHM = ident' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) ALGORITHM = `default`
< error
| 1064 42000 …near 'ALGORITHM = `default`' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) ALGORITHM = DEFAULT LOCK = DEFAULT
< error
| 1064 42000 …near 'ALGORITHM = DEFAULT LOCK = DEFAULT' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) LOCK = DEFAULT ALGORITHM = DEFAULT
< error
| 1064 42000 …near 'LOCK = DEFAULT ALGORITHM = DEFAULT' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) ALGORITHM
< error
| 1064 42000 …near 'ALGORITHM' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) LOCK
< error
| 1064 42000 …near 'LOCK' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) ALGORITHM =
< error
| 1064 42000 …near 'ALGORITHM =' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) ALGORITHM = 1
< error
| 1064 42000 …near 'ALGORITHM = 1' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) ALGORITHM = 'INPLACE'
< error
| 1064 42000 …near 'ALGORITHM = 'INPLACE'' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) ALGORITHM = select
< error
| 1064 42000 …near 'ALGORITHM = select' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) ALGORITHM = INPLACE ALGORITHM = COPY
< error
| 1064 42000 …near 'ALGORITHM = INPLACE ALGORITHM = COPY' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) LOCK = SHARED LOCK = NONE
< error
| 1064 42000 …near 'LOCK = SHARED LOCK = NONE' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a) ALGORITHM = INPLACE LOCK = NONE ALGORITHM = COPY
< error
| 1064 42000 …near 'ALGORITHM = INPLACE LOCK = NONE ALGORITHM = COPY' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| DROP INDEX obidx_ia ON obidx_t
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| DROP INDEX IF EXISTS obidx_ia ON obidx_t
< error
| 1064 42000 …near 'IF EXISTS obidx_ia ON obidx_t' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| DROP INDEX IF EXISTS obidx_ia ON obidx_d.obidx_t
< error
| 1064 42000 …near 'IF EXISTS obidx_ia ON obidx_d.obidx_t' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| DROP INDEX IF NOT EXISTS obidx_ia ON obidx_t
< error
| 1064 42000 …near 'IF NOT EXISTS obidx_ia ON obidx_t' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| DROP INDEX obidx_ia IF EXISTS ON obidx_t
< error
| 1064 42000 …near 'IF EXISTS ON obidx_t' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| DROP INDEX IF EXISTS ON obidx_t
< error
| 1064 42000 …near 'IF EXISTS ON obidx_t' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| DROP INDEX obidx_ia ON obidx_t ALGORITHM = DEFAULT
< error
| 1064 42000 …near 'ALGORITHM = DEFAULT' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| DROP INDEX obidx_ia ON obidx_t ALGORITHM DEFAULT
< error
| 1064 42000 …near 'ALGORITHM DEFAULT' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| DROP INDEX obidx_ia ON obidx_t LOCK = DEFAULT
< error
| 1064 42000 …near 'LOCK = DEFAULT' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| DROP INDEX obidx_ia ON obidx_t LOCK DEFAULT
< error
| 1064 42000 …near 'LOCK DEFAULT' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| DROP INDEX obidx_ia ON obidx_t ALGORITHM = INPLACE
< error
| 1064 42000 …near 'ALGORITHM = INPLACE' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| DROP INDEX obidx_ia ON obidx_t LOCK = SHARED
< error
| 1064 42000 …near 'LOCK = SHARED' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| DROP INDEX obidx_ia ON obidx_t ALGORITHM = DEFAULT LOCK = DEFAULT
< error
| 1064 42000 …near 'ALGORITHM = DEFAULT LOCK = DEFAULT' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| DROP INDEX obidx_ia ON obidx_t LOCK = DEFAULT ALGORITHM = DEFAULT
< error
| 1064 42000 …near 'LOCK = DEFAULT ALGORITHM = DEFAULT' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| DROP INDEX IF EXISTS obidx_ia ON obidx_t ALGORITHM = DEFAULT
< error
| 1064 42000 …near 'IF EXISTS obidx_ia ON obidx_t ALGORITHM = DEFAULT' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| DROP INDEX obidx_ia ON obidx_t ALGORITHM
< error
| 1064 42000 …near 'ALGORITHM' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| DROP INDEX obidx_ia ON obidx_t LOCK
< error
| 1064 42000 …near 'LOCK' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| DROP INDEX obidx_ia ON obidx_t ALGORITHM =
< error
| 1064 42000 …near 'ALGORITHM =' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| DROP INDEX obidx_ia ON obidx_t ALGORITHM = select
< error
| 1064 42000 …near 'ALGORITHM = select' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| DROP INDEX obidx_ia ON obidx_t ALGORITHM = INPLACE ALGORITHM = COPY
< error
| 1064 42000 …near 'ALGORITHM = INPLACE ALGORITHM = COPY' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| DROP INDEX obidx_ia ON obidx_t LOCK = SHARED LOCK = NONE
< error
| 1064 42000 …near 'LOCK = SHARED LOCK = NONE' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| DROP INDEX obidx_ia ON obidx_t, ALGORITHM = INPLACE
< error
| 1064 42000 …near ', ALGORITHM = INPLACE' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE INDEX obidx_i ON obidx_t (a,)
< error
| 1064 42000 …near ')' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| CREATE UNIQUE INDEX obidx_i ON obidx_d.obidx_t (a, b ASC,) TYPE BTREE
< error
| 1064 42000 …near ') TYPE BTREE' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| LOCK TABLES obidx_t READ
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| LOCK TABLE obidx_t READ
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| LOCK = obidx_t READ
< error
| 1064 42000 …near '= obidx_t READ' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec ^
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| UNLOCK TABLE
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| UNLOCK = obidx_t
< error
| 1064 42000 …near '= obidx_t' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| UNLOCK obidx_t
< error
| 1064 42000 …near 'obidx_t' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| LOCK STATS obidx_t
< error
| 1064 42000 …near 'STATS obidx_t' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| UNLOCK STATS obidx_t
< error
| 1064 42000 …near 'STATS obidx_t' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| UNLOCK STATS
< error
| 1064 42000 …near 'STATS' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| LOCK STATS obidx_t PARTITION (p0)
< error
| 1064 42000 …near 'STATS obidx_t PARTITION (p0)' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| USE obidx
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP DATABASE IF EXISTS obidx
< ok 8 0
> exec
| DROP DATABASE IF EXISTS obidx_d
< ok 1 0
