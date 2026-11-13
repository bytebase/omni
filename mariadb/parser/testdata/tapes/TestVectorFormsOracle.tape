-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestVectorFormsOracle default#1
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t (a VECTOR)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(3))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VECTOR(0))
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a VECTOR(0))
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a VECTOR(0.9))
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a VECTOR(0.9))
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a VECTOR(0,2))
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a VECTOR(0,2))
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a VECTOR(1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VECTOR(1,0))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VECTOR(16383))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VECTOR(16383.5))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VECTOR(16383,2))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VECTOR(16384.5))
< error
| 1074 42000 Column length too big for column 'a' (max = 16383); use BLOB or TEXT instead
> exec ^
: string
| CREATE TABLE t (a VECTOR(16384.5))
| ; SELECT FROM
< error
| 1074 42000 Column length too big for column 'a' (max = 16383); use BLOB or TEXT instead
> exec ^
: string
| CREATE TABLE t (a VECTOR(16384,2))
< error
| 1074 42000 Column length too big for column 'a' (max = 16383); use BLOB or TEXT instead
> exec ^
: string
| CREATE TABLE t (a VECTOR(16384,2))
| ; SELECT FROM
< error
| 1074 42000 Column length too big for column 'a' (max = 16383); use BLOB or TEXT instead
> exec ^
: string
| CREATE TABLE t (a `VECTOR`(0))
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a `VECTOR`(0))
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| ALTER TABLE t MODIFY a VECTOR(0)
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| ALTER TABLE t MODIFY a VECTOR(0)
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| ALTER TABLE t ADD COLUMN a VECTOR(16384)
< error
| 1074 42000 Column length too big for column 'a' (max = 16383); use BLOB or TEXT instead
> exec ^
: string
| ALTER TABLE t ADD COLUMN a VECTOR(16384)
| ; SELECT FROM
< error
| 1074 42000 Column length too big for column 'a' (max = 16383); use BLOB or TEXT instead
> exec ^
: string
| CREATE FUNCTION f(a VECTOR(0)) RETURNS INT RETURN 1
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE FUNCTION f(a VECTOR(0)) RETURNS INT RETURN 1
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE FUNCTION f() RETURNS VECTOR(16384) RETURN 1
< error
| 1074 42000 Column length too big for column '' (max = 16383); use BLOB or TEXT instead
> exec ^
: string
| CREATE FUNCTION f() RETURNS VECTOR(16384) RETURN 1
| ; SELECT FROM
< error
| 1074 42000 Column length too big for column '' (max = 16383); use BLOB or TEXT instead
> exec ^
: string
| CREATE PROCEDURE p(a VECTOR(0)) SELECT 1
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE PROCEDURE p(a VECTOR(0)) SELECT 1
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE b VECTOR(16384); END
< error
| 1074 42000 Column length too big for column 'b' (max = 16383); use BLOB or TEXT instead
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE b VECTOR(16384); END
| ; SELECT FROM
< error
| 1074 42000 Column length too big for column 'b' (max = 16383); use BLOB or TEXT instead
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) COLLATE ascii_bin)
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) COLLATE ascii_bin)
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) COLLATE 'utf8mb4_bin')
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) COLLATE 'utf8mb4_bin')
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) COLLATE binary)
< error
| 1064 42000 …near 'binary)' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) COLLATE utf8mb4_bin NOT NULL)
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) COLLATE utf8mb4_bin NOT NULL)
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a `VECTOR`(3) COLLATE utf8mb4_bin)
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a `VECTOR`(3) COLLATE utf8mb4_bin)
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| ALTER TABLE t CHANGE a b VECTOR(3) COLLATE utf8mb4_bin
< error
| 1063 42000 Incorrect column specifier for column 'b'
> exec ^
: string
| ALTER TABLE t CHANGE a b VECTOR(3) COLLATE utf8mb4_bin
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'b'
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) AUTO_INCREMENT)
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) AUTO_INCREMENT)
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) AUTO_INCREMENT PRIMARY KEY)
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) AUTO_INCREMENT PRIMARY KEY)
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) PRIMARY KEY AUTO_INCREMENT)
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) PRIMARY KEY AUTO_INCREMENT)
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) SERIAL DEFAULT VALUE)
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) SERIAL DEFAULT VALUE)
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a `VECTOR`(3) AUTO_INCREMENT)
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a `VECTOR`(3) AUTO_INCREMENT)
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| ALTER TABLE t MODIFY a VECTOR(3) AUTO_INCREMENT
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| ALTER TABLE t MODIFY a VECTOR(3) AUTO_INCREMENT
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) ON UPDATE CURRENT_TIMESTAMP)
< error
| 1294 HY000 Invalid ON UPDATE clause for 'a' column
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) ON UPDATE CURRENT_TIMESTAMP)
| ; SELECT FROM
< error
| 1294 HY000 Invalid ON UPDATE clause for 'a' column
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) ON UPDATE NOW())
< error
| 1294 HY000 Invalid ON UPDATE clause for 'a' column
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) ON UPDATE NOW())
| ; SELECT FROM
< error
| 1294 HY000 Invalid ON UPDATE clause for 'a' column
> exec ^
: string
| CREATE TABLE t (a `VECTOR`(3) ON UPDATE CURRENT_TIMESTAMP)
< error
| 1294 HY000 Invalid ON UPDATE clause for 'a' column
> exec ^
: string
| CREATE TABLE t (a `VECTOR`(3) ON UPDATE CURRENT_TIMESTAMP)
| ; SELECT FROM
< error
| 1294 HY000 Invalid ON UPDATE clause for 'a' column
> exec ^
: string
| ALTER TABLE t ADD COLUMN a VECTOR(3) ON UPDATE CURRENT_TIMESTAMP
< error
| 1294 HY000 Invalid ON UPDATE clause for 'a' column
> exec ^
: string
| ALTER TABLE t ADD COLUMN a VECTOR(3) ON UPDATE CURRENT_TIMESTAMP
| ; SELECT FROM
< error
| 1294 HY000 Invalid ON UPDATE clause for 'a' column
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) DEFAULT 1)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) INVISIBLE)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) UNIQUE)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) GENERATED ALWAYS AS (x))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VECTOR(16383))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VECTOR(16384))
< error
| 1074 42000 Column length too big for column 'a' (max = 16383); use BLOB or TEXT instead
> exec ^
: string
| CREATE TABLE t (a VECTOR(16384))
| ; SELECT FROM
< error
| 1074 42000 Column length too big for column 'a' (max = 16383); use BLOB or TEXT instead
> exec ^
: string
| CREATE TABLE t (a VECTOR(99999999999999999999))
< error
| 1074 42000 Column length too big for column 'a' (max = 16383); use BLOB or TEXT instead
> exec ^
: string
| CREATE TABLE t (a VECTOR(99999999999999999999))
| ; SELECT FROM
< error
| 1074 42000 Column length too big for column 'a' (max = 16383); use BLOB or TEXT instead
> exec ^
: string
| CREATE TABLE t (a VECTOR<FLOAT>)
< error
| 1064 42000 …near '<FLOAT>)' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR<FLOAT>(3))
< error
| 1064 42000 …near '<FLOAT>(3))' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(5)<FLOAT>)
< error
| 1064 42000 …near '<FLOAT>)' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR())
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(3,2))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VECTOR('3'))
< error
| 1064 42000 …near ''3'))' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(b))
< error
| 1064 42000 …near 'b))' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(1.5))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VECTOR(0x3))
< error
| 1064 42000 …near '0x3))' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(-1))
< error
| 1064 42000 …near '-1))' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(3)
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(3, FLOAT))
< error
| 1064 42000 …near 'FLOAT))' at line 1
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
| 1064 42000 …near 'UNSIGNED)' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) ZEROFILL)
< error
| 1064 42000 …near 'ZEROFILL)' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) BINARY)
< error
| 1064 42000 …near 'BINARY)' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) CHARSET utf8mb4)
< error
| 1064 42000 …near 'CHARSET utf8mb4)' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) CHARACTER SET utf8mb4)
< error
| 1064 42000 …near 'CHARACTER SET utf8mb4)' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) COLLATE utf8mb4_bin)
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) COLLATE utf8mb4_bin)
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) PRIMARY KEY)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) ARRAY)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) GENERATED ALWAYS AS ('[1]'))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) AS ('[1]') STORED)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a `VECTOR`)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| CREATE TABLE t (a `VECTOR`(3))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a vector(3))
< ok 0 0
> exec ^
: string
| CREATE TEMPORARY TABLE tt (a VECTOR(3))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX (v))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR KEY (v))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR (v))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR i (v))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR KEY i (v))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX IF NOT EXISTS i (v))
< error
| 1064 42000 …near '(v))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i USING HNSW (v))
< error
| 1064 42000 …near 'USING HNSW (v))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) USING HNSW)
< error
| 1064 42000 …near 'USING HNSW)' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) USING BTREE)
< error
| 1064 42000 …near 'USING BTREE)' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) TYPE HNSW)
< error
| 1064 42000 …near 'TYPE HNSW)' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX ((VEC_COSINE_DISTANCE(v))))
< error
| 1064 42000 …near '(VEC_COSINE_DISTANCE(v))))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX (v, w))
< error
| 1064 42000 …near ' w))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX (v(3)))
< error
| 1064 42000 …near '(3)))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX (v DESC))
< error
| 1064 42000 …near 'DESC))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M=6)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M 6)
< error
| 1064 42000 …near '6)' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M=abc)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M='6')
< ok 0 0
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M=6 DISTANCE=euclidean)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) DISTANCE=cosine)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) DISTANCE=foo)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) DISTANCE='cosine')
< ok 0 0
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) DISTANCE=default)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) M=default)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) COMMENT 'c')
< ok 0 0
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) VISIBLE)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) IGNORED)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) NOT IGNORED)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) KEY_BLOCK_SIZE=1)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH PARSER p)
< error
| 1064 42000 …near 'WITH PARSER p)' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) FOO=1)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=l2, type=hnsw, lib=vsag))
< error
| 1064 42000 …near 'WITH (distance=l2, type=hnsw, lib=vsag))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=l2, type=hnsw, lib=vsag, m=16, ef_construction=200, ef_search=64))
< error
| 1064 42000 …near 'WITH (distance=l2, type=hnsw, lib=vsag, m=16, ef_construction=200, ef_search=...' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR KEY i (v) WITH (distance=l2, type=hnsw, lib=vsag))
< error
| 1064 42000 …near 'WITH (distance=l2, type=hnsw, lib=vsag))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=cosine, type=hnsw))
< error
| 1064 42000 …near 'WITH (distance=cosine, type=hnsw))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=inner_product, type=hnsw_sq, lib=vsag))
< error
| 1064 42000 …near 'WITH (distance=inner_product, type=hnsw_sq, lib=vsag))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=l2, type=ivf_flat))
< error
| 1064 42000 …near 'WITH (distance=l2, type=ivf_flat))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH ())
< error
| 1064 42000 …near 'WITH ())' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=foo, type=bar))
< error
| 1064 42000 …near 'WITH (distance=foo, type=bar))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (foo=bar))
< error
| 1064 42000 …near 'WITH (foo=bar))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=l2 type=hnsw))
< error
| 1064 42000 …near 'WITH (distance=l2 type=hnsw))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH distance=l2)
< error
| 1064 42000 …near 'WITH distance=l2)' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX i (v) WITH (distance=l2, type=hnsw, lib=vsag) COMMENT 'c')
< error
| 1064 42000 …near 'WITH (distance=l2, type=hnsw, lib=vsag) COMMENT 'c')' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, VECTOR INDEX (v) WITH (distance=l2, type=hnsw, lib=vsag))
< error
| 1064 42000 …near 'WITH (distance=l2, type=hnsw, lib=vsag))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, CONSTRAINT c VECTOR INDEX (v))
< error
| 1064 42000 …near 'VECTOR INDEX (v))' at line 1
> exec ^
: string
| CREATE TABLE t (v VECTOR(3) NOT NULL, CONSTRAINT VECTOR INDEX (v))
< error
| 1064 42000 …near 'VECTOR INDEX (v))' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v)
< ok 0 0
> exec ^
: string
| CREATE VECTOR KEY i ON t (v)
< error
| 1064 42000 …near 'KEY i ON t (v)' at line 1
> exec ^
: string
| CREATE VECTOR i ON t (v)
< error
| 1064 42000 …near 'i ON t (v)' at line 1
> exec ^
: string
| CREATE VECTOR INDEX ON t (v)
< error
| 1064 42000 …near 'ON t (v)' at line 1
> exec ^
: string
| CREATE VECTOR INDEX
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t ()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| CREATE VECTOR
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| CREATE OR REPLACE VECTOR INDEX i ON t (v)
< ok 0 0
> exec ^
: string
| CREATE VECTOR INDEX IF NOT EXISTS i ON t (v)
< ok 0 0
> exec ^
: string
| CREATE UNIQUE VECTOR INDEX i ON t (v)
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v)' at line 1
> exec ^
: string
| CREATE VECTOR UNIQUE INDEX i ON t (v)
< error
| 1064 42000 …near 'UNIQUE INDEX i ON t (v)' at line 1
> exec ^
: string
| CREATE TEMPORARY VECTOR INDEX i ON t (v)
< error
| 1064 42000 …near 'VECTOR INDEX i ON t (v)' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i USING HNSW ON t (v)
< error
| 1064 42000 …near 'USING HNSW ON t (v)' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i TYPE HNSW ON t (v)
< error
| 1064 42000 …near 'TYPE HNSW ON t (v)' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i USING BTREE ON t (v)
< error
| 1064 42000 …near 'USING BTREE ON t (v)' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) USING HNSW
< error
| 1064 42000 …near 'USING HNSW' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) USING BTREE
< error
| 1064 42000 …near 'USING BTREE' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t ((VEC_COSINE_DISTANCE(v)))
< error
| 1064 42000 …near '(VEC_COSINE_DISTANCE(v)))' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v, w)
< error
| 1064 42000 …near ' w)' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v(3))
< error
| 1064 42000 …near '(3))' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v DESC)
< error
| 1064 42000 …near 'DESC)' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) M=6
< ok 0 0
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) M=6 DISTANCE=cosine
< ok 0 0
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) DISTANCE=euclidean M=6
< ok 0 0
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) M=abc
< ok 0 0
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) DISTANCE=foo
< ok 0 0
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) COMMENT 'c'
< ok 0 0
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) ALGORITHM=DEFAULT
< ok 0 0
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) ALGORITHM=INPLACE LOCK=NONE
< ok 0 0
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) LOCK=NONE
< ok 0 0
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) WAIT 3
< ok 0 0
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) NOWAIT
< ok 0 0
> exec ^
: string
| CREATE VECTOR INDEX i WAIT 3 ON t (v)
< error
| 1064 42000 …near 'WAIT 3 ON t (v)' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) M=6 ALGORITHM=DEFAULT
< ok 0 0
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) WITH (distance=l2, type=hnsw, lib=vsag)
< error
| 1064 42000 …near 'WITH (distance=l2, type=hnsw, lib=vsag)' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) WITH (distance=l2, type=hnsw, lib=vsag) COMMENT 'c'
< error
| 1064 42000 …near 'WITH (distance=l2, type=hnsw, lib=vsag) COMMENT 'c'' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) WITH (distance=cosine, type=hnsw)
< error
| 1064 42000 …near 'WITH (distance=cosine, type=hnsw)' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) WITH (foo=bar)
< error
| 1064 42000 …near 'WITH (foo=bar)' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) WITH ()
< error
| 1064 42000 …near 'WITH ()' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) WITH (distance=l2,)
< error
| 1064 42000 …near 'WITH (distance=l2,)' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) WITH (distance=l2) ALGORITHM=INPLACE
< error
| 1064 42000 …near 'WITH (distance=l2) ALGORITHM=INPLACE' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) SELECT 1
< error
| 1064 42000 …near 'SELECT 1' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) x
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v),
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| CREATE VECTOR INDEX vector ON vector (vector)
< error
| 1064 42000 …near 'vector ON vector (vector)' at line 1
> exec ^
: string
| CREATE VECTOR INDEX `index` ON t (v)
< ok 0 0
> exec ^
: string
| CREATE VECTOR INDEX select ON t (v)
< error
| 1064 42000 …near 'select ON t (v)' at line 1
> exec ^
: string
| CREATE TABLE vector (vector INT)
< error
| 1064 42000 …near 'vector (vector INT)' at line 1
> exec ^
: string
| CREATE TABLE vector (vector vector(3))
< error
| 1064 42000 …near 'vector (vector vector(3))' at line 1
> exec ^
: string
| CREATE TABLE t (vector)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| CREATE TABLE t (vector index)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| CREATE TABLE t (vector index INT)
< error
| 1064 42000 …near 'INT)' at line 1
> exec ^
: string
| CREATE TABLE t (vector INT, INDEX (vector))
< error
| 1064 42000 …near 'INT, INDEX (vector))' at line 1
> exec ^
: string
| SELECT vector FROM vector
< error
| 1064 42000 …near 'vector FROM vector' at line 1
> exec ^
: string
| SELECT 1 vector
< error
| 1064 42000 …near 'vector' at line 1
> exec ^
: string
| SELECT vector(3)
< error
| 1064 42000 …near 'vector(3)' at line 1
> exec ^
: string
| SELECT vector.vector FROM vector
< error
| 1064 42000 …near 'vector' at line 1
> exec ^
: string
| CREATE INDEX vector ON t (vector)
< error
| 1064 42000 …near 'vector ON t (vector)' at line 1
> exec ^
: string
| DROP INDEX vector ON t
< error
| 1064 42000 …near 'vector ON t' at line 1
> exec ^
: string
| SELECT CAST('[1]' AS VECTOR)
< error
| 4162 HY000 Operator does not exist: 'CAST(expr AS vector)'
> exec ^
: string
| SELECT CAST('[1]' AS VECTOR)
| ; SELECT FROM
< error
| 4162 HY000 Operator does not exist: 'CAST(expr AS vector)'
> exec ^
: string
| SELECT CAST('[1]' AS VECTOR(3))
< error
| 1064 42000 …near '(3))' at line 1
> exec ^
: string
| SELECT CAST('[1]' AS VECTOR(3) ARRAY)
< error
| 1064 42000 …near '(3) ARRAY)' at line 1
> exec ^
: string
| SELECT CONVERT('[1]', VECTOR)
< error
| 4162 HY000 Operator does not exist: 'CAST(expr AS vector)'
> exec ^
: string
| SELECT CONVERT('[1]', VECTOR)
| ; SELECT FROM
< error
| 4162 HY000 Operator does not exist: 'CAST(expr AS vector)'
> exec ^
: string
| SELECT CONVERT('[1]', VECTOR(3))
< error
| 1064 42000 …near '(3))' at line 1
> exec ^
: string
| ALTER TABLE t ADD COLUMN a VECTOR(3)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD a VECTOR(3)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD COLUMN (a VECTOR(3), b INT)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD vector vector(3)
< error
| 1064 42000 …near 'vector(3)' at line 1
> exec ^
: string
| ALTER TABLE t ADD vector INT
< error
| 1064 42000 …near 'INT' at line 1
> exec ^
: string
| ALTER TABLE t ADD COLUMN vector INT
< error
| 1064 42000 …near 'vector INT' at line 1
> exec ^
: string
| ALTER TABLE t ADD VECTOR
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| ALTER TABLE t ADD VECTOR (a)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD VECTOR KEY (a)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD VECTOR KEY i (a)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX (a)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX i (a)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX IF NOT EXISTS i (a)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD VECTOR i (a)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX i (a) M=6 DISTANCE=cosine
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX i (a) WITH (distance=l2, type=hnsw, lib=vsag)
< error
| 1064 42000 …near 'WITH (distance=l2, type=hnsw, lib=vsag)' at line 1
> exec ^
: string
| ALTER TABLE t ADD VECTOR KEY i (a) WITH (distance=l2, type=hnsw, lib=vsag)
< error
| 1064 42000 …near 'WITH (distance=l2, type=hnsw, lib=vsag)' at line 1
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX i (a) USING HNSW
< error
| 1064 42000 …near 'USING HNSW' at line 1
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX ((a))
< error
| 1064 42000 …near '(a))' at line 1
> exec ^
: string
| ALTER TABLE t ADD CONSTRAINT c VECTOR INDEX (a)
< error
| 1064 42000 …near 'VECTOR INDEX (a)' at line 1
> exec ^
: string
| ALTER TABLE t MODIFY a VECTOR(3)
< ok 0 0
> exec ^
: string
| ALTER TABLE t MODIFY COLUMN a VECTOR(3) NOT NULL
< ok 0 0
> exec ^
: string
| ALTER TABLE t MODIFY a VECTOR
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| ALTER TABLE t CHANGE a b VECTOR(3)
< ok 0 0
> exec ^
: string
| ALTER TABLE t DROP INDEX i
< ok 0 0
> exec ^
: string
| SELECT VEC_FROMTEXT('[1]'), VEC_TOTEXT(a), VEC_DISTANCE_COSINE(a,a), VEC_DISTANCE_EUCLIDEAN(a,a), VEC_DISTANCE(a,a) FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
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
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT l2_distance(a,'[1]'), cosine_distance(a,'[1]'), inner_product(a,'[1]') FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE INDEX i ON t (a) FOO=1
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) M=6
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) FOO=1)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) VISIBLE)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) IGNORED)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, FULLTEXT i (a) FOO=1)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, FULLTEXT i (a) VISIBLE)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, PRIMARY KEY (a) FOO=1)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD INDEX i (a) FOO=1
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD FULLTEXT (a) FOO=1
< ok 0 0
> exec ^
: string
| SELECT CAST('x' AS `VECTOR`)
< error
| 4162 HY000 Operator does not exist: 'CAST(expr AS vector)'
> exec ^
: string
| SELECT CAST('x' AS `VECTOR`)
| ; SELECT FROM
< error
| 4162 HY000 Operator does not exist: 'CAST(expr AS vector)'
> exec ^
: string
| SELECT CAST('x' AS INET6)
< ok 0 0
> exec ^
: string
| SELECT CAST('x' AS vector)
< error
| 4162 HY000 Operator does not exist: 'CAST(expr AS vector)'
> exec ^
: string
| SELECT CAST('x' AS vector)
| ; SELECT FROM
< error
| 4162 HY000 Operator does not exist: 'CAST(expr AS vector)'
> exec ^
: string
| CREATE TABLE t (a vector(3,))
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a vector(,3))
< error
| 1064 42000 …near '3))' at line 1
> exec ^
: string
| CREATE TABLE t (a vector(3,2,1))
< error
| 1064 42000 …near '1))' at line 1
> exec ^
: string
| CREATE TABLE t (a vector(1e5))
< error
| 1064 42000 …near '1e5))' at line 1
> exec ^
: string
| CREATE TABLE t (a vector(1.5,2))
< error
| 1064 42000 …near '2))' at line 1
> exec ^
: string
| CREATE TABLE t (a vector(3) NOT NULL, b vector(4) AUTO_INCREMENT)
< error
| 1063 42000 Incorrect column specifier for column 'b'
> exec ^
: string
| CREATE TABLE t (a vector(3) NOT NULL, b vector(4) AUTO_INCREMENT)
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'b'
> exec ^
: string
| CREATE TABLE t (v vector(3), VECTOR KEY `index` (v))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (v vector(3), VECTOR `index` (v))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (v vector(3), VECTOR KEY i (v, ))
< error
| 1064 42000 …near ' ))' at line 1
> exec ^
: string
| CREATE TABLE t (v vector(3), VECTOR KEY i ())
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| CREATE TABLE t (v vector(3), VECTOR KEY i (v) M=)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| CREATE TABLE t (v vector(3), VECTOR KEY i (v) M)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| CREATE TABLE t (v vector(3), VECTOR KEY i (v) M=-1)
< error
| 1064 42000 …near '-1)' at line 1
> exec ^
: string
| CREATE TABLE t (v vector(3), VECTOR KEY i (v) M=1.5)
< error
| 1064 42000 Only integers allowed as number here near '1.5)' at line 1
> exec ^
: string
| CREATE TABLE t (v vector(3), VECTOR KEY i (v) M=6, DISTANCE=cosine)
< error
| 1064 42000 …near '=cosine)' at line 1
> exec ^
: string
| CREATE TABLE t (v vector(3), VECTOR KEY i (v) `M`=6)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (v vector(3), VECTOR KEY i (v) select=6)
< error
| 1064 42000 …near 'select=6)' at line 1
> exec ^
: string
| CREATE TABLE t (v vector(3), VECTOR KEY i (v) COMMENT 'c' M=6)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (v vector(3), VECTOR KEY i (v) M=6 COMMENT 'c')
< ok 0 0
> exec ^
: string
| CREATE TABLE t (v vector(3), VECTOR KEY i (v) M=6 M=7)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (v vector(3), VECTOR KEY i (v) WITHOUT OVERLAPS)
< error
| 1064 42000 …near 'WITHOUT OVERLAPS)' at line 1
> exec ^
: string
| CREATE TABLE t (v vector(3), VECTOR KEY i (v) CLUSTERING=YES)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (v vector(3), VECTOR KEY i (v) ENGINE_ATTRIBUTE='x')
< ok 0 0
> exec ^
: string
| CREATE TABLE t (v vector(3), VECTOR KEY i (v) GLOBAL)
< error
| 1064 42000 …near 'GLOBAL)' at line 1
> exec ^
: string
| CREATE TABLE t (v vector(3), UNIQUE VECTOR KEY i (v))
< error
| 1064 42000 …near 'VECTOR KEY i (v))' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) M=
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) M=6, DISTANCE=cosine
< error
| 1064 42000 …near ' DISTANCE=cosine' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) COMMENT 'c' M=6
< ok 0 0
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) M=6 COMMENT 'c'
< ok 0 0
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) M=6 LOCK=NONE M=7
< error
| 1064 42000 …near 'M=7' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) IGNORED
< ok 0 0
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) NOT IGNORED
< ok 0 0
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) VISIBLE
< ok 0 0
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) KEY_BLOCK_SIZE=1
< ok 0 0
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) WITH PARSER p
< error
| 1064 42000 …near 'WITH PARSER p' at line 1
> exec ^
: string
| CREATE VECTOR INDEX `vector` ON `vector` (`vector`)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD VECTOR i (a) M=6, ADD COLUMN b INT
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX i (a, b)
< error
| 1064 42000 …near ' b)' at line 1
> exec ^
: string
| ALTER TABLE t ADD VECTOR i USING HNSW (a)
< error
| 1064 42000 …near 'USING HNSW (a)' at line 1
> exec ^
: string
| ALTER TABLE t ADD VECTOR IF NOT EXISTS i (a)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD VECTOR KEY IF NOT EXISTS i (a)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD VECTOR IF NOT EXISTS (a)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD `vector` INT
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD `vector` vector(3)
< ok 0 0
> exec ^
: string
| CREATE TABLE `vector` (`vector` vector(3))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (`vector` INT, `vector2` vector(3), VECTOR KEY `vector` (`vector2`))
< ok 0 0
> exec ^
: string
| CREATE TABLE vector2 (vector2 INT)
< ok 0 0
> exec ^
: string
| SELECT 1 AS vector
< error
| 1064 42000 …near 'vector' at line 1
> exec ^
: string
| SELECT `vector` FROM t vector
< error
| 1064 42000 …near 'vector' at line 1
> exec ^
: string
| SELECT a FROM vector
< error
| 1064 42000 …near 'vector' at line 1
> exec ^
: string
| INSERT INTO vector VALUES (1)
< error
| 1064 42000 …near 'vector VALUES (1)' at line 1
> exec ^
: string
| UPDATE t SET vector = 1
< error
| 1064 42000 …near 'vector = 1' at line 1
> exec ^
: string
| DROP TABLE vector
< error
| 1064 42000 …near 'vector' at line 1
> exec ^
: string
| DROP TABLE IF EXISTS vector
< error
| 1064 42000 …near 'vector' at line 1
> exec ^
: string
| ALTER TABLE vector ADD a INT
< error
| 1064 42000 …near 'vector ADD a INT' at line 1
> exec ^
: string
| ALTER TABLE t DROP vector
< error
| 1064 42000 …near 'vector' at line 1
> exec ^
: string
| ALTER TABLE t DROP COLUMN vector
< error
| 1064 42000 …near 'vector' at line 1
> exec ^
: string
| ALTER TABLE t DROP INDEX vector
< error
| 1064 42000 …near 'vector' at line 1
> exec ^
: string
| ALTER TABLE t RENAME INDEX a TO vector
< error
| 1064 42000 …near 'vector' at line 1
> exec ^
: string
| SET @vector = 1
< ok 0 0
> exec ^
: string
| SELECT @vector
< ok 0 0
> exec ^
: string
| SELECT t.vector FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT t.vector FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| USE vector
< error
| 1064 42000 …near 'vector' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) ENGINE=1)
< error
| 1064 42000 …near 'ENGINE=1)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) STATUS=1)
< error
| 1064 42000 …near 'STATUS=1)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) USER=1)
< error
| 1064 42000 …near 'USER=1)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) KEY=1)
< error
| 1064 42000 …near 'KEY=1)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) COMMENT=1)
< error
| 1064 42000 …near '=1)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) COMMENT='c')
< error
| 1064 42000 …near '='c')' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) TYPE=1)
< error
| 1064 42000 …near '=1)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) ALGORITHM=1)
< error
| 1064 42000 …near 'ALGORITHM=1)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) LOCK=NONE)
< error
| 1064 42000 …near 'LOCK=NONE)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) VISIBLE=1)
< error
| 1064 42000 …near '=1)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) IGNORED=1)
< error
| 1064 42000 …near '=1)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) TIMESTAMP=1)
< error
| 1064 42000 …near 'TIMESTAMP=1)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) FILE=1)
< error
| 1064 42000 …near 'FILE=1)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) FORMAT=1)
< error
| 1064 42000 …near 'FORMAT=1)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) NAME=1)
< error
| 1064 42000 …near 'NAME=1)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) SIMPLE=1)
< error
| 1064 42000 …near 'SIMPLE=1)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) M:=6)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) M=0x10)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) M=b'1')
< error
| 1064 42000 …near 'b'1')' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) M=TRUE)
< error
| 1064 42000 …near 'TRUE)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) M=NULL)
< error
| 1064 42000 …near 'NULL)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) M=1 N=2)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) M=a.b)
< error
| 1064 42000 …near '.b)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) M=select)
< error
| 1064 42000 …near 'select)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) M=status)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) M=99999999999999999999)
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) M=`x`)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) M=N'x')
< error
| 1064 42000 …near 'N'x')' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) `select`=1)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) 'x'=1)
< error
| 1064 42000 …near ''x'=1)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) M=1, b INT)
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) M=1 ALGORITHM=DEFAULT
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) ALGORITHM=DEFAULT M=1
< error
| 1064 42000 …near 'M=1' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a) M=1 LOCK=NONE
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) USING BTREE M=1
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) M=1 USING BTREE
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD INDEX i (a) M=1, ALGORITHM=DEFAULT
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD INDEX i (a) M=1 ALGORITHM=DEFAULT
< error
| 1064 42000 …near 'ALGORITHM=DEFAULT' at line 1
> exec ^
: string
| ALTER TABLE t ADD INDEX i (a) ALGORITHM=DEFAULT
< error
| 1064 42000 …near 'ALGORITHM=DEFAULT' at line 1
> exec ^
: string
| CREATE TABLE t (a vector(3), VECTOR i (a ASC))
< error
| 1064 42000 …near 'ASC))' at line 1
> exec ^
: string
| CREATE TABLE t (a vector(3), VECTOR i (`a`))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a vector(3), VECTOR i (t.a))
< error
| 1064 42000 …near '.a))' at line 1
> exec ^
: string
| CREATE TABLE t (a vector(3), VECTOR INDEX select (a))
< error
| 1064 42000 …near 'select (a))' at line 1
> exec ^
: string
| CREATE TABLE t (a vector(3), VECTOR INDEX vector (a))
< error
| 1064 42000 …near 'vector (a))' at line 1
> exec ^
: string
| CREATE TABLE t (a `vector`(3,2))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a `vector`(1.5))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a `vector`())
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a vector (3))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a `INET6`)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a `UUID`)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (`vector` vector(3), VECTOR KEY (`vector`))
< ok 0 0
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) WAIT 3 M=6
< ok 0 0
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) M=6 WAIT 3
< error
| 1064 42000 …near 'WAIT 3' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) NOWAIT M=6
< ok 0 0
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) M=6 NOWAIT
< error
| 1064 42000 …near 'NOWAIT' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) WAIT 3 ALGORITHM=DEFAULT
< ok 0 0
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) ALGORITHM=DEFAULT WAIT 3
< error
| 1064 42000 …near 'WAIT 3' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) WAIT
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| CREATE VECTOR INDEX i ON t (v) WAIT a
< error
| 1064 42000 …near 'a' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, INDEX IF NOT EXISTS i (a))
< error
| 1064 42000 …near '(a))' at line 1
> exec ^
: string
| ALTER TABLE t ADD (VECTOR INDEX IF NOT EXISTS i (a))
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD (VECTOR INDEX i (a))
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD COLUMN (b VECTOR(3) NOT NULL, VECTOR KEY (b))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) M=X'10')
< error
| 1064 42000 …near 'X'10')' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) M=X'0011')
< error
| 1064 42000 …near 'X'0011')' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) M=vector)
< error
| 1064 42000 …near 'vector)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, VECTOR i (a) M=vector)
< error
| 1064 42000 …near 'vector)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, VECTOR)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, VECTOR KEY)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, VECTOR KEY i)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| CREATE TEMPORARY TABLE tt (v vector(3) not null, VECTOR KEY (v))
< ok 0 0
> exec ^
: string
| CREATE FUNCTION f(a VECTOR(3)) RETURNS VECTOR(3) RETURN a
< ok 0 0
> exec ^
: string
| CREATE PROCEDURE p(a VECTOR(3)) BEGIN DECLARE b VECTOR(3); END
< ok 0 0
> exec ^
: string
| SELECT * FROM JSON_TABLE('[]', '$[*]' COLUMNS (a VECTOR(3) PATH '$')) AS j
< error
| 1064 42000 …near 'VECTOR(3) PATH '$')) AS j' at line 1
