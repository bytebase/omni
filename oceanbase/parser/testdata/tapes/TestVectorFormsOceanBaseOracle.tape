-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestVectorFormsOceanBaseOracle default#1
> exec
| DROP DATABASE IF EXISTS obvec
< ok 0 0
> exec
| CREATE DATABASE obvec
< ok 1 0
> exec
| USE obvec
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, v VECTOR(3) NOT NULL)
< ok 0 0
> exec
| CREATE TABLE t (a VECTOR)
< error
| 1064 42000 …near ')' at line 1
> exec
| CREATE TABLE t (a VECTOR(3))
< error
| 1050 42S01 Table 't' already exists
> exec
| CREATE TABLE t (a VECTOR(0))
< error
| 1235 0A000 vector column dim less or equal to zero or larger than 16000 is not supported
> exec
| CREATE TABLE t (a VECTOR(16383))
< error
| 1235 0A000 vector column dim less or equal to zero or larger than 16000 is not supported
> exec
| CREATE TABLE t (a VECTOR(16384))
< error
| 1235 0A000 vector column dim less or equal to zero or larger than 16000 is not supported
> exec
| CREATE TABLE t (a VECTOR(99999999999999999999))
< error
| 1064 42000 …near '99999999999999999999))' at line 1
> exec
| CREATE TABLE t (a VECTOR<FLOAT>)
< error
| 1064 42000 …near '<FLOAT>)' at line 1
> exec
| CREATE TABLE t (a VECTOR<FLOAT>(3))
< error
| 1064 42000 …near '<FLOAT>(3))' at line 1
> exec
| CREATE TABLE t (a VECTOR(5)<FLOAT>)
< error
| 1064 42000 …near '<FLOAT>)' at line 1
> exec
| CREATE TABLE t (a VECTOR())
< error
| 1064 42000 …near '))' at line 1
> exec
| CREATE TABLE t (a VECTOR(3,2))
< error
| 1064 42000 …near ',2))' at line 1
> exec
| CREATE TABLE t (a VECTOR('3'))
< error
| 1064 42000 …near ''3'))' at line 1
> exec
| CREATE TABLE t (a VECTOR(b))
< error
| 1064 42000 …near 'b))' at line 1
> exec
| CREATE TABLE t (a VECTOR(1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec
| CREATE TABLE t (a VECTOR(0x3))
< error
| 1064 42000 …near '0x3))' at line 1
> exec
| CREATE TABLE t (a VECTOR(-1))
< error
| 1064 42000 …near '-1))' at line 1
> exec
| CREATE TABLE t (a VECTOR(3)
< error
| 1064 42000 …near ')' at line 1
> exec
| CREATE TABLE t (a VECTOR(3, FLOAT))
< error
| 1064 42000 …near ', FLOAT))' at line 1
> exec
| CREATE TABLE t (a VECTOR(3) NOT NULL)
< error
| 1050 42S01 Table 't' already exists
> exec
| CREATE TABLE t (a VECTOR(3) NULL)
< error
| 1050 42S01 Table 't' already exists
> exec
| CREATE TABLE t (a VECTOR(3) COMMENT 'x')
< error
| 1050 42S01 Table 't' already exists
> exec
| CREATE TABLE t (a VECTOR(3) DEFAULT NULL)
< error
| 1050 42S01 Table 't' already exists
> exec
| CREATE TABLE t (a VECTOR(3) DEFAULT '[1,2,3]')
< error
| 1101 42000 BLOB, TEXT, GEOMETRY or JSON column 'a' can't have a default value
> exec
| CREATE TABLE t (a VECTOR(3) UNSIGNED)
< error
| 1064 42000 …near 'UNSIGNED)' at line 1
> exec
| CREATE TABLE t (a VECTOR(3) ZEROFILL)
< error
| 1064 42000 …near 'ZEROFILL)' at line 1
> exec
| CREATE TABLE t (a VECTOR(3) BINARY)
< error
| 1064 42000 …near 'BINARY)' at line 1
> exec
| CREATE TABLE t (a VECTOR(3) CHARSET utf8mb4)
< error
| 1064 42000 …near 'CHARSET utf8mb4)' at line 1
> exec
| CREATE TABLE t (a VECTOR(3) CHARACTER SET utf8mb4)
< error
| 1064 42000 …near 'CHARACTER SET utf8mb4)' at line 1
> exec
| CREATE TABLE t (a VECTOR(3) COLLATE utf8mb4_bin)
< error
| 1050 42S01 Table 't' already exists
> exec
| CREATE TABLE t (a VECTOR(3) PRIMARY KEY)
< error
| 1167 42000 The used storage engine can't index column 'a'
> exec
| CREATE TABLE t (a VECTOR(3) ARRAY)
< error
| 1064 42000 …near 'ARRAY)' at line 1
> exec
| CREATE TABLE t (a VECTOR(3) GENERATED ALWAYS AS ('[1]'))
< error
| 1050 42S01 Table 't' already exists
> exec
| CREATE TABLE t (a VECTOR(3) AS ('[1]') STORED)
< error
| 1050 42S01 Table 't' already exists
> exec
| CREATE TABLE t (a `VECTOR`)
< error
| 1064 42000 …near '`VECTOR`)' at line 1
> exec
| CREATE TABLE t (a `VECTOR`(3))
< error
| 1064 42000 …near '`VECTOR`(3))' at line 1
> exec
| CREATE TABLE t (a vector(3))
< error
| 1050 42S01 Table 't' already exists
> exec
| CREATE TEMPORARY TABLE tt (a VECTOR(3))
< ok 0 0
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX (v))
< error
| 1235 0A000 vector index params not set is not supported
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR KEY (v))
< error
| 1235 0A000 vector index params not set is not supported
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR (v))
< error
| 1064 42000 …near '(v))' at line 1
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR i (v))
< error
| 1064 42000 …near 'i (v))' at line 1
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v))
< error
| 1235 0A000 vector index params not set is not supported
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR KEY i (v))
< error
| 1235 0A000 vector index params not set is not supported
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX IF NOT EXISTS i (v))
< error
| 1064 42000 …near 'IF NOT EXISTS i (v))' at line 1
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i USING HNSW (v))
< error
| 1064 42000 …near 'USING HNSW (v))' at line 1
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) USING HNSW)
< error
| 1064 42000 …near 'HNSW)' at line 1
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) USING BTREE)
< error
| 1235 0A000 vector index params not set is not supported
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) TYPE HNSW)
< error
| 1064 42000 …near 'TYPE HNSW)' at line 1
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX ((VEC_COSINE_DISTANCE(v))))
< error
| 1235 0A000 using udf as generated column not supported
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX (v, w))
< error
| 1235 0A000 multi vector index column is not supported
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX (v(3)))
< error
| 1089 HY000 Incorrect prefix key; the used key part isn't a string, the used length is longer than the key part, or the storage engine doesn't support unique prefix keys
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX (v DESC))
< error
| 1235 0A000 vector index params not set is not supported
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M=6)
< error
| 1064 42000 …near 'M=6)' at line 1
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M 6)
< error
| 1064 42000 …near 'M 6)' at line 1
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M=abc)
< error
| 1064 42000 …near 'M=abc)' at line 1
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M='6')
< error
| 1064 42000 …near 'M='6')' at line 1
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M=6 DISTANCE=euclidean)
< error
| 1064 42000 …near 'M=6 DISTANCE=euclidean)' at line 1
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) DISTANCE=cosine)
< error
| 1064 42000 …near 'DISTANCE=cosine)' at line 1
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) DISTANCE=foo)
< error
| 1064 42000 …near 'DISTANCE=foo)' at line 1
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) DISTANCE='cosine')
< error
| 1064 42000 …near 'DISTANCE='cosine')' at line 1
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) DISTANCE=default)
< error
| 1064 42000 …near 'DISTANCE=default)' at line 1
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M=default)
< error
| 1064 42000 …near 'M=default)' at line 1
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) COMMENT 'c')
< error
| 1235 0A000 vector index params not set is not supported
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) VISIBLE)
< error
| 1235 0A000 vector index params not set is not supported
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) IGNORED)
< error
| 1064 42000 …near 'IGNORED)' at line 1
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) NOT IGNORED)
< error
| 1064 42000 …near 'NOT IGNORED)' at line 1
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) KEY_BLOCK_SIZE=1)
< error
| 1235 0A000 vector index params not set is not supported
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH PARSER p)
< error
| 1064 42000 …near 'PARSER p)' at line 1
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) FOO=1)
< error
| 1064 42000 …near 'FOO=1)' at line 1
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=l2, type=hnsw, lib=vsag))
< error
| 1050 42S01 Table 't' already exists
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=l2, type=hnsw, lib=vsag, m=16, ef_construction=200, ef_search=64))
< error
| 1050 42S01 Table 't' already exists
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR KEY i (v) WITH (distance=l2, type=hnsw, lib=vsag))
< error
| 1050 42S01 Table 't' already exists
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=cosine, type=hnsw))
< error
| 1050 42S01 Table 't' already exists
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=inner_product, type=hnsw_sq, lib=vsag))
< error
| 1050 42S01 Table 't' already exists
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=l2, type=ivf_flat))
< error
| 1050 42S01 Table 't' already exists
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH ())
< error
| 1064 42000 …near '))' at line 1
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=foo, type=bar))
< error
| 1235 0A000 this type of vector index distance algorithm is not supported
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (foo=bar))
< error
| 1235 0A000 vector index params not set distance and type is not supported
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=l2 type=hnsw))
< error
| 1064 42000 …near 'type=hnsw))' at line 1
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH distance=l2)
< error
| 1064 42000 …near 'distance=l2)' at line 1
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=l2, type=hnsw, lib=vsag) COMMENT 'c')
< error
| 1050 42S01 Table 't' already exists
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX (v) WITH (distance=l2, type=hnsw, lib=vsag))
< error
| 1050 42S01 Table 't' already exists
> exec
| CREATE TABLE t (v VECTOR(3) NOT NULL, CONSTRAINT c VECTOR INDEX (v))
< error
| 1064 42000 …near 'VECTOR INDEX (v))' at line 1
> exec
| CREATE VECTOR INDEX i ON t (v)
< error
| 1235 0A000 vector index params not set is not supported
> exec
| CREATE VECTOR KEY i ON t (v)
< error
| 1064 42000 …near 'KEY i ON t (v)' at line 1
> exec
| CREATE VECTOR i ON t (v)
< error
| 1064 42000 …near 'i ON t (v)' at line 1
> exec
| CREATE VECTOR INDEX ON t (v)
< error
| 1064 42000 …near 'ON t (v)' at line 1
> exec
| CREATE VECTOR INDEX
< error
| 1064 42000 …near 'INDEX' at line 1
> exec
| CREATE VECTOR INDEX i ON t
< error
| 1064 42000 …near 't' at line 1
> exec
| CREATE VECTOR INDEX i ON t ()
< error
| 1064 42000 …near ')' at line 1
> exec
| CREATE VECTOR
< error
| 1064 42000 …near 'VECTOR' at line 1
> exec
| CREATE OR REPLACE VECTOR INDEX i ON t (v)
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v)' at line 1
> exec
| CREATE VECTOR INDEX IF NOT EXISTS i ON t (v)
< error
| 1235 0A000 vector index params not set is not supported
> exec
| CREATE OR REPLACE VECTOR INDEX IF NOT EXISTS i ON t (v)
< error
| 1064 42000 …near 'VECTOR INDEX IF NOT EXISTS i ON t (v)' at line 1
> exec
| CREATE UNIQUE VECTOR INDEX i ON t (v)
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v)' at line 1
> exec
| CREATE VECTOR UNIQUE INDEX i ON t (v)
< error
| 1064 42000 …near 'UNIQUE INDEX i ON t (v)' at line 1
> exec
| CREATE TEMPORARY VECTOR INDEX i ON t (v)
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v)' at line 1
> exec
| CREATE VECTOR INDEX i USING HNSW ON t (v)
< error
| 1064 42000 …near 'HNSW ON t (v)' at line 1
> exec
| CREATE VECTOR INDEX i TYPE HNSW ON t (v)
< error
| 1064 42000 …near 'TYPE HNSW ON t (v)' at line 1
> exec
| CREATE VECTOR INDEX i USING BTREE ON t (v)
< error
| 1235 0A000 vector index params not set is not supported
> exec
| CREATE VECTOR INDEX i ON t (v) USING HNSW
< error
| 1064 42000 …near 'HNSW' at line 1
> exec
| CREATE VECTOR INDEX i ON t (v) USING BTREE
< error
| 1235 0A000 vector index params not set is not supported
> exec
| CREATE VECTOR INDEX i ON t ((VEC_COSINE_DISTANCE(v)))
< error
| 3762 HY000 Functional index on a column is not supported. Consider using a regular index instead.
> exec
| CREATE VECTOR INDEX i ON t (v, w)
< error
| 1235 0A000 multi vector index column is not supported
> exec
| CREATE VECTOR INDEX i ON t (v(3))
< error
| 1235 0A000 vector index params not set is not supported
> exec
| CREATE VECTOR INDEX i ON t (v DESC)
< error
| 1235 0A000 vector index params not set is not supported
> exec
| CREATE VECTOR INDEX i ON t (v) M=6
< error
| 1064 42000 …near 'M=6' at line 1
> exec
| CREATE VECTOR INDEX i ON t (v) M=6 DISTANCE=cosine
< error
| 1064 42000 …near 'M=6 DISTANCE=cosine' at line 1
> exec
| CREATE VECTOR INDEX i ON t (v) DISTANCE=euclidean M=6
< error
| 1064 42000 …near 'DISTANCE=euclidean M=6' at line 1
> exec
| CREATE VECTOR INDEX i ON t (v) M=abc
< error
| 1064 42000 …near 'M=abc' at line 1
> exec
| CREATE VECTOR INDEX i ON t (v) DISTANCE=foo
< error
| 1064 42000 …near 'DISTANCE=foo' at line 1
> exec
| CREATE VECTOR INDEX i ON t (v) COMMENT 'c'
< error
| 1235 0A000 vector index params not set is not supported
> exec
| CREATE VECTOR INDEX i ON t (v) ALGORITHM=DEFAULT
< error
| 1064 42000 …near 'ALGORITHM=DEFAULT' at line 1
> exec
| CREATE VECTOR INDEX i ON t (v) ALGORITHM=INPLACE LOCK=NONE
< error
| 1064 42000 …near 'ALGORITHM=INPLACE LOCK=NONE' at line 1
> exec
| CREATE VECTOR INDEX i ON t (v) LOCK=NONE
< error
| 1064 42000 …near 'LOCK=NONE' at line 1
> exec
| CREATE VECTOR INDEX i ON t (v) WAIT 3
< error
| 1064 42000 …near 'WAIT 3' at line 1
> exec
| CREATE VECTOR INDEX i ON t (v) NOWAIT
< error
| 1064 42000 …near 'NOWAIT' at line 1
> exec
| CREATE VECTOR INDEX i WAIT 3 ON t (v)
< error
| 1064 42000 …near 'WAIT 3 ON t (v)' at line 1
> exec
| CREATE VECTOR INDEX i ON t (v) M=6 ALGORITHM=DEFAULT
< error
| 1064 42000 …near 'M=6 ALGORITHM=DEFAULT' at line 1
> exec
| CREATE VECTOR INDEX i ON t (v) WITH (distance=l2, type=hnsw, lib=vsag)
< ok 0 0
> exec
| CREATE VECTOR INDEX i ON t (v) WITH (distance=l2, type=hnsw, lib=vsag) COMMENT 'c'
< error
| 1235 0A000 create vector index on column has vector index is not supported
> exec
| CREATE VECTOR INDEX i ON t (v) WITH (distance=cosine, type=hnsw)
< error
| 1235 0A000 create vector index on column has vector index is not supported
> exec
| CREATE VECTOR INDEX i ON t (v) WITH (foo=bar)
< error
| 1235 0A000 create vector index on column has vector index is not supported
> exec
| CREATE VECTOR INDEX i ON t (v) WITH ()
< error
| 1064 42000 …near ')' at line 1
> exec
| CREATE VECTOR INDEX i ON t (v) WITH (distance=l2,)
< error
| 1064 42000 …near ')' at line 1
> exec
| CREATE VECTOR INDEX i ON t (v) WITH (distance=l2) ALGORITHM=INPLACE
< error
| 1064 42000 …near 'ALGORITHM=INPLACE' at line 1
> exec
| CREATE VECTOR INDEX i ON t (v) SELECT 1
< error
| 1064 42000 …near 'SELECT 1' at line 1
> exec
| CREATE VECTOR INDEX i ON t (v) x
< error
| 1064 42000 …near 'x' at line 1
> exec
| CREATE VECTOR INDEX i ON t (v),
< error
| 1064 42000 …near ',' at line 1
> exec
| CREATE VECTOR INDEX vector ON vector (vector)
< error
| 1146 42S02 Table 'obvec.vector' doesn't exist
> exec
| CREATE VECTOR INDEX `index` ON t (v)
< error
| 1235 0A000 create vector index on column has vector index is not supported
> exec
| CREATE VECTOR INDEX select ON t (v)
< error
| 1064 42000 …near 'select ON t (v)' at line 1
> exec
| CREATE TABLE vector (vector INT)
< ok 0 0
> exec
| CREATE TABLE vector (vector vector(3))
< error
| 1050 42S01 Table 'vector' already exists
> exec
| CREATE TABLE t (vector)
< error
| 1064 42000 …near ')' at line 1
> exec
| CREATE TABLE t (vector index)
< error
| 1064 42000 …near ')' at line 1
> exec
| CREATE TABLE t (vector index INT)
< error
| 1064 42000 …near 'INT)' at line 1
> exec
| CREATE TABLE t (vector INT, INDEX (vector))
< error
| 1050 42S01 Table 't' already exists
> exec
| SELECT 1 vector
< ok 0 0
> exec
| CREATE INDEX vector ON t (vector)
< error
| 1072 42000 Key column 'vector' doesn't exist in table
> exec
| DROP INDEX vector ON t
< error
| 1091 42000 Can't DROP 'vector'; check that column/key exists
> exec
| SELECT CAST('[1]' AS VECTOR)
< error
| 1064 42000 …near 'VECTOR)' at line 1
> exec
| SELECT CAST('[1]' AS VECTOR(3))
< error
| 1064 42000 …near 'VECTOR(3))' at line 1
> exec
| SELECT CAST('[1]' AS VECTOR(3) ARRAY)
< error
| 1064 42000 …near 'VECTOR(3) ARRAY)' at line 1
> exec
| SELECT CONVERT('[1]', VECTOR)
< error
| 1064 42000 …near 'VECTOR)' at line 1
> exec
| SELECT CONVERT('[1]', VECTOR(3))
< error
| 1064 42000 …near 'VECTOR(3))' at line 1
> exec
| ALTER TABLE t ADD COLUMN a VECTOR(3)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| ALTER TABLE t ADD a VECTOR(3)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| ALTER TABLE t ADD COLUMN (a VECTOR(3), b INT)
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| ALTER TABLE t ADD vector vector(3)
< ok 0 0
> exec
| ALTER TABLE t ADD vector INT
< error
| 1060 42S21 Duplicate column name 'vector'
> exec
| ALTER TABLE t ADD COLUMN vector INT
< error
| 1060 42S21 Duplicate column name 'vector'
> exec
| ALTER TABLE t ADD VECTOR
< error
| 1064 42000 …near 'VECTOR' at line 1
> exec
| ALTER TABLE t ADD VECTOR (a)
< error
| 1064 42000 …near '(a)' at line 1
> exec
| ALTER TABLE t ADD VECTOR KEY (a)
< error
| 1064 42000 …near 'KEY (a)' at line 1
> exec
| ALTER TABLE t ADD VECTOR KEY i (a)
< error
| 1064 42000 …near 'KEY i (a)' at line 1
> exec
| ALTER TABLE t ADD VECTOR INDEX (a)
< error
| 1064 42000 …near 'INDEX (a)' at line 1
> exec
| ALTER TABLE t ADD VECTOR INDEX i (a)
< error
| 1064 42000 …near 'INDEX i (a)' at line 1
> exec
| ALTER TABLE t ADD VECTOR INDEX IF NOT EXISTS i (a)
< error
| 1064 42000 …near 'INDEX IF NOT EXISTS i (a)' at line 1
> exec
| ALTER TABLE t ADD VECTOR i (a)
< error
| 1064 42000 …near 'i (a)' at line 1
> exec
| ALTER TABLE t ADD VECTOR INDEX i (a) M=6 DISTANCE=cosine
< error
| 1064 42000 …near 'INDEX i (a) M=6 DISTANCE=cosine' at line 1
> exec
| ALTER TABLE t ADD VECTOR INDEX i (a) WITH (distance=l2, type=hnsw, lib=vsag)
< error
| 1064 42000 …near 'INDEX i (a) WITH (distance=l2, type=hnsw, lib=vsag)' at line 1
> exec
| ALTER TABLE t ADD VECTOR KEY i (a) WITH (distance=l2, type=hnsw, lib=vsag)
< error
| 1064 42000 …near 'KEY i (a) WITH (distance=l2, type=hnsw, lib=vsag)' at line 1
> exec
| ALTER TABLE t ADD VECTOR INDEX i (a) USING HNSW
< error
| 1064 42000 …near 'INDEX i (a) USING HNSW' at line 1
> exec
| ALTER TABLE t ADD VECTOR INDEX ((a))
< error
| 1064 42000 …near 'INDEX ((a))' at line 1
> exec
| ALTER TABLE t ADD CONSTRAINT c VECTOR INDEX (a)
< error
| 1064 42000 …near 'VECTOR INDEX (a)' at line 1
> exec
| ALTER TABLE t MODIFY a VECTOR(3)
< error
| 5083 22000 Invalid data type for the operation
> exec
| ALTER TABLE t MODIFY COLUMN a VECTOR(3) NOT NULL
< error
| 5083 22000 Invalid data type for the operation
> exec
| ALTER TABLE t MODIFY a VECTOR
< error
| 1064 42000 …near 'VECTOR' at line 1
> exec
| ALTER TABLE t CHANGE a b VECTOR(3)
< error
| 1060 42S21 Duplicate column name 'b'
> exec
| ALTER TABLE t DROP INDEX i
< ok 0 0
> exec
| CREATE VECTOR INDEX i1 ON t (v) COMMENT 'c' WITH (distance=l2, type=hnsw, lib=vsag)
< ok 0 0
> exec
| CREATE VECTOR INDEX i2 ON t (v) WITH (distance=l2, type=hnsw, lib=vsag) WITH (distance=l2, type=hnsw, lib=vsag)
< error
| 1235 0A000 create vector index on column has vector index is not supported
> exec
| CREATE VECTOR INDEX i3 ON t (v) USING BTREE WITH (distance=l2, type=hnsw, lib=vsag)
< error
| 1235 0A000 create vector index on column has vector index is not supported
> exec
| CREATE VECTOR INDEX i4 ON t (v) WITH (distance=l2, type=hnsw, lib=vsag) USING BTREE
< error
| 1235 0A000 create vector index on column has vector index is not supported
> exec
| CREATE VECTOR INDEX i5 USING BTREE ON t (v) WITH (distance=l2, type=hnsw, lib=vsag)
< error
| 1235 0A000 create vector index on column has vector index is not supported
> exec
| CREATE VECTOR INDEX i6 USING HASH ON t (v) WITH (distance=l2, type=hnsw, lib=vsag)
< error
| 1235 0A000 create vector index on column has vector index is not supported
> exec
| CREATE VECTOR INDEX i7 ON t (v) WITH (distance=l2, distance=l2)
< error
| 1235 0A000 create vector index on column has vector index is not supported
> exec
| CREATE VECTOR INDEX i8 ON t (v) WITH (distance='l2', type='hnsw', lib='vsag')
< error
| 1064 42000 …near ''l2', type='hnsw', lib='vsag')' at line 1
> exec
| CREATE VECTOR INDEX i9 ON t (v) WITH (m=16)
< error
| 1235 0A000 create vector index on column has vector index is not supported
> exec
| CREATE VECTOR INDEX i10 ON t (v) WITH (m='16')
< error
| 1064 42000 …near ''16')' at line 1
> exec
| CREATE VECTOR INDEX i11 ON t (v) WITH (ef_search=1.5)
< error
| 1235 0A000 create vector index on column has vector index is not supported
> exec
| CREATE VECTOR INDEX i12 ON t (v) WITH (`distance`=l2)
< error
| 1235 0A000 create vector index on column has vector index is not supported
> exec
| CREATE VECTOR INDEX i13 ON t (v) WITH (distance=`l2`)
< error
| 1235 0A000 create vector index on column has vector index is not supported
> exec
| CREATE VECTOR INDEX i14 ON t (v) WITH (distance l2)
< error
| 1064 42000 …near 'l2)' at line 1
> exec
| CREATE VECTOR INDEX i15 ON t (v) WITH (distance=select)
< error
| 1064 42000 …near 'select)' at line 1
> exec
| CREATE VECTOR INDEX i16 ON t (v) WITH (select=1)
< error
| 1064 42000 …near 'select=1)' at line 1
> exec
| CREATE VECTOR INDEX i17 ON t (v) WITH (distance=)
< error
| 1064 42000 …near ')' at line 1
> exec
| CREATE VECTOR INDEX i18 ON t (v) WITH (=l2)
< error
| 1064 42000 …near '=l2)' at line 1
> exec
| CREATE VECTOR INDEX i19 ON t (v) WITH (a=-1)
< error
| 1064 42000 …near '-1)' at line 1
> exec
| CREATE VECTOR INDEX i20 ON t (v) WITH (a=+1)
< error
| 1064 42000 …near '+1)' at line 1
> exec
| CREATE VECTOR INDEX i21 ON t (v) WITH (a=NULL)
< error
| 1064 42000 …near 'NULL)' at line 1
> exec
| CREATE VECTOR INDEX i22 ON t (v) WITH (a=DEFAULT)
< error
| 1064 42000 …near 'DEFAULT)' at line 1
> exec
| CREATE VECTOR INDEX i23 ON t (v) WITH (a=TRUE)
< error
| 1235 0A000 create vector index on column has vector index is not supported
> exec
| CREATE VECTOR INDEX i24 ON t (v) WITH (a=(1))
< error
| 1064 42000 …near '(1))' at line 1
> exec
| CREATE VECTOR INDEX i25 ON t (v) WITH (a=1+1)
< error
| 1064 42000 …near '+1)' at line 1
> exec
| CREATE VECTOR INDEX i26 ON t (v) WITH (a=b.c)
< error
| 1064 42000 …near 'b.c)' at line 1
> exec
| CREATE VECTOR INDEX i30 ON t (v) WITH (type=hnsw) GLOBAL
< error
| 1235 0A000 create vector index on column has vector index is not supported
> exec
| CREATE VECTOR INDEX i31 ON t (v) WITH (type=hnsw) LOCAL
< error
| 1235 0A000 create vector index on column has vector index is not supported
> exec
| CREATE VECTOR INDEX i32 ON t (v) GLOBAL WITH (type=hnsw)
< error
| 1235 0A000 create vector index on column has vector index is not supported
> exec
| CREATE VECTOR INDEX i34 ON t (v) WITH (type=hnsw) BLOCK_SIZE 16384
< error
| 1235 0A000 create vector index on column has vector index is not supported
> exec
| CREATE VECTOR INDEX i35 ON t (v) WITH (type=hnsw) KEY_BLOCK_SIZE=1
< error
| 1235 0A000 create vector index on column has vector index is not supported
> exec
| CREATE VECTOR INDEX i36 ON t (v) WITH (type=hnsw) VISIBLE
< error
| 1235 0A000 create vector index on column has vector index is not supported
> exec
| CREATE VECTOR INDEX i37 ON t (v) INVISIBLE WITH (type=hnsw)
< error
| 1235 0A000 create vector index on column has vector index is not supported
> exec
| CREATE VECTOR INDEX i38 ON t (v) WITH PARSER p
< error
| 1064 42000 …near 'PARSER p' at line 1
> exec
| CREATE VECTOR INDEX i39 ON t (v) WITH
< error
| 1064 42000 …near 'WITH' at line 1
> exec
| CREATE VECTOR INDEX i40 ON t (v) WITH (
< error
| 1064 42000 …near '(' at line 1
> exec
| CREATE VECTOR INDEX i41 ON t (v) WITH (a=1
< error
| 1064 42000 …near '1' at line 1
> exec
| CREATE VECTOR INDEX i42 ON t (v) WITH (a=1),
< error
| 1064 42000 …near ',' at line 1
> exec
| CREATE VECTOR INDEX i43 ON t (v) WITH (a=1) x
< error
| 1064 42000 …near 'x' at line 1
> exec
| CREATE VECTOR INDEX i44 ON d.t (v) WITH (a=1)
< error
| 1049 42000 Unknown database
> exec
| CREATE VECTOR INDEX IF EXISTS i46 ON t (v)
< error
| 1064 42000 …near 'EXISTS i46 ON t (v)' at line 1
> exec
| CREATE VECTOR INDEX i47 ON t (v ASC)
< error
| 1235 0A000 create vector index on column has vector index is not supported
> exec
| CREATE VECTOR INDEX i48 ON t (v) WITH (TYPE=HNSW, DISTANCE=L2, LIB=VSAG)
< error
| 1235 0A000 create vector index on column has vector index is not supported
> exec
| ALTER TABLE t ADD (VECTOR INDEX i (v))
< error
| 1064 42000 …near 'INDEX i (v))' at line 1
> exec
| SELECT `vector`(3)
< error
| 1305 42000 FUNCTION vector does not exist
> exec
| CREATE TABLE s6 (a VECTOR(2147483647))
< error
| 1235 0A000 vector column dim less or equal to zero or larger than 16000 is not supported
> exec
| CREATE TABLE s7 (a VECTOR(2147483648))
< error
| 1235 0A000 vector column dim less or equal to zero or larger than 16000 is not supported
> exec
| CREATE TABLE s8 (a VECTOR(4294967295))
< error
| 1235 0A000 vector column dim less or equal to zero or larger than 16000 is not supported
> exec
| CREATE TABLE s9 (a VECTOR(4294967296))
< error
| 1235 0A000 vector column dim less or equal to zero or larger than 16000 is not supported
> exec
| CREATE TABLE s10 (a VECTOR(9223372036854775807))
< error
| 1235 0A000 vector column dim less or equal to zero or larger than 16000 is not supported
> exec
| CREATE TABLE s11 (a VECTOR(9223372036854775808))
< error
| 4016 HY000 Internal error
> exec
| CREATE TABLE s12 (a VECTOR(18446744073709551615))
< error
| 4016 HY000 Internal error
> exec
| CREATE TABLE s13 (a VECTOR(18446744073709551616))
< error
| 1064 42000 …near '18446744073709551616))' at line 1
> exec
| CREATE VECTOR INDEX w1 ON t (v) WITH (a=FALSE)
< error
| 1235 0A000 create vector index on column has vector index is not supported
> exec
| CREATE VECTOR INDEX w2 ON t (v) WITH (a=0x10)
< error
| 1064 42000 …near '0x10)' at line 1
> exec
| CREATE VECTOR INDEX w3 ON t (v) WITH (a=1e5)
< error
| 1235 0A000 create vector index on column has vector index is not supported
> exec
| CREATE VECTOR INDEX w4 ON t (v) WITH (a=.5)
< error
| 1235 0A000 create vector index on column has vector index is not supported
> exec
| SELECT * FROM JSON_TABLE('[]', '$[*]' COLUMNS (a VECTOR(3) PATH '$')) AS j
< error
| 4006 HY000 The object is not initialized
> exec
| SELECT * FROM JSON_TABLE('[]', '$[*]' COLUMNS (a INT PATH '$')) AS j
< ok 0 0
> exec
| CREATE FUNCTION f1(a VECTOR(3)) RETURNS INT RETURN 1
< error
| 1064 42000 …near '(a VECTOR(3)) RETURNS INT RETURN 1' at line 1
> exec
| CREATE INDEX i50 ON t (a) WITH (foo=bar)
< error
| 1235 0A000 set index params in not vector index is not supported
> exec
| CREATE UNIQUE INDEX i51 ON t (a) WITH (foo=bar)
< error
| 1235 0A000 set index params in not vector index is not supported
> exec
| ALTER TABLE t ADD INDEX i56 (a) WITH (foo=bar)
< error
| 1235 0A000 set index params in not vector index is not supported
> exec
| CREATE TABLE tk1 (a INT, KEY k (a) WITH (foo=bar))
< error
| 1235 0A000 set index params in not vector index is not supported
> exec
| CREATE TABLE tk2 (a INT, UNIQUE KEY k (a) WITH (foo=bar))
< error
| 1235 0A000 set index params in not vector index is not supported
> exec
| CREATE TABLE tv1 (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) COMMENT 'c' WITH (distance=l2, type=hnsw, lib=vsag))
< ok 0 0
> exec
| CREATE TABLE tv2 (v VECTOR(3) NOT NULL, VECTOR INDEX `key` (v) WITH (distance=l2, type=hnsw, lib=vsag))
< ok 0 0
> exec
| CREATE TABLE tv3 (v VECTOR(3) NOT NULL, VECTOR INDEX select (v))
< error
| 1064 42000 …near 'select (v))' at line 1
> exec
| CREATE TABLE tv4 (vector VECTOR(3) NOT NULL, VECTOR INDEX vector (vector) WITH (distance=l2, type=hnsw, lib=vsag))
< ok 0 0
> exec
| CREATE TABLE tv5 (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) LOCAL WITH (distance=l2, type=hnsw, lib=vsag))
< ok 0 0
> exec
| CREATE TABLE tv6 (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) GLOBAL WITH (distance=l2, type=hnsw, lib=vsag))
< error
| 1235 0A000 Not supported feature or function
> exec
| CREATE TABLE tv7 (v VECTOR(3) NOT NULL, VECTOR INDEX i USING BTREE (v) WITH (distance=l2, type=hnsw, lib=vsag))
< error
| 1064 42000 …near 'USING BTREE (v) WITH (distance=l2, type=hnsw, lib=vsag))' at line 1
> exec
| CREATE TABLE tv8 (v VECTOR(3) NOT NULL, VECTOR KEY USING BTREE (v) WITH (distance=l2, type=hnsw, lib=vsag))
< error
| 1064 42000 …near 'USING BTREE (v) WITH (distance=l2, type=hnsw, lib=vsag))' at line 1
> exec
| CREATE TABLE tv9 (v VECTOR(3) NOT NULL, VECTOR INDEX USING HASH (v) WITH (distance=l2, type=hnsw, lib=vsag))
< error
| 1064 42000 …near 'USING HASH (v) WITH (distance=l2, type=hnsw, lib=vsag))' at line 1
> exec
| CREATE TABLE tv10 (vector INT KEY, b INT)
< ok 0 0
> exec
| CREATE TABLE tv11 (vector key)
< error
| 1064 42000 …near ')' at line 1
> exec
| CREATE TABLE tv12 (vector key INT)
< error
| 1064 42000 …near 'INT)' at line 1
> exec
| ALTER TABLE t ADD (b2 VECTOR(3), VECTOR INDEX i (b2) WITH (distance=l2, type=hnsw, lib=vsag))
< error
| 1064 42000 …near 'INDEX i (b2) WITH (distance=l2, type=hnsw, lib=vsag))' at line 1
> exec
| CREATE FUNCTION f2() RETURNS VECTOR(3) RETURN 1
< error
| 1064 42000 …near '() RETURNS VECTOR(3) RETURN 1' at line 1
> exec
| CREATE PROCEDURE p2(a VECTOR(3)) SELECT 1
< error
| 1064 42000 …near '(' at line 1
> exec
| CREATE PROCEDURE p3() BEGIN DECLARE b VECTOR(3); END
< error
| 1064 42000 …near '(' at line 1
> exec
| DROP DATABASE IF EXISTS obvec
< ok 45 0
