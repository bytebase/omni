-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestIndexClauseFormsOracle default#1
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE INDEX i ON t (a) USING BTREE
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) USING HASH
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) USING RTREE
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) TYPE BTREE
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) TYPE HASH
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) TYPE RTREE
< ok 0 0
> exec ^
: string
| CREATE INDEX i USING BTREE ON t (a)
< ok 0 0
> exec ^
: string
| CREATE INDEX i TYPE BTREE ON t (a)
< ok 0 0
> exec ^
: string
| CREATE UNIQUE INDEX i ON d.t (a, b ASC) TYPE BTREE
< ok 0 0
> exec ^
: string
| CREATE UNIQUE INDEX i ON d.t (a, b ASC) TYPE HASH
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) TYPE BTREE COMMENT 'c'
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) COMMENT 'c' TYPE BTREE
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) USING BTREE TYPE HASH
< ok 0 0
> exec ^
: string
| CREATE INDEX i TYPE BTREE ON t (a) TYPE HASH
< ok 0 0
> exec ^
: string
| CREATE TABLE ct11 (a INT, KEY k (a) USING BTREE)
< ok 0 0
> exec ^
: string
| CREATE TABLE ct12 (a INT, KEY k USING BTREE (a))
< ok 0 0
> exec ^
: string
| CREATE TABLE ct13 (a INT, PRIMARY KEY USING BTREE (a))
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD KEY k (a) USING HASH
< ok 0 0
> exec ^
: string
| CREATE TABLE ct1 (a INT, KEY k (a) TYPE BTREE)
< ok 0 0
> exec ^
: string
| CREATE TABLE ct2 (a INT, KEY k TYPE BTREE (a))
< ok 0 0
> exec ^
: string
| CREATE TABLE ct3 (a INT, PRIMARY KEY pk TYPE BTREE (a))
< ok 0 0
> exec ^
: string
| CREATE TABLE ct4 (a INT, UNIQUE KEY u TYPE BTREE (a))
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD KEY k (a) TYPE BTREE
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD KEY k TYPE BTREE (a)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD UNIQUE KEY u (a) TYPE HASH
< ok 0 0
> exec ^
: string
| CREATE TABLE ct5 (a INT, KEY TYPE BTREE (a))
< error
| 1064 42000 …near 'BTREE (a))' at line 1
> exec ^
: string
| CREATE TABLE ct6 (a INT, PRIMARY KEY TYPE BTREE (a))
< error
| 1064 42000 …near 'BTREE (a))' at line 1
> exec ^
: string
| CREATE TABLE ct7 (a INT, UNIQUE KEY TYPE BTREE (a))
< error
| 1064 42000 …near 'BTREE (a))' at line 1
> exec ^
: string
| ALTER TABLE t ADD KEY TYPE BTREE (a)
< error
| 1064 42000 …near 'BTREE (a)' at line 1
> exec ^
: string
| ALTER TABLE t ADD PRIMARY KEY TYPE BTREE (a)
< error
| 1064 42000 …near 'BTREE (a)' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a) USING nosuchtype
< error
| 1064 42000 …near 'nosuchtype' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a) TYPE nosuchtype
< error
| 1064 42000 …near 'nosuchtype' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a) USING `BTREE`
< error
| 1064 42000 …near '`BTREE`' at line 1
> exec ^
: string
| CREATE INDEX i USING nosuchtype ON t (a)
< error
| 1064 42000 …near 'nosuchtype ON t (a)' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a) USING
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a) TYPE
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a) TYPE 1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a) TYPE 'BTREE'
< error
| 1064 42000 …near ''BTREE'' at line 1
> exec ^
: string
| CREATE TABLE ct8 (a INT, KEY k (a) USING nosuchtype)
< error
| 1064 42000 …near 'nosuchtype)' at line 1
> exec ^
: string
| CREATE TABLE ct9 (a INT, KEY k USING nosuchtype (a))
< error
| 1064 42000 …near 'nosuchtype (a))' at line 1
> exec ^
: string
| ALTER TABLE t ADD KEY k (a) USING nosuchtype
< error
| 1064 42000 …near 'nosuchtype' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a) USING HNSW
< error
| 1064 42000 …near 'HNSW' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a) USING hnsw
< error
| 1064 42000 …near 'hnsw' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a) TYPE HNSW
< error
| 1064 42000 …near 'HNSW' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a) USING `HNSW`
< error
| 1064 42000 …near '`HNSW`' at line 1
> exec ^
: string
| CREATE TABLE ct10 (a INT, KEY k (a) USING HNSW)
< error
| 1064 42000 …near 'HNSW)' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a) ALGORITHM = DEFAULT
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) ALGORITHM DEFAULT
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) LOCK = DEFAULT
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) LOCK DEFAULT
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) ALGORITHM = INPLACE
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) LOCK = NONE
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) LOCK = EXCLUSIVE
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) ALGORITHM = ident
< error
| 1800 HY000 Unknown ALGORITHM 'ident'
> exec ^
: string
| CREATE INDEX i ON t (a) ALGORITHM = ident
| ; SELECT FROM
< error
| 1800 HY000 Unknown ALGORITHM 'ident'
> exec ^
: string
| CREATE INDEX i ON t (a) ALGORITHM = `default`
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) ALGORITHM = DEFAULT LOCK = DEFAULT
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) LOCK = DEFAULT ALGORITHM = DEFAULT
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) ALGORITHM
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a) LOCK
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a) ALGORITHM =
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a) ALGORITHM = 1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a) ALGORITHM = 'INPLACE'
< error
| 1064 42000 …near ''INPLACE'' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a) ALGORITHM = select
< error
| 1064 42000 …near 'select' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a) ALGORITHM = INPLACE ALGORITHM = COPY
< error
| 1064 42000 …near 'ALGORITHM = COPY' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a) LOCK = SHARED LOCK = NONE
< error
| 1064 42000 …near 'LOCK = NONE' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a) ALGORITHM = INPLACE LOCK = NONE ALGORITHM = COPY
< error
| 1064 42000 …near 'ALGORITHM = COPY' at line 1
> exec ^
: string
| DROP INDEX idx ON t
< ok 0 0
> exec ^
: string
| DROP INDEX IF EXISTS idx ON t
< ok 0 0
> exec ^
: string
| DROP INDEX IF EXISTS idx ON d.t
< ok 0 0
> exec ^
: string
| DROP INDEX IF NOT EXISTS idx ON t
< error
| 1064 42000 …near 'NOT EXISTS idx ON t' at line 1
> exec ^
: string
| DROP INDEX idx IF EXISTS ON t
< error
| 1064 42000 …near 'IF EXISTS ON t' at line 1
> exec ^
: string
| DROP INDEX IF EXISTS ON t
< error
| 1064 42000 …near 'ON t' at line 1
> exec ^
: string
| DROP INDEX idx ON t ALGORITHM = DEFAULT
< error
| 1064 42000 …near 'ALGORITHM = DEFAULT' at line 1
> exec ^
: string
| DROP INDEX idx ON t ALGORITHM DEFAULT
< error
| 1064 42000 …near 'ALGORITHM DEFAULT' at line 1
> exec ^
: string
| DROP INDEX idx ON t LOCK = DEFAULT
< error
| 1064 42000 …near 'LOCK = DEFAULT' at line 1
> exec ^
: string
| DROP INDEX idx ON t LOCK DEFAULT
< error
| 1064 42000 …near 'LOCK DEFAULT' at line 1
> exec ^
: string
| DROP INDEX idx ON t ALGORITHM = INPLACE
< error
| 1064 42000 …near 'ALGORITHM = INPLACE' at line 1
> exec ^
: string
| DROP INDEX idx ON t LOCK = SHARED
< error
| 1064 42000 …near 'LOCK = SHARED' at line 1
> exec ^
: string
| DROP INDEX idx ON t ALGORITHM = DEFAULT LOCK = DEFAULT
< error
| 1064 42000 …near 'ALGORITHM = DEFAULT LOCK = DEFAULT' at line 1
> exec ^
: string
| DROP INDEX idx ON t LOCK = DEFAULT ALGORITHM = DEFAULT
< error
| 1064 42000 …near 'LOCK = DEFAULT ALGORITHM = DEFAULT' at line 1
> exec ^
: string
| DROP INDEX IF EXISTS idx ON t ALGORITHM = DEFAULT
< error
| 1064 42000 …near 'ALGORITHM = DEFAULT' at line 1
> exec ^
: string
| DROP INDEX idx ON t ALGORITHM
< error
| 1064 42000 …near 'ALGORITHM' at line 1
> exec ^
: string
| DROP INDEX idx ON t LOCK
< error
| 1064 42000 …near 'LOCK' at line 1
> exec ^
: string
| DROP INDEX idx ON t ALGORITHM =
< error
| 1064 42000 …near 'ALGORITHM =' at line 1
> exec ^
: string
| DROP INDEX idx ON t ALGORITHM = select
< error
| 1064 42000 …near 'ALGORITHM = select' at line 1
> exec ^
: string
| DROP INDEX idx ON t ALGORITHM = INPLACE ALGORITHM = COPY
< error
| 1064 42000 …near 'ALGORITHM = INPLACE ALGORITHM = COPY' at line 1
> exec ^
: string
| DROP INDEX idx ON t LOCK = SHARED LOCK = NONE
< error
| 1064 42000 …near 'LOCK = SHARED LOCK = NONE' at line 1
> exec ^
: string
| DROP INDEX idx ON t, ALGORITHM = INPLACE
< error
| 1064 42000 …near ' ALGORITHM = INPLACE' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a,)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| CREATE UNIQUE INDEX i ON d.t (a, b ASC,) TYPE BTREE
< error
| 1064 42000 …near ') TYPE BTREE' at line 1
> exec ^
: string
| LOCK TABLES t READ
< ok 0 0
> exec ^
: string
| LOCK TABLE t READ
< ok 0 0
> exec ^
: string
| LOCK = t READ
< error
| 1064 42000 …near '= t READ' at line 1
> exec ^
: string
| UNLOCK TABLES
< ok 0 0
> exec ^
: string
| UNLOCK TABLE
< ok 0 0
> exec ^
: string
| UNLOCK nosuch
< error
| 1064 42000 …near 'nosuch' at line 1
> exec ^
: string
| LOCK STATS t
< error
| 1064 42000 …near 'STATS t' at line 1
> exec ^
: string
| UNLOCK STATS t
< error
| 1064 42000 …near 'STATS t' at line 1
> exec ^
: string
| LOCK STATS
< error
| 1064 42000 …near 'STATS' at line 1
> exec ^
: string
| UNLOCK STATS
< error
| 1064 42000 …near 'STATS' at line 1
