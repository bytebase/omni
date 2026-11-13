-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestIndexClauseFormsOracle default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| CREATE INDEX i ON t (a) USING BTREE
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i ON t (a) USING BTREE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX i ON t (a) USING HASH
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i ON t (a) USING HASH
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX i ON t (a) USING RTREE
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i ON t (a) USING RTREE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX i ON t (a) TYPE BTREE
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i ON t (a) TYPE BTREE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX i ON t (a) TYPE HASH
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i ON t (a) TYPE HASH
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX i ON t (a) TYPE RTREE
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i ON t (a) TYPE RTREE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX i USING BTREE ON t (a)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i USING BTREE ON t (a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX i TYPE BTREE ON t (a)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i TYPE BTREE ON t (a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE UNIQUE INDEX i ON d.t (a, b ASC) TYPE BTREE
< error
| 1146 42S02 Table 'd.t' doesn't exist
> exec ^
: string
| CREATE UNIQUE INDEX i ON d.t (a, b ASC) TYPE BTREE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE UNIQUE INDEX i ON d.t (a, b ASC) TYPE HASH
< error
| 1146 42S02 Table 'd.t' doesn't exist
> exec ^
: string
| CREATE UNIQUE INDEX i ON d.t (a, b ASC) TYPE HASH
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX i ON t (a) TYPE BTREE COMMENT 'c'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i ON t (a) TYPE BTREE COMMENT 'c'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX i ON t (a) COMMENT 'c' TYPE BTREE
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i ON t (a) COMMENT 'c' TYPE BTREE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX i ON t (a) USING BTREE TYPE HASH
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i ON t (a) USING BTREE TYPE HASH
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX i TYPE BTREE ON t (a) TYPE HASH
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i TYPE BTREE ON t (a) TYPE HASH
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
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
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD KEY k (a) USING HASH
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
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
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD KEY k (a) TYPE BTREE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD KEY k TYPE BTREE (a)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD KEY k TYPE BTREE (a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD UNIQUE KEY u (a) TYPE HASH
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD UNIQUE KEY u (a) TYPE HASH
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE ct5 (a INT, KEY TYPE BTREE (a))
< error
| 1064 42000 …line 1 column 39 near "BTREE (a))" 
> exec ^
: string
| CREATE TABLE ct6 (a INT, PRIMARY KEY TYPE BTREE (a))
< error
| 1064 42000 …line 1 column 47 near "BTREE (a))" 
> exec ^
: string
| CREATE TABLE ct7 (a INT, UNIQUE KEY TYPE BTREE (a))
< error
| 1064 42000 …line 1 column 46 near "BTREE (a))" 
> exec ^
: string
| ALTER TABLE t ADD KEY TYPE BTREE (a)
< error
| 1064 42000 …line 1 column 32 near "BTREE (a)" 
> exec ^
: string
| ALTER TABLE t ADD PRIMARY KEY TYPE BTREE (a)
< error
| 1064 42000 …line 1 column 40 near "BTREE (a)" 
> exec ^
: string
| CREATE INDEX i ON t (a) USING nosuchtype
< error
| 1064 42000 …line 1 column 40 near "nosuchtype" 
> exec ^
: string
| CREATE INDEX i ON t (a) TYPE nosuchtype
< error
| 1064 42000 …line 1 column 39 near "nosuchtype" 
> exec ^
: string
| CREATE INDEX i ON t (a) USING `BTREE`
< error
| 1064 42000 …line 1 column 37 near "`BTREE`" 
> exec ^
: string
| CREATE INDEX i USING nosuchtype ON t (a)
< error
| 1064 42000 …line 1 column 31 near "nosuchtype ON t (a)" 
> exec ^
: string
| CREATE INDEX i ON t (a) USING
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| CREATE INDEX i ON t (a) TYPE
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| CREATE INDEX i ON t (a) TYPE 1
< error
| 1064 42000 …line 1 column 30 near "1" 
> exec ^
: string
| CREATE INDEX i ON t (a) TYPE 'BTREE'
< error
| 1064 42000 …line 1 column 36 near "'BTREE'" 
> exec ^
: string
| CREATE TABLE ct8 (a INT, KEY k (a) USING nosuchtype)
< error
| 1064 42000 …line 1 column 51 near "nosuchtype)" 
> exec ^
: string
| CREATE TABLE ct9 (a INT, KEY k USING nosuchtype (a))
< error
| 1064 42000 …line 1 column 47 near "nosuchtype (a))" 
> exec ^
: string
| ALTER TABLE t ADD KEY k (a) USING nosuchtype
< error
| 1064 42000 …line 1 column 44 near "nosuchtype" 
> exec ^
: string
| CREATE INDEX i ON t (a) USING HNSW
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE INDEX i ON t (a) USING HNSW
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE INDEX i ON t (a) USING hnsw
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE INDEX i ON t (a) USING hnsw
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE INDEX i ON t (a) TYPE HNSW
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE INDEX i ON t (a) TYPE HNSW
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE INDEX i ON t (a) USING `HNSW`
< error
| 1064 42000 …line 1 column 36 near "`HNSW`" 
> exec ^
: string
| CREATE TABLE ct10 (a INT, KEY k (a) USING HNSW)
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) ALGORITHM = DEFAULT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i ON t (a) ALGORITHM = DEFAULT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX i ON t (a) ALGORITHM DEFAULT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i ON t (a) ALGORITHM DEFAULT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX i ON t (a) LOCK = DEFAULT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i ON t (a) LOCK = DEFAULT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX i ON t (a) LOCK DEFAULT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i ON t (a) LOCK DEFAULT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX i ON t (a) ALGORITHM = INPLACE
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i ON t (a) ALGORITHM = INPLACE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX i ON t (a) LOCK = NONE
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i ON t (a) LOCK = NONE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX i ON t (a) LOCK = EXCLUSIVE
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i ON t (a) LOCK = EXCLUSIVE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX i ON t (a) ALGORITHM = ident
< error
| 1800 HY000 Unknown ALGORITHM 'ALGORITHM'
> exec ^
: string
| CREATE INDEX i ON t (a) ALGORITHM = ident
| ; SELECT FROM
< error
| 1800 HY000 Unknown ALGORITHM 'ALGORITHM'
> exec ^
: string
| CREATE INDEX i ON t (a) ALGORITHM = `default`
< error
| 1800 HY000 Unknown ALGORITHM 'ALGORITHM'
> exec ^
: string
| CREATE INDEX i ON t (a) ALGORITHM = `default`
| ; SELECT FROM
< error
| 1800 HY000 Unknown ALGORITHM 'ALGORITHM'
> exec ^
: string
| CREATE INDEX i ON t (a) ALGORITHM = DEFAULT LOCK = DEFAULT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i ON t (a) ALGORITHM = DEFAULT LOCK = DEFAULT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX i ON t (a) LOCK = DEFAULT ALGORITHM = DEFAULT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i ON t (a) LOCK = DEFAULT ALGORITHM = DEFAULT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX i ON t (a) ALGORITHM
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| CREATE INDEX i ON t (a) LOCK
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| CREATE INDEX i ON t (a) ALGORITHM =
< error
| 1064 42000 …line 1 column 35 near "" 
> exec ^
: string
| CREATE INDEX i ON t (a) ALGORITHM = 1
< error
| 1064 42000 …line 1 column 37 near "1" 
> exec ^
: string
| CREATE INDEX i ON t (a) ALGORITHM = 'INPLACE'
< error
| 1064 42000 …line 1 column 45 near "'INPLACE'" 
> exec ^
: string
| CREATE INDEX i ON t (a) ALGORITHM = select
< error
| 1064 42000 …line 1 column 42 near "select" 
> exec ^
: string
| CREATE INDEX i ON t (a) ALGORITHM = INPLACE ALGORITHM = COPY
< error
| 1064 42000 …line 1 column 53 near "ALGORITHM = COPY" 
> exec ^
: string
| CREATE INDEX i ON t (a) LOCK = SHARED LOCK = NONE
< error
| 1064 42000 …line 1 column 42 near "LOCK = NONE" 
> exec ^
: string
| CREATE INDEX i ON t (a) ALGORITHM = INPLACE LOCK = NONE ALGORITHM = COPY
< error
| 1064 42000 …line 1 column 65 near "ALGORITHM = COPY" 
> exec ^
: string
| DROP INDEX idx ON t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DROP INDEX idx ON t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DROP INDEX IF EXISTS idx ON t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DROP INDEX IF EXISTS idx ON t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DROP INDEX IF EXISTS idx ON d.t
< error
| 1146 42S02 Table 'd.t' doesn't exist
> exec ^
: string
| DROP INDEX IF EXISTS idx ON d.t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DROP INDEX IF NOT EXISTS idx ON t
< error
| 1064 42000 …line 1 column 17 near "NOT EXISTS idx ON t" 
> exec ^
: string
| DROP INDEX idx IF EXISTS ON t
< error
| 1064 42000 …line 1 column 17 near "IF EXISTS ON t" 
> exec ^
: string
| DROP INDEX IF EXISTS ON t
< error
| 1064 42000 …line 1 column 23 near "ON t" 
> exec ^
: string
| DROP INDEX idx ON t ALGORITHM = DEFAULT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DROP INDEX idx ON t ALGORITHM = DEFAULT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DROP INDEX idx ON t ALGORITHM DEFAULT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DROP INDEX idx ON t ALGORITHM DEFAULT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DROP INDEX idx ON t LOCK = DEFAULT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DROP INDEX idx ON t LOCK = DEFAULT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DROP INDEX idx ON t LOCK DEFAULT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DROP INDEX idx ON t LOCK DEFAULT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DROP INDEX idx ON t ALGORITHM = INPLACE
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DROP INDEX idx ON t ALGORITHM = INPLACE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DROP INDEX idx ON t LOCK = SHARED
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DROP INDEX idx ON t LOCK = SHARED
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DROP INDEX idx ON t ALGORITHM = DEFAULT LOCK = DEFAULT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DROP INDEX idx ON t ALGORITHM = DEFAULT LOCK = DEFAULT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DROP INDEX idx ON t LOCK = DEFAULT ALGORITHM = DEFAULT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DROP INDEX idx ON t LOCK = DEFAULT ALGORITHM = DEFAULT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DROP INDEX IF EXISTS idx ON t ALGORITHM = DEFAULT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DROP INDEX IF EXISTS idx ON t ALGORITHM = DEFAULT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DROP INDEX idx ON t ALGORITHM
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| DROP INDEX idx ON t LOCK
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| DROP INDEX idx ON t ALGORITHM =
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| DROP INDEX idx ON t ALGORITHM = select
< error
| 1064 42000 …line 1 column 38 near "select" 
> exec ^
: string
| DROP INDEX idx ON t ALGORITHM = INPLACE ALGORITHM = COPY
< error
| 1064 42000 …line 1 column 49 near "ALGORITHM = COPY" 
> exec ^
: string
| DROP INDEX idx ON t LOCK = SHARED LOCK = NONE
< error
| 1064 42000 …line 1 column 38 near "LOCK = NONE" 
> exec ^
: string
| DROP INDEX idx ON t, ALGORITHM = INPLACE
< error
| 1064 42000 …line 1 column 20 near ", ALGORITHM = INPLACE" 
> exec ^
: string
| CREATE INDEX i ON t (a,)
< error
| 1064 42000 …line 1 column 24 near ")" 
> exec ^
: string
| CREATE UNIQUE INDEX i ON d.t (a, b ASC,) TYPE BTREE
< error
| 1064 42000 …line 1 column 40 near ") TYPE BTREE" 
> exec ^
: string
| LOCK TABLES t READ
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| LOCK TABLES t READ
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| LOCK TABLE t READ
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| LOCK TABLE t READ
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| LOCK = t READ
< error
| 1064 42000 …line 1 column 6 near "= t READ" 
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
| 1064 42000 …line 1 column 13 near "nosuch" 
> exec ^
: string
| LOCK STATS t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| LOCK STATS t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| UNLOCK STATS t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| UNLOCK STATS t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| LOCK STATS
< error
| 1064 42000 …line 1 column 10 near "" 
> exec ^
: string
| UNLOCK STATS
< error
| 1064 42000 …line 1 column 12 near "" 
