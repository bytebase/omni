-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestVectorFormsOracle default#1
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a VECTOR)
< error
| 1064 42000 …near 'VECTOR)' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(3))
< error
| 1064 42000 …near 'VECTOR(3))' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(0))
< error
| 1064 42000 …near 'VECTOR(0))' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(16383))
< error
| 1064 42000 …near 'VECTOR(16383))' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(16384))
< error
| 1064 42000 …near 'VECTOR(16384))' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(99999999999999999999))
< error
| 1064 42000 …near 'VECTOR(99999999999999999999))' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR<FLOAT>)
< error
| 1064 42000 …near 'VECTOR<FLOAT>)' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR<FLOAT>(3))
< error
| 1064 42000 …near 'VECTOR<FLOAT>(3))' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(5)<FLOAT>)
< error
| 1064 42000 …near 'VECTOR(5)<FLOAT>)' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR())
< error
| 1064 42000 …near 'VECTOR())' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(3,2))
< error
| 1064 42000 …near 'VECTOR(3,2))' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR('3'))
< error
| 1064 42000 …near 'VECTOR('3'))' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(b))
< error
| 1064 42000 …near 'VECTOR(b))' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(1.5))
< error
| 1064 42000 …near 'VECTOR(1.5))' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(0x3))
< error
| 1064 42000 …near 'VECTOR(0x3))' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(-1))
< error
| 1064 42000 …near 'VECTOR(-1))' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(3)
< error
| 1064 42000 …near 'VECTOR(3)' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(3, FLOAT))
< error
| 1064 42000 …near 'VECTOR(3, FLOAT))' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) NOT NULL)
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL)' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) NULL)
< error
| 1064 42000 …near 'VECTOR(3) NULL)' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) COMMENT 'x')
< error
| 1064 42000 …near 'VECTOR(3) COMMENT 'x')' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) DEFAULT NULL)
< error
| 1064 42000 …near 'VECTOR(3) DEFAULT NULL)' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) DEFAULT '[1,2,3]')
< error
| 1064 42000 …near 'VECTOR(3) DEFAULT '[1,2,3]')' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) UNSIGNED)
< error
| 1064 42000 …near 'VECTOR(3) UNSIGNED)' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) ZEROFILL)
< error
| 1064 42000 …near 'VECTOR(3) ZEROFILL)' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) BINARY)
< error
| 1064 42000 …near 'VECTOR(3) BINARY)' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) CHARSET utf8mb4)
< error
| 1064 42000 …near 'VECTOR(3) CHARSET utf8mb4)' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) CHARACTER SET utf8mb4)
< error
| 1064 42000 …near 'VECTOR(3) CHARACTER SET utf8mb4)' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) COLLATE utf8mb4_bin)
< error
| 1064 42000 …near 'VECTOR(3) COLLATE utf8mb4_bin)' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) PRIMARY KEY)
< error
| 1064 42000 …near 'VECTOR(3) PRIMARY KEY)' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) ARRAY)
< error
| 1064 42000 …near 'VECTOR(3) ARRAY)' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) GENERATED ALWAYS AS ('[1]'))
< error
| 1064 42000 …near 'VECTOR(3) GENERATED ALWAYS AS ('[1]'))' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) AS ('[1]') STORED)
< error
| 1064 42000 …near 'VECTOR(3) AS ('[1]') STORED)' at line 1
> exec ^
: string
| CREATE TABLE t (a `VECTOR`)
< error
| 1064 42000 …near '`VECTOR`)' at line 1
> exec ^
: string
| CREATE TABLE t (a `VECTOR`(3))
< error
| 1064 42000 …near '`VECTOR`(3))' at line 1
> exec ^
: string
| CREATE TABLE t (a vector(3))
< error
| 1064 42000 …near 'vector(3))' at line 1
> exec ^
: string
| CREATE TEMPORARY TABLE tt (a VECTOR(3))
< error
| 1064 42000 …near 'VECTOR(3))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX (v))
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX (v))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR KEY (v))
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR KEY (v))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR (v))
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR (v))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR i (v))
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR i (v))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v))
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX i (v))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR KEY i (v))
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR KEY i (v))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX IF NOT EXISTS i (v))
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX IF NOT EXISTS i (v))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i USING HNSW (v))
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX i USING HNSW (v))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) USING HNSW)
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX i (v) USING HNSW)' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) USING BTREE)
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX i (v) USING BTREE)' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) TYPE HNSW)
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX i (v) TYPE HNSW)' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX ((VEC_COSINE_DISTANCE(v))))
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX ((VEC_COSINE_DISTANCE(v))))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX (v, w))
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX (v, w))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX (v(3)))
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX (v(3)))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX (v DESC))
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX (v DESC))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M=6)
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX i (v) M=6)' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M 6)
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX i (v) M 6)' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M=abc)
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX i (v) M=abc)' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M='6')
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX i (v) M='6')' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M=6 DISTANCE=euclidean)
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX i (v) M=6 DISTANCE=euclidean)' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) DISTANCE=cosine)
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX i (v) DISTANCE=cosine)' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) DISTANCE=foo)
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX i (v) DISTANCE=foo)' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) DISTANCE='cosine')
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX i (v) DISTANCE='cosine')' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) DISTANCE=default)
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX i (v) DISTANCE=default)' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M=default)
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX i (v) M=default)' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) COMMENT 'c')
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX i (v) COMMENT 'c')' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) VISIBLE)
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX i (v) VISIBLE)' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) IGNORED)
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX i (v) IGNORED)' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) NOT IGNORED)
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX i (v) NOT IGNORED)' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) KEY_BLOCK_SIZE=1)
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX i (v) KEY_BLOCK_SIZE=1)' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH PARSER p)
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH PARSER p)' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) FOO=1)
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX i (v) FOO=1)' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=l2, type=hnsw, lib=vsag))
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=l2, type=hnsw, lib=vsag))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=l2, type=hnsw, lib=vsag, m=16, ef_construction=200, ef_search=64))
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=l2, type=hnsw, lib=vsag, m' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR KEY i (v) WITH (distance=l2, type=hnsw, lib=vsag))
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR KEY i (v) WITH (distance=l2, type=hnsw, lib=vsag))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=cosine, type=hnsw))
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=cosine, type=hnsw))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=inner_product, type=hnsw_sq, lib=vsag))
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=inner_product, type=hnsw_s' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=l2, type=ivf_flat))
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=l2, type=ivf_flat))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH ())
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH ())' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=foo, type=bar))
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=foo, type=bar))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (foo=bar))
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (foo=bar))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=l2 type=hnsw))
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=l2 type=hnsw))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH distance=l2)
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH distance=l2)' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=l2, type=hnsw, lib=vsag) COMMENT 'c')
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=l2, type=hnsw, lib=vsag) C' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX (v) WITH (distance=l2, type=hnsw, lib=vsag))
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, VECTOR INDEX (v) WITH (distance=l2, type=hnsw, lib=vsag))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, CONSTRAINT c VECTOR INDEX (v))
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, CONSTRAINT c VECTOR INDEX (v))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, CONSTRAINT VECTOR INDEX (v))
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL, CONSTRAINT VECTOR INDEX (v))' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v)
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v)' at line 1
> exec ^
: string
| CREATE VECTOR KEY i ON t (v)
< error
| 1064 42000 …near 'VECTOR KEY i ON t (v)' at line 1
> exec ^
: string
| CREATE VECTOR i ON t (v)
< error
| 1064 42000 …near 'VECTOR i ON t (v)' at line 1
> exec ^
: string
| CREATE VECTOR INDEX ON t (v)
< error
| 1064 42000 …near 'VECTOR INDEX ON t (v)' at line 1
> exec ^
: string
| CREATE VECTOR INDEX
< error
| 1064 42000 …near 'VECTOR INDEX' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t
< error
| 1064 42000 …near 'VECTOR INDEX i ON t' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t ()
< error
| 1064 42000 …near 'VECTOR INDEX i ON t ()' at line 1
> exec ^
: string
| CREATE VECTOR
< error
| 1064 42000 …near 'VECTOR' at line 1
> exec ^
: string
| CREATE OR REPLACE VECTOR INDEX i ON t (v)
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v)' at line 1
> exec ^
: string
| CREATE VECTOR INDEX IF NOT EXISTS i ON t (v)
< error
| 1064 42000 …near 'VECTOR INDEX IF NOT EXISTS i ON t (v)' at line 1
> exec ^
: string
| CREATE OR REPLACE VECTOR INDEX IF NOT EXISTS i ON t (v)
< error
| 1064 42000 …near 'VECTOR INDEX IF NOT EXISTS i ON t (v)' at line 1
> exec ^
: string
| CREATE UNIQUE VECTOR INDEX i ON t (v)
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v)' at line 1
> exec ^
: string
| CREATE VECTOR UNIQUE INDEX i ON t (v)
< error
| 1064 42000 …near 'VECTOR UNIQUE INDEX i ON t (v)' at line 1
> exec ^
: string
| CREATE TEMPORARY VECTOR INDEX i ON t (v)
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v)' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i USING HNSW ON t (v)
< error
| 1064 42000 …near 'VECTOR INDEX i USING HNSW ON t (v)' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i TYPE HNSW ON t (v)
< error
| 1064 42000 …near 'VECTOR INDEX i TYPE HNSW ON t (v)' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i USING BTREE ON t (v)
< error
| 1064 42000 …near 'VECTOR INDEX i USING BTREE ON t (v)' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) USING HNSW
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v) USING HNSW' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) USING BTREE
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v) USING BTREE' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t ((VEC_COSINE_DISTANCE(v)))
< error
| 1064 42000 …near 'VECTOR INDEX i ON t ((VEC_COSINE_DISTANCE(v)))' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v, w)
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v, w)' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v(3))
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v(3))' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v DESC)
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v DESC)' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) M=6
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v) M=6' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) M=6 DISTANCE=cosine
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v) M=6 DISTANCE=cosine' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) DISTANCE=euclidean M=6
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v) DISTANCE=euclidean M=6' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) M=abc
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v) M=abc' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) DISTANCE=foo
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v) DISTANCE=foo' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) COMMENT 'c'
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v) COMMENT 'c'' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) ALGORITHM=DEFAULT
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v) ALGORITHM=DEFAULT' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) ALGORITHM=INPLACE LOCK=NONE
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v) ALGORITHM=INPLACE LOCK=NONE' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) LOCK=NONE
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v) LOCK=NONE' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) WAIT 3
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v) WAIT 3' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) NOWAIT
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v) NOWAIT' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i WAIT 3 ON t (v)
< error
| 1064 42000 …near 'VECTOR INDEX i WAIT 3 ON t (v)' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) M=6 ALGORITHM=DEFAULT
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v) M=6 ALGORITHM=DEFAULT' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) WITH (distance=l2, type=hnsw, lib=vsag)
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v) WITH (distance=l2, type=hnsw, lib=vsag)' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) WITH (distance=l2, type=hnsw, lib=vsag) COMMENT 'c'
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v) WITH (distance=l2, type=hnsw, lib=vsag) COMMENT 'c'' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) WITH (distance=cosine, type=hnsw)
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v) WITH (distance=cosine, type=hnsw)' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) WITH (foo=bar)
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v) WITH (foo=bar)' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) WITH ()
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v) WITH ()' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) WITH (distance=l2,)
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v) WITH (distance=l2,)' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) WITH (distance=l2) ALGORITHM=INPLACE
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v) WITH (distance=l2) ALGORITHM=INPLACE' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) SELECT 1
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v) SELECT 1' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) x
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v) x' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v),
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v),' at line 1
> exec ^
: string
| CREATE VECTOR INDEX vector ON vector (vector)
< error
| 1064 42000 …near 'VECTOR INDEX vector ON vector (vector)' at line 1
> exec ^
: string
| CREATE VECTOR INDEX `index` ON t (v)
< error
| 1064 42000 …near 'VECTOR INDEX `index` ON t (v)' at line 1
> exec ^
: string
| CREATE VECTOR INDEX select ON t (v)
< error
| 1064 42000 …near 'VECTOR INDEX select ON t (v)' at line 1
> exec ^
: string
| CREATE TABLE vector (vector INT)
< ok 0 0
> exec ^
: string
| CREATE TABLE vector (vector vector(3))
< error
| 1064 42000 …near 'vector(3))' at line 1
> exec ^
: string
| CREATE TABLE t (vector)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| CREATE TABLE t (vector index)
< error
| 1064 42000 …near 'index)' at line 1
> exec ^
: string
| CREATE TABLE t (vector index INT)
< error
| 1064 42000 …near 'index INT)' at line 1
> exec ^
: string
| CREATE TABLE t (vector INT, INDEX (vector))
< ok 0 0
> exec ^
: string
| SELECT vector FROM vector
< error
| 1146 42S02 Table 'omni_t_mysql_parser.vector' doesn't exist
> exec ^
: string
| SELECT vector FROM vector
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 vector
< ok 0 0
> exec ^
: string
| SELECT vector(3)
< error
| 1305 42000 FUNCTION omni_t_mysql_parser.vector does not exist
> exec ^
: string
| SELECT vector(3)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT vector.vector FROM vector
< error
| 1146 42S02 Table 'omni_t_mysql_parser.vector' doesn't exist
> exec ^
: string
| SELECT vector.vector FROM vector
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE INDEX vector ON t (vector)
< ok 0 0
> exec ^
: string
| DROP INDEX vector ON t
< ok 0 0
> exec ^
: string
| SELECT CAST('[1]' AS VECTOR)
< error
| 1064 42000 …near 'VECTOR)' at line 1
> exec ^
: string
| SELECT CAST('[1]' AS VECTOR(3))
< error
| 1064 42000 …near 'VECTOR(3))' at line 1
> exec ^
: string
| SELECT CAST('[1]' AS VECTOR(3) ARRAY)
< error
| 1064 42000 …near 'VECTOR(3) ARRAY)' at line 1
> exec ^
: string
| SELECT CONVERT('[1]', VECTOR)
< error
| 1064 42000 …near 'VECTOR)' at line 1
> exec ^
: string
| SELECT CONVERT('[1]', VECTOR(3))
< error
| 1064 42000 …near 'VECTOR(3))' at line 1
> exec ^
: string
| ALTER TABLE t ADD COLUMN a VECTOR(3)
< error
| 1064 42000 …near 'VECTOR(3)' at line 1
> exec ^
: string
| ALTER TABLE t ADD a VECTOR(3)
< error
| 1064 42000 …near 'VECTOR(3)' at line 1
> exec ^
: string
| ALTER TABLE t ADD COLUMN (a VECTOR(3), b INT)
< error
| 1064 42000 …near 'VECTOR(3), b INT)' at line 1
> exec ^
: string
| ALTER TABLE t ADD vector vector(3)
< error
| 1064 42000 …near 'vector(3)' at line 1
> exec ^
: string
| ALTER TABLE t ADD vector INT
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD COLUMN vector INT
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD VECTOR
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| ALTER TABLE t ADD VECTOR (a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| ALTER TABLE t ADD VECTOR KEY (a)
< error
| 1064 42000 …near 'KEY (a)' at line 1
> exec ^
: string
| ALTER TABLE t ADD VECTOR KEY i (a)
< error
| 1064 42000 …near 'KEY i (a)' at line 1
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX (a)
< error
| 1064 42000 …near 'INDEX (a)' at line 1
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX i (a)
< error
| 1064 42000 …near 'INDEX i (a)' at line 1
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX IF NOT EXISTS i (a)
< error
| 1064 42000 …near 'INDEX IF NOT EXISTS i (a)' at line 1
> exec ^
: string
| ALTER TABLE t ADD VECTOR i (a)
< error
| 1064 42000 …near 'i (a)' at line 1
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX i (a) M=6 DISTANCE=cosine
< error
| 1064 42000 …near 'INDEX i (a) M=6 DISTANCE=cosine' at line 1
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX i (a) WITH (distance=l2, type=hnsw, lib=vsag)
< error
| 1064 42000 …near 'INDEX i (a) WITH (distance=l2, type=hnsw, lib=vsag)' at line 1
> exec ^
: string
| ALTER TABLE t ADD VECTOR KEY i (a) WITH (distance=l2, type=hnsw, lib=vsag)
< error
| 1064 42000 …near 'KEY i (a) WITH (distance=l2, type=hnsw, lib=vsag)' at line 1
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX i (a) USING HNSW
< error
| 1064 42000 …near 'INDEX i (a) USING HNSW' at line 1
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX ((a))
< error
| 1064 42000 …near 'INDEX ((a))' at line 1
> exec ^
: string
| ALTER TABLE t ADD CONSTRAINT c VECTOR INDEX (a)
< error
| 1064 42000 …near 'VECTOR INDEX (a)' at line 1
> exec ^
: string
| ALTER TABLE t MODIFY a VECTOR(3)
< error
| 1064 42000 …near 'VECTOR(3)' at line 1
> exec ^
: string
| ALTER TABLE t MODIFY COLUMN a VECTOR(3) NOT NULL
< error
| 1064 42000 …near 'VECTOR(3) NOT NULL' at line 1
> exec ^
: string
| ALTER TABLE t MODIFY a VECTOR
< error
| 1064 42000 …near 'VECTOR' at line 1
> exec ^
: string
| ALTER TABLE t CHANGE a b VECTOR(3)
< error
| 1064 42000 …near 'VECTOR(3)' at line 1
> exec ^
: string
| ALTER TABLE t DROP INDEX i
< ok 0 0
> exec ^
: string
| SELECT VEC_FROMTEXT('[1]'), VEC_TOTEXT(a), VEC_DISTANCE_COSINE(a,a), VEC_DISTANCE_EUCLIDEAN(a,a), VEC_DISTANCE(a,a) FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT VEC_FROMTEXT('[1]'), VEC_TOTEXT(a), VEC_DISTANCE_COSINE(a,a), VEC_DISTANCE_EUCLIDEAN(a,a), VEC_DISTANCE(a,a) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT l2_distance(a,'[1]'), cosine_distance(a,'[1]'), inner_product(a,'[1]') FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT l2_distance(a,'[1]'), cosine_distance(a,'[1]'), inner_product(a,'[1]') FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT * FROM t ORDER BY a <-> '[1,2,3]' LIMIT 1
< error
| 1064 42000 …near '-> '[1,2,3]' LIMIT 1' at line 1
> exec ^
: string
| SELECT * FROM t ORDER BY a <=> '[1,2,3]' APPROXIMATE LIMIT 1
< error
| 1064 42000 …near 'APPROXIMATE LIMIT 1' at line 1
> exec ^
: string
| SELECT * FROM t ORDER BY a <@> '[1,2,3]' APPROX LIMIT 1
< error
| 1064 42000 …near 'APPROX LIMIT 1' at line 1
> exec ^
: string
| ALTER TABLE t ADD (vector INT)
< ok 0 0
> exec ^
: string
| CREATE FUNCTION f(a VECTOR(3)) RETURNS INT RETURN 1
< error
| 1064 42000 …near 'VECTOR(3)) RETURNS INT RETURN 1' at line 1
