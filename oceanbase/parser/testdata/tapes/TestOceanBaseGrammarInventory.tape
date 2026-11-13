-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOceanBaseGrammarInventory default#1
> exec
| CREATE DATABASE IF NOT EXISTS obprobe
< ok 1 0
> exec
| USE obprobe
< ok 0 0
> exec
| DROP TABLE IF EXISTS obprobe_t
< ok 0 0
> exec
| CREATE TABLE obprobe_t (a INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLEGROUP obprobe_tg1 PARTITION BY HASH PARTITIONS 8
< ok 0 0
> exec
| ALTER TABLE obprobe_t SET TABLEGROUP obprobe_tg1
< ok 0 0
> exec
| CREATE TABLE obprobe_t2 (a INT) TABLEGROUP = obprobe_tg1
< ok 0 0
> exec
| CREATE OUTLINE obprobe_ol1 ON SELECT * FROM obprobe_t
< ok 0 0
> exec
| CREATE RESOURCE UNIT obprobe_u1 MAX_CPU = 2, MEMORY_SIZE = '2G'
< error
| 1227 42501 Access denied
> exec
| CREATE RESOURCE POOL obprobe_p1 UNIT = 'obprobe_u1', UNIT_NUM = 1
< error
| 1227 42501 Access denied
> exec
| CREATE EXTERNAL TABLE obprobe_ext (a INT) LOCATION = '/tmp/obprobe' FORMAT (TYPE = 'CSV')
< ok 0 0
> exec
| CREATE MATERIALIZED VIEW obprobe_mv1 REFRESH COMPLETE ON DEMAND AS SELECT a FROM obprobe_t
< ok 0 0
> exec
| CREATE TABLE obprobe_t3 (a INT) WITH COLUMN GROUP(all columns, each column)
< ok 0 0
> exec
| ALTER SYSTEM MAJOR FREEZE
< ok 0 0
> exec
| SELECT * FROM ((SELECT 1) d)
< ok 0 0
> exec
| SELECT * FROM obprobe_t, (LATERAL (SELECT 1) d)
< ok 0 0
> exec
| SELECT * FROM (JSON_TABLE('[1]', '$[*]' COLUMNS (v INT PATH '$')) AS jt)
< ok 0 0
> exec
| CREATE TABLE obprobe_oi (a INT PRIMARY KEY) ORGANIZATION INDEX
< ok 0 0
> exec
| CREATE TABLE obprobe_so (a INT) REPLICA_NUM = 1 BLOCK_SIZE = 16384 TABLET_SIZE = 134217728 PCTFREE = 0 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE
< ok 0 0
> exec
| CREATE TABLE obprobe_ai (a INT AUTO_INCREMENT PRIMARY KEY) AUTO_INCREMENT_MODE = 'ORDER'
< ok 0 0
> exec
| CREATE TABLE obprobe_li (a INT, KEY obprobe_ia (a) LOCAL)
< ok 0 0
> exec
| CREATE TABLE obprobe_gi (a INT, KEY obprobe_ia (a) GLOBAL)
< ok 0 0
> exec
| CREATE TABLE obprobe_ib (a INT, KEY obprobe_ia (a) BLOCK_SIZE 16384)
< ok 0 0
> exec
| CREATE TABLE obprobe_ibl (a INT, KEY obprobe_ia (a) BLOCK_SIZE 16384 LOCAL)
< ok 0 0
> exec
| ALTER TABLE obprobe_t ADD UNIQUE INDEX obprobe_uk (a) GLOBAL
< ok 0 0
> exec
| CREATE TABLE obprobe_fl (id INT PRIMARY KEY, c VARCHAR(100), FULLTEXT KEY obprobe_fk (c) LOCAL)
< ok 0 0
> exec
| CREATE TABLE obprobe_fg (id INT PRIMARY KEY, c VARCHAR(100), FULLTEXT KEY obprobe_fk (c) GLOBAL)
< error
| 1235 0A000 global fulltext index is not supported
> exec
| CREATE TABLE obprobe_sg (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0, SPATIAL KEY obprobe_sk (g) GLOBAL)
< ok 0 0
> exec
| CREATE TABLE obprobe_rs (id INT PRIMARY KEY, a INT, KEY obprobe_ka (a) LOCAL GLOBAL)
< ok 0 0
> exec
| ALTER TABLE obprobe_t ADD PRIMARY KEY (a) GLOBAL
< error
| 1068 42000 Multiple primary key defined
> exec
| CREATE TABLE obprobe_fp (id INT PRIMARY KEY, c VARCHAR(100), FULLTEXT KEY obprobe_fk (c) WITH PARSER space PARSER_PROPERTIES=(min_token_size=3,max_token_size=84) BLOCK_SIZE 16384)
< ok 0 0
> exec
| CREATE TABLE obprobe_sp (a INT, b INT) PARTITION BY RANGE COLUMNS(a) SUBPARTITION BY KEY(b) SUBPARTITION TEMPLATE (SUBPARTITION p0, SUBPARTITION p1) (PARTITION q0 VALUES LESS THAN (10))
< ok 0 0
> exec
| CREATE TABLE obprobe_rsp (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY RANGE(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES LESS THAN (5), SUBPARTITION s1 VALUES LESS THAN (MAXVALUE)) (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> exec
| CREATE TABLE obprobe_dp (a INT) DYNAMIC_PARTITION_POLICY = (ENABLE = TRUE, TIME_UNIT = 'DAY')
< error
| 1235 0A000 specify dynamic partition on non partitioned table is not supported
> exec
| CREATE TABLE obprobe_oh (a INT PRIMARY KEY) ORGANIZATION HEAP
< ok 0 0
> exec
| ALTER TABLE obprobe_t SET PCTFREE = 5
< ok 0 0
> exec
| ALTER TABLE obprobe_t SET DYNAMIC_PARTITION_POLICY (ENABLE = FALSE)
< error
| 1235 0A000 specify dynamic partition on non partitioned table is not supported
> exec
| ALTER TABLE obprobe_t RENAME COLUMN a TO b
< error
| 1054 42S22 Unknown column 'obprobe.obprobe_t.a' in 'field list'
> exec
| WITH c AS (SELECT 1 AS x) SELECT x FROM c
< ok 0 0
> exec
| SELECT ROW_NUMBER() OVER (ORDER BY a) FROM obprobe_t
< ok 0 0
> exec
| CREATE TABLE obprobe_t4 (a INT, CONSTRAINT obprobe_ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE obprobe_t5 (a INT, b INT GENERATED ALWAYS AS (a + 1) STORED)
< ok 0 0
> exec
| SELECT * FROM JSON_TABLE('[1]', '$[*]' COLUMNS (v INT PATH '$')) AS jt
< ok 0 0
> exec
| SELECT * FROM obprobe_t, LATERAL (SELECT 1) AS d
< ok 0 0
> exec
| CREATE TABLE obprobe_t6 (a INT INVISIBLE)
< error
| 4028 HY000 table must have at least one column that is not invisible
> exec
| CREATE TABLE obprobe_t7 (a INT, INDEX ((a + 1)))
< ok 0 0
> exec
| VALUES ROW(1), ROW(2)
< ok 0 0
> exec
| SELECT 1 INTERSECT SELECT 1
< ok 0 0
> exec
| SELECT 1 EXCEPT SELECT 2
< ok 0 0
> exec
| CREATE TABLE obprobe_so2 (a INT) COMMENT = 'x' CONNECTION 'y' PASSWORD = 'z'
< ok 0 0
> exec
| CREATE TABLE obprobe_cc (a INT COMMENT 'x' NOT NULL, b INT NOT NULL COMMENT "y" DEFAULT 1)
< ok 0 0
> exec
| CREATE TABLE obprobe_kc (a INT, PRIMARY KEY (a) COMMENT 'x', KEY obprobe_ia (a) COMMENT 'y' USING BTREE)
< ok 0 0
> exec
| ALTER TABLE obprobe_t ADD INDEX obprobe_ic (a) COMMENT 'x', ADD COLUMN obprobe_c INT
< error
| 1235 0A000 double table long running ddl on table required by materialized view is not supported
> exec
| CREATE PROCEDURE obprobe_p() COMMENT 'x' SELECT 1
< ok 0 0
> exec
| ALTER TABLE obprobe_t DROP INDEX IF EXISTS obprobe_ia
< error
| 1064 42000 …near 'IF EXISTS obprobe_ia' at line 1
> exec
| ALTER TABLE obprobe_t DROP COLUMN IF EXISTS a
< error
| 1064 42000 …near 'IF EXISTS a' at line 1
> exec
| DROP INDEX IF EXISTS obprobe_ia ON obprobe_t
< error
| 1064 42000 …near 'IF EXISTS obprobe_ia ON obprobe_t' at line 1
> exec
| CREATE TABLE obprobe_tv (a INT) COMMENT
< error
| 1064 42000 …near 'COMMENT' at line 1
> exec
| CREATE TABLE obprobe_tv (a INT) COMMENT =
< error
| 1064 42000 …near '=' at line 1
> exec
| CREATE TABLE obprobe_tv (a INT) CONNECTION
< error
| 1064 42000 …near 'CONNECTION' at line 1
> exec
| CREATE TABLE obprobe_tv (a INT) PASSWORD =
< error
| 1064 42000 …near '=' at line 1
> exec
| ALTER TABLE obprobe_t COMMENT
< error
| 1064 42000 …near 'COMMENT' at line 1
> exec
| CREATE TABLE obprobe_tv (a INT) KEY_BLOCK_SIZE
< error
| 1064 42000 …near 'KEY_BLOCK_SIZE' at line 1
> exec
| CREATE TABLE obprobe_tv (a INT COMMENT)
< error
| 1064 42000 …near ')' at line 1
> exec
| CREATE TABLE obprobe_tv (a INT COMMENT NOT NULL)
< error
| 1064 42000 …near 'NOT NULL)' at line 1
> exec
| ALTER TABLE obprobe_t MODIFY a INT COMMENT
< error
| 1064 42000 …near 'COMMENT' at line 1
> exec
| CREATE TABLE obprobe_tv (a INT, KEY obprobe_ia (a) COMMENT)
< error
| 1064 42000 …near ')' at line 1
> exec
| CREATE TABLE obprobe_tv (a INT, PRIMARY KEY (a) COMMENT, b INT)
< error
| 1064 42000 …near ', b INT)' at line 1
> exec
| CREATE TABLE obprobe_tv (a INT, KEY obprobe_ia (a) BLOCK_SIZE)
< error
| 1064 42000 …near ')' at line 1
> exec
| CREATE TABLE obprobe_tv (a INT, KEY obprobe_ia (a) USING)
< error
| 1064 42000 …near ')' at line 1
> exec
| ALTER TABLE obprobe_t ADD INDEX obprobe_ia (a) COMMENT
< error
| 1064 42000 …near 'COMMENT' at line 1
> exec
| CREATE PROCEDURE obprobe_pv() COMMENT SELECT 1
< error
| 1064 42000 …near 'SELECT' at line 1
> exec
| ALTER PROCEDURE obprobe_pv COMMENT
< error
| 1064 42000 …near 'COMMENT' at line 1
> exec
| ALTER EVENT obprobe_e COMMENT
< error
| 1064 42000 …near 'COMMENT' at line 1
> exec
| CREATE USER obprobe_u COMMENT
< error
| 1064 42000 …near 'COMMENT' at line 1
> exec
| ALTER USER obprobe_u ATTRIBUTE
< error
| 1064 42000 …near 'ATTRIBUTE' at line 1
> exec
| CREATE TABLESPACE obprobe_ts ADD DATAFILE 'obprobe_ts.ibd' COMMENT
< error
| 1064 42000 …near 'COMMENT' at line 1
> exec
| CREATE LOGFILE GROUP obprobe_lg ADD UNDOFILE 'obprobe_u.dat' COMMENT =
< error
| 1064 42000 …near '=' at line 1
> exec
| CREATE TABLE obprobe_bv1 (a INT, b INT) PARTITION BY HASH(a) SUBPARTITION BY HASH(b) SUBPARTITION TEMPLATE (SUBPARTITION s0 VALUES)
< error
| 1064 42000 …near 'VALUES)' at line 1
> exec
| CREATE TABLE obprobe_bv2 (a INT, b INT) PARTITION BY RANGE(a) SUBPARTITION BY HASH(b) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0 VALUES ENGINE=InnoDB))
< error
| 1064 42000 …near 'ENGINE=InnoDB))' at line 1
> exec
| CREATE TABLE obprobe_bv3 (a INT) PARTITION BY RANGE(a) (PARTITION p0 VALUES)
< error
| 1064 42000 …near ')' at line 1
> exec
| CREATE TABLE obprobe_bv4 (a INT) PARTITION BY HASH(a) (PARTITION p0 VALUES ENGINE = InnoDB)
< error
| 1064 42000 …near 'VALUES ENGINE = InnoDB)' at line 1
> exec
| DROP DATABASE IF EXISTS obprobe
< ok 46 0
