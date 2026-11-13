-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestVectorFormsOracle default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| CREATE VECTOR INDEX idx ON t ((VEC_COSINE_DISTANCE(a)))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((VEC_COSINE_DISTANCE(a)))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE VECTOR INDEX idx ON t (a)
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE VECTOR INDEX idx ON t (a)
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE VECTOR INDEX idx ON t (a) USING HNSW
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE VECTOR INDEX idx ON t (a) USING HNSW
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE VECTOR INDEX idx ON t (a, b) USING HNSW
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE VECTOR INDEX idx ON t (a, b) USING HNSW
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((VEC_COSINE_DISTANCE(a))) TYPE BTREE
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((VEC_COSINE_DISTANCE(a))) TYPE BTREE
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE VECTOR INDEX idx USING BTREE ON t ((VEC_COSINE_DISTANCE(a)))
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE VECTOR INDEX idx USING BTREE ON t ((VEC_COSINE_DISTANCE(a)))
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE VECTOR INDEX idx USING HASH ON t ((VEC_COSINE_DISTANCE(a)))
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE VECTOR INDEX idx USING HASH ON t ((VEC_COSINE_DISTANCE(a)))
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((VEC_COSINE_DISTANCE(a))) USING HYPO
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((VEC_COSINE_DISTANCE(a))) USING HYPO
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE VECTOR INDEX idx ON t USING HNSW ((VEC_COSINE_DISTANCE(a)))
< error
| 1064 42000 …line 1 column 34 near "USING HNSW ((VEC_COSINE_DISTANCE(a)))" 
> exec ^
: string
| CREATE VECTOR idx ON t ((VEC_COSINE_DISTANCE(a))) USING HNSW
< error
| 1064 42000 …line 1 column 17 near "idx ON t ((VEC_COSINE_DISTANCE(a))) USING HNSW" 
> exec ^
: string
| CREATE VECTOR KEY idx ON t ((VEC_COSINE_DISTANCE(a))) USING HNSW
< error
| 1064 42000 …line 1 column 17 near "KEY idx ON t ((VEC_COSINE_DISTANCE(a))) USING HNSW" 
> exec ^
: string
| CREATE VECTOR INDEX IF NOT EXISTS idx ON t ((VEC_COSINE_DISTANCE(a))) TYPE HNSW
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE VECTOR INDEX IF NOT EXISTS idx ON t ((VEC_COSINE_DISTANCE(a))) TYPE HNSW
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE VECTOR INDEX IF NOT idx ON t ((a))
< error
| 1064 42000 …line 1 column 30 near "idx ON t ((a))" 
> exec ^
: string
| CREATE VECTOR INDEX IF EXISTS idx ON t ((a))
< error
| 1064 42000 …line 1 column 29 near "EXISTS idx ON t ((a))" 
> exec ^
: string
| CREATE VECTOR INDEX ident TYPE HNSW ON d_n.t_n ((VEC_COSINE_DISTANCE(a)))
< error
| 1146 42S02 Table 'd_n.t_n' doesn't exist
> exec ^
: string
| CREATE VECTOR INDEX ident TYPE HNSW ON d_n.t_n ((VEC_COSINE_DISTANCE(a)))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE VECTOR INDEX idx USING HNSW ON t ((VEC_COSINE_DISTANCE(a)))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE VECTOR INDEX idx USING HNSW ON t ((VEC_COSINE_DISTANCE(a)))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE VECTOR INDEX idx USING `HNSW` ON t ((VEC_COSINE_DISTANCE(a)))
< error
| 1064 42000 …line 1 column 36 near "`HNSW` ON t ((VEC_COSINE_DISTANCE(a)))" 
> exec ^
: string
| CREATE VECTOR INDEX idx USING hnsw ON t ((VEC_COSINE_DISTANCE(a))) TYPE hnsw
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE VECTOR INDEX idx USING hnsw ON t ((VEC_COSINE_DISTANCE(a))) TYPE hnsw
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE VECTOR INDEX idx USING HNSW ON t ((VEC_COSINE_DISTANCE(a))) USING BTREE
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE VECTOR INDEX idx USING HNSW ON t ((VEC_COSINE_DISTANCE(a))) USING BTREE
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE VECTOR INDEX idx USING BTREE ON t ((VEC_COSINE_DISTANCE(a))) USING HNSW
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE VECTOR INDEX idx USING BTREE ON t ((VEC_COSINE_DISTANCE(a))) USING HNSW
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE VECTOR INDEX ON t ((a))
< error
| 1064 42000 …line 1 column 22 near "ON t ((a))" 
> exec ^
: string
| CREATE VECTOR INDEX
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| CREATE VECTOR INDEX idx
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| CREATE VECTOR INDEX idx ON
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| CREATE VECTOR INDEX idx ON t
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| CREATE VECTOR INDEX idx ON t (
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ()
< error
| 1064 42000 …line 1 column 31 near ")" 
> exec ^
: string
| CREATE VECTOR
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| CREATE VECTOR UNIQUE INDEX idx ON t ((a))
< error
| 1064 42000 …line 1 column 20 near "UNIQUE INDEX idx ON t ((a))" 
> exec ^
: string
| CREATE UNIQUE VECTOR INDEX idx ON t ((a))
< error
| 1064 42000 …line 1 column 20 near "VECTOR INDEX idx ON t ((a))" 
> exec ^
: string
| CREATE FULLTEXT VECTOR INDEX idx ON t ((a))
< error
| 1064 42000 …line 1 column 22 near "VECTOR INDEX idx ON t ((a))" 
> exec ^
: string
| CREATE TEMPORARY VECTOR INDEX idx ON t ((a))
< error
| 1064 42000 …line 1 column 23 near "VECTOR INDEX idx ON t ((a))" 
> exec ^
: string
| CREATE GLOBAL TEMPORARY VECTOR INDEX idx ON t ((a))
< error
| 1064 42000 …line 1 column 30 near "VECTOR INDEX idx ON t ((a))" 
> exec ^
: string
| CREATE OR REPLACE VECTOR INDEX idx ON t ((a))
< error
| 1064 42000 …line 1 column 24 near "VECTOR INDEX idx ON t ((a))" 
> exec ^
: string
| CREATE VECTOR INDEX vector ON vector ((vec(vector)))
< error
| 1146 42S02 Table 'test.vector' doesn't exist
> exec ^
: string
| CREATE VECTOR INDEX vector ON vector ((vec(vector)))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE VECTOR INDEX `index` ON t ((a))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE VECTOR INDEX `index` ON t ((a))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE VECTOR INDEX select ON t ((a))
< error
| 1064 42000 …line 1 column 26 near "select ON t ((a))" 
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a) DESC)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a) DESC)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a) ASC)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a) ASC)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE VECTOR INDEX idx ON t (a(3))
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE VECTOR INDEX idx ON t (a(3))
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a) DESC) COMMENT 'c' ALGORITHM = INPLACE LOCK = NONE
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a) DESC) COMMENT 'c' ALGORITHM = INPLACE LOCK = NONE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a)) COMMENT 'c'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a)) COMMENT 'c'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a)) USING HNSW COMMENT 'c' VISIBLE
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a)) USING HNSW COMMENT 'c' VISIBLE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a)) INVISIBLE
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a)) INVISIBLE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a)) GLOBAL
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a)) GLOBAL
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a)) LOCAL
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a)) LOCAL
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a)) CLUSTERED
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a)) CLUSTERED
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a)) KEY_BLOCK_SIZE = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a)) KEY_BLOCK_SIZE = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a)) WITH PARSER p
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a)) WITH PARSER p
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a)) ALGORITHM = DEFAULT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a)) ALGORITHM = DEFAULT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a)) LOCK = SHARED ALGORITHM = COPY
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a)) LOCK = SHARED ALGORITHM = COPY
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a)) ALGORITHM = INPLACE ALGORITHM = COPY
< error
| 1064 42000 …line 1 column 64 near "ALGORITHM = COPY" 
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a)) WITH (distance=l2)
< error
| 1064 42000 …line 1 column 41 near "(distance=l2)" 
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a)) M=6
< error
| 1064 42000 …line 1 column 36 near "M=6" 
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a)) DISTANCE=cosine
< error
| 1064 42000 …line 1 column 43 near "DISTANCE=cosine" 
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a)) ADD_COLUMNAR_REPLICA_ON_DEMAND
< error
| 1064 42000 …line 1 column 65 near "ADD_COLUMNAR_REPLICA_ON_DEMAND" 
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a)),
< error
| 1064 42000 …line 1 column 35 near "," 
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a)) SELECT 1
< error
| 1064 42000 …line 1 column 41 near "SELECT 1" 
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a)) x
< error
| 1064 42000 …line 1 column 36 near "x" 
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a), (b)) USING HNSW
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE VECTOR INDEX idx ON t ((a), (b)) USING HNSW
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE INDEX i USING HNSW ON t (a)
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE INDEX i USING HNSW ON t (a)
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
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
| CREATE UNIQUE INDEX ident USING HNSW ON d_n.t_n ( ident , ident ASC )
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE UNIQUE INDEX ident USING HNSW ON d_n.t_n ( ident , ident ASC )
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE FULLTEXT INDEX i ON t (a) USING HNSW
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE FULLTEXT INDEX i ON t (a) USING HNSW
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE SPATIAL INDEX i ON t (a) USING HNSW
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE SPATIAL INDEX i ON t (a) USING HNSW
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| ALTER TABLE t ADD INDEX (a) USING HNSW
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD INDEX (a) USING HNSW
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a INT, INDEX (a) USING HNSW)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VECTOR)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VECTOR(3))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VECTOR(0))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VECTOR(16383))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VECTOR(16384))
< error
| 1105 HY000 vector cannot have more than 16383 dimensions
> exec ^
: string
| CREATE TABLE t (a VECTOR(16384))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a VECTOR(99999999999999999999))
< error
| 1064 42000 …line 1 column 45 near "99999999999999999999))" 
> exec ^
: string
| CREATE TABLE t (a VECTOR<FLOAT>)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VECTOR<float>)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VECTOR<FLOAT>(3))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VECTOR < FLOAT > ( 3 ))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VECTOR<DOUBLE>)
< error
| 1064 42000 …line 1 column 33 near ")"Only VECTOR is supported for now 
> exec ^
: string
| CREATE TABLE t (a VECTOR<DOUBLE>(3))
< error
| 1064 42000 …line 1 column 36 near ")"Only VECTOR is supported for now 
> exec ^
: string
| CREATE TABLE t (a VECTOR<INT>)
< error
| 1064 42000 …line 1 column 28 near "INT>)" 
> exec ^
: string
| CREATE TABLE t (a VECTOR<ABC>)
< error
| 1064 42000 …line 1 column 28 near "ABC>)" 
> exec ^
: string
| CREATE TABLE t (a VECTOR<>)
< error
| 1064 42000 …line 1 column 26 near "<>)" 
> exec ^
: string
| CREATE TABLE t (a VECTOR<FLOAT)
< error
| 1064 42000 …line 1 column 31 near ")" 
> exec ^
: string
| CREATE TABLE t (a VECTOR FLOAT>)
< error
| 1064 42000 …line 1 column 30 near "FLOAT>)" 
> exec ^
: string
| CREATE TABLE t (a VECTOR<FLOAT>>)
< error
| 1064 42000 …line 1 column 32 near ">>)" 
> exec ^
: string
| CREATE TABLE t (a VECTOR<<FLOAT>)
< error
| 1064 42000 …line 1 column 26 near "<<FLOAT>)" 
> exec ^
: string
| CREATE TABLE t (a VECTOR<`FLOAT`>)
< error
| 1064 42000 …line 1 column 32 near "`FLOAT`>)" 
> exec ^
: string
| CREATE TABLE t (a VECTOR<FLOAT(3)>)
< error
| 1064 42000 …line 1 column 31 near "(3)>)" 
> exec ^
: string
| CREATE TABLE t (a VECTOR(5)<FLOAT>)
< error
| 1064 42000 …line 1 column 28 near "<FLOAT>)" 
> exec ^
: string
| CREATE TABLE t (a VECTOR())
< error
| 1064 42000 …line 1 column 26 near "))" 
> exec ^
: string
| CREATE TABLE t (a VECTOR(3,2))
< error
| 1064 42000 …line 1 column 27 near ",2))" 
> exec ^
: string
| CREATE TABLE t (a VECTOR('3'))
< error
| 1064 42000 …line 1 column 28 near "'3'))" 
> exec ^
: string
| CREATE TABLE t (a VECTOR(b))
< error
| 1064 42000 …line 1 column 26 near "b))" 
> exec ^
: string
| CREATE TABLE t (a VECTOR(1.5))
< error
| 1064 42000 …line 1 column 28 near "1.5))" 
> exec ^
: string
| CREATE TABLE t (a VECTOR(0x3))
< error
| 1064 42000 …line 1 column 28 near "0x3))" 
> exec ^
: string
| CREATE TABLE t (a VECTOR(-1))
< error
| 1064 42000 …line 1 column 26 near "-1))" 
> exec ^
: string
| CREATE TABLE t (a VECTOR(3)
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) NOT NULL)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) NULL)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) COMMENT 'x')
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) DEFAULT NULL)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) DEFAULT '[1,2,3]')
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) UNSIGNED)
< error
| 1064 42000 …line 1 column 36 near "UNSIGNED)" 
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) ZEROFILL)
< error
| 1064 42000 …line 1 column 36 near "ZEROFILL)" 
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) BINARY)
< error
| 1064 42000 …line 1 column 34 near "BINARY)" 
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) CHARSET utf8mb4)
< error
| 1064 42000 …line 1 column 35 near "CHARSET utf8mb4)" 
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) CHARACTER SET utf8mb4)
< error
| 1064 42000 …line 1 column 37 near "CHARACTER SET utf8mb4)" 
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) COLLATE utf8mb4_bin)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) PRIMARY KEY)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) ARRAY)
< error
| 1064 42000 …line 1 column 33 near "ARRAY)" 
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) GENERATED ALWAYS AS ('[1]'))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VECTOR AS ('[1]') STORED)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a `VECTOR`)
< error
| 1064 42000 …line 1 column 26 near "`VECTOR`)" 
> exec ^
: string
| CREATE TABLE t (a VEC)
< error
| 1064 42000 …line 1 column 21 near "VEC)" 
> exec ^
: string
| CREATE TABLE t (a int, b VECTOR(3), VECTOR INDEX ((VEC_COSINE_DISTANCE(b))))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a int, b VECTOR(3), VECTOR INDEX i TYPE HNSW ((VEC_COSINE_DISTANCE(b))))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a int, b VECTOR(3), VECTOR INDEX i USING HNSW ((VEC_COSINE_DISTANCE(b))))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a int, b VECTOR(3), VECTOR INDEX TYPE HNSW ((VEC_COSINE_DISTANCE(b))))
< error
| 1064 42000 …line 1 column 58 near "HNSW ((VEC_COSINE_DISTANCE(b))))" 
> exec ^
: string
| CREATE TABLE t (a int, b VECTOR(3), VECTOR INDEX USING HNSW ((VEC_COSINE_DISTANCE(b))))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a int, b VECTOR(3), VECTOR INDEX IF NOT EXISTS i ((VEC_COSINE_DISTANCE(b))))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a int, b VECTOR(3), VECTOR INDEX (b))
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE TABLE t (a int, b VECTOR(3), VECTOR INDEX (b))
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE TABLE t (a int, b VECTOR(3), VECTOR INDEX ((VEC_COSINE_DISTANCE(b))) USING BTREE)
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE TABLE t (a int, b VECTOR(3), VECTOR INDEX ((VEC_COSINE_DISTANCE(b))) USING BTREE)
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE TABLE t (a int, b VECTOR(3), VECTOR KEY ((VEC_COSINE_DISTANCE(b))))
< error
| 1064 42000 …line 1 column 46 near "KEY ((VEC_COSINE_DISTANCE(b))))" 
> exec ^
: string
| CREATE TABLE t (a int, b VECTOR(3), VECTOR (b))
< error
| 1064 42000 …line 1 column 44 near "(b))" 
> exec ^
: string
| CREATE TABLE t (a int, b VECTOR(3), CONSTRAINT c VECTOR INDEX ((VEC_COSINE_DISTANCE(b))))
< error
| 1064 42000 …line 1 column 55 near "VECTOR INDEX ((VEC_COSINE_DISTANCE(b))))" 
> exec ^
: string
| CREATE TABLE t (a int, b VECTOR(3), CONSTRAINT VECTOR INDEX ((VEC_COSINE_DISTANCE(b))))
< ok 0 0
> exec ^
: string
| CREATE TABLE vector (vector vector)
< ok 0 0
> exec ^
: string
| CREATE TABLE vector (vector vector(3), vector2 vector<float>(3), VECTOR INDEX vector ((vec(vector))))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (vector INT)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (vector INT, INDEX (vector))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (vector)
< error
| 1064 42000 …line 1 column 23 near ")" 
> exec ^
: string
| CREATE TABLE t (vector index)
< error
| 1064 42000 …line 1 column 29 near ")" 
> exec ^
: string
| CREATE TABLE t (vector index INT)
< error
| 1064 42000 …line 1 column 32 near "INT)" 
> exec ^
: string
| CREATE TABLE t (`vector` INT)
< ok 0 0
> exec ^
: string
| SELECT 1 vector
< ok 0 0
> exec ^
: string
| CREATE INDEX vector ON t (vector)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX vector ON t (vector)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a int) COMMENT 'x' vector
< error
| 1064 42000 …line 1 column 41 near "vector" 
> exec ^
: string
| DROP INDEX vector ON t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DROP INDEX vector ON t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS VECTOR)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS VECTOR)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS VECTOR(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS VECTOR(3))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS VECTOR<FLOAT>)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS VECTOR<FLOAT>)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS VECTOR<FLOAT>(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS VECTOR<FLOAT>(3))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS VECTOR<DOUBLE>)
< error
| 1064 42000 …line 1 column 32 near ")"Only VECTOR is supported for now 
> exec ^
: string
| SELECT CAST(a AS VECTOR<DOUBLE>(3))
< error
| 1064 42000 …line 1 column 35 near ")"Only VECTOR is supported for now 
> exec ^
: string
| SELECT CAST(a AS VECTOR(3)<FLOAT>)
< error
| 1064 42000 …line 1 column 27 near "<FLOAT>)" 
> exec ^
: string
| SELECT CAST(a AS VECTOR())
< error
| 1064 42000 …line 1 column 25 near "))" 
> exec ^
: string
| SELECT CAST(a AS VECTOR(3) ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS VECTOR(3) ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS VECTOR(3) UNSIGNED)
< error
| 1064 42000 …line 1 column 35 near "UNSIGNED)" 
> exec ^
: string
| SELECT CONVERT(a, VECTOR)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CONVERT(a, VECTOR)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CONVERT(a, VECTOR(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CONVERT(a, VECTOR(3))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CONVERT(a, VECTOR<FLOAT>(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CONVERT(a, VECTOR<FLOAT>(3))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CONVERT(a USING VECTOR)
< error
| 1115 42000 Unknown character set: 'VECTOR'
> exec ^
: string
| SELECT a::VECTOR
< error
| 1064 42000 …line 1 column 9 near "::VECTOR" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN a VECTOR(3)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN a VECTOR(3)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD a VECTOR
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD a VECTOR
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD a VECTOR<FLOAT>
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD a VECTOR<FLOAT>
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD a VECTOR<DOUBLE>
< error
| 1064 42000 …line 1 column 34 near ""Only VECTOR is supported for now 
> exec ^
: string
| ALTER TABLE t ADD COLUMN (a VECTOR(3), b INT)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN (a VECTOR(3), b INT)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD vector vector(3)
< error
| 1064 42000 …line 1 column 31 near "vector(3)" 
> exec ^
: string
| ALTER TABLE t ADD vector INT
< error
| 1064 42000 …line 1 column 28 near "INT" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN vector INT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN vector INT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD vector int, ADD VECTOR INDEX ((a))
< error
| 1064 42000 …line 1 column 28 near "int, ADD VECTOR INDEX ((a))" 
> exec ^
: string
| ALTER TABLE t ADD VECTOR
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| ALTER TABLE t ADD VECTOR (a)
< error
| 1064 42000 …line 1 column 26 near "(a)" 
> exec ^
: string
| ALTER TABLE t ADD VECTOR KEY ((a))
< error
| 1064 42000 …line 1 column 28 near "KEY ((a))" 
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX ((a))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX ((a))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX (a)
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX (a)
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| ALTER TABLE t ADD CONSTRAINT c VECTOR INDEX ((a))
< error
| 1064 42000 …line 1 column 37 near "VECTOR INDEX ((a))" 
> exec ^
: string
| ALTER TABLE t MODIFY a VECTOR(3)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t MODIFY a VECTOR(3)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t MODIFY COLUMN a VECTOR<FLOAT>(3) NOT NULL
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t MODIFY COLUMN a VECTOR<FLOAT>(3) NOT NULL
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t MODIFY a VECTOR<INT>
< error
| 1064 42000 …line 1 column 33 near "INT>" 
> exec ^
: string
| ALTER TABLE t CHANGE a b VECTOR(3)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t CHANGE a b VECTOR(3)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t CHANGE COLUMN a b VECTOR
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t CHANGE COLUMN a b VECTOR
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE FUNCTION f(a VECTOR) RETURNS VECTOR RETURN a
< error
| 1064 42000 …line 1 column 15 near "FUNCTION f(a VECTOR) RETURNS VECTOR RETURN a" 
> exec ^
: string
| CREATE TABLE t2 LIKE t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE TABLE t2 LIKE t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TEMPORARY TABLE tt (a VECTOR(3))
< ok 0 0
> exec ^
: string
| SELECT VEC_COSINE_DISTANCE(a, '[1,2,3]') FROM t ORDER BY VEC_COSINE_DISTANCE(a, '[1,2,3]') LIMIT 3
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT VEC_COSINE_DISTANCE(a, '[1,2,3]') FROM t ORDER BY VEC_COSINE_DISTANCE(a, '[1,2,3]') LIMIT 3
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT VEC_FROM_TEXT('[1]'), VEC_DIMS(a), VEC_L2_NORM(a), VEC_AS_TEXT(a) FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT VEC_FROM_TEXT('[1]'), VEC_DIMS(a), VEC_L2_NORM(a), VEC_AS_TEXT(a) FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t SET TIFLASH REPLICA 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t SET TIFLASH REPLICA 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD `vector` INT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD `vector` INT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD comment INT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD comment INT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN IF NOT EXISTS vector INT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN IF NOT EXISTS vector INT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t MODIFY vector INT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t MODIFY vector INT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t CHANGE vector vector INT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t CHANGE vector vector INT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD (vector INT)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD (vector INT)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD (vector VECTOR(3))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD (vector VECTOR(3))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD (VECTOR INDEX ((a)))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD (VECTOR INDEX ((a)))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN (VECTOR INDEX ((a)))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN (VECTOR INDEX ((a)))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN IF NOT EXISTS (VECTOR INDEX ((a)))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN IF NOT EXISTS (VECTOR INDEX ((a)))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a VECTOR(4294967296))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a VECTOR(4294967296))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a VECTOR(9223372036854775808))
< error
| 1105 HY000 dimensions for type vector must be at least 0
> exec ^
: string
| CREATE TABLE t (a VECTOR(9223372036854775808))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a VECTOR(18446744073709551615))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VECTOR(18446744073709551616))
< error
| 1064 42000 …line 1 column 45 near "18446744073709551616))" 
> exec ^
: string
| CREATE TABLE t (a VECTOR(1,))
< error
| 1064 42000 …line 1 column 27 near ",))" 
> exec ^
: string
| CREATE TABLE t (a VECTOR<FLOAT>())
< error
| 1064 42000 …line 1 column 33 near "))" 
> exec ^
: string
| CREATE TABLE t (a VECTOR<FLOAT> NOT NULL)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VECTOR/* c */<FLOAT>/* d */(3))
< ok 0 0
> exec ^
: string
| SELECT CAST(a AS VECTOR<FLOAT> ARRAY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS VECTOR<FLOAT> ARRAY)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TEMPORARY TABLE t (a int, VECTOR INDEX ((a)))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) CHECK (a IS NOT NULL))
< ok 0 0
> exec ^
: string
| CREATE VIEW vv AS SELECT CAST(a AS VECTOR)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| CREATE VIEW vv AS SELECT CAST(a AS VECTOR)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a VECTOR(3), CONSTRAINT vector UNIQUE (a))
< ok 0 0
