-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestIndexOptionFormsOceanBaseOracle default#1
> exec
| DROP DATABASE IF EXISTS obo
< ok 0 0
> exec
| DROP DATABASE IF EXISTS probe
< ok 0 0
> exec
| CREATE DATABASE obo
< ok 1 0
> exec
| CREATE DATABASE probe
< ok 1 0
> exec
| USE probe
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, c TEXT, g GEOMETRY NOT NULL)
< ok 0 0
> exec
| CREATE TABLE obo.ot (a INT, b INT)
< ok 0 0
> exec
| CREATE INDEX p1 ON t (a) PARALLEL 4
< ok 0 0
> exec
| CREATE INDEX p2 ON t (a) PARALLEL = 4
< ok 0 0
> exec
| CREATE INDEX p3 ON t (a) PARALLEL=1
< ok 0 0
> exec
| CREATE INDEX p4 ON t (a) PARALLEL a
< error
| 1064 42000 …near 'a' at line 1
> exec
| CREATE INDEX p5 ON t (a) PARALLEL 0
< error
| 1064 42000 …near 'PARALLEL 0' at line 1
> exec
| CREATE INDEX p6 ON t (a) PARALLEL
< error
| 1064 42000 …near 'PARALLEL' at line 1
> exec
| CREATE INDEX p7 ON t (a) PARALLEL 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec
| CREATE INDEX p8 ON t (a) PARALLEL '4'
< error
| 1064 42000 …near ''4'' at line 1
> exec
| CREATE INDEX p9 ON t (a) PARALLEL -1
< error
| 1064 42000 …near '-1' at line 1
> exec
| CREATE INDEX p10 ON t (a) PARALLEL =
< error
| 1064 42000 …near '=' at line 1
> exec
| CREATE INDEX p11 ON t (a) PARALLEL 2 PARALLEL 3
< ok 0 0
> exec
| CREATE INDEX p12 ON t (a) PARALLEL 2 COMMENT 'c' NOPARALLEL
< ok 0 0
> exec
| CREATE INDEX p13 ON t (a) NOPARALLEL
< ok 0 0
> exec
| CREATE INDEX p14 ON t (a) NOPARALLEL 4
< error
| 1064 42000 …near '4' at line 1
> exec
| CREATE INDEX p15 ON t (a) NOPARALLEL = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec
| CREATE INDEX `parallel` ON t (a)
< ok 0 0
> exec
| CREATE INDEX parallel ON t (a)
< error
| 1061 42000 Duplicate key name 'parallel'
> exec
| CREATE INDEX noparallel ON t (a)
< ok 0 0
> exec
| CREATE INDEX storing ON t (a)
< ok 0 0
> exec
| CREATE INDEX rowid ON t (a)
< ok 0 0
> exec
| CREATE TABLE parallel (parallel INT, noparallel INT, storing INT, rowid INT, data_table_id INT, index_table_id INT, virtual_column_id INT, max_used_part_id INT, storage_cache_policy INT)
< ok 0 0
> exec
| CREATE INDEX s1 ON t (a) STORING (b)
< ok 0 0
> exec
| CREATE INDEX s2 ON t (a) STORING (b, c)
< error
| 1167 42000 The used storage engine can't index column 'c'
> exec
| CREATE INDEX s3 ON t (a) STORING ()
< error
| 1064 42000 …near ')' at line 1
> exec
| CREATE INDEX s4 ON t (a) STORING b
< error
| 1064 42000 …near 'b' at line 1
> exec
| CREATE INDEX s5 ON t (a) STORING (`b`)
< ok 0 0
> exec
| CREATE INDEX s6 ON t (a) STORING (t.b)
< error
| 1064 42000 …near 't.b)' at line 1
> exec
| CREATE INDEX s7 ON t (a) STORING (nosuch)
< error
| 1054 42S22 Unknown column
> exec
| CREATE INDEX s8 ON t (a) STORING (b,)
< error
| 1064 42000 …near ')' at line 1
> exec
| CREATE INDEX s9 ON t (a) STORING (b DESC)
< error
| 1064 42000 …near 'DESC)' at line 1
> exec
| CREATE INDEX s10 ON t (a) STORING (b(3))
< error
| 1064 42000 …near '(3))' at line 1
> exec
| CREATE INDEX s11 ON t (a) STORING (b) STORING (c)
< error
| 1167 42000 The used storage engine can't index column 'c'
> exec
| CREATE INDEX s12 ON t (a) STORING
< error
| 1064 42000 …near 'STORING' at line 1
> exec
| CREATE INDEX s13 ON t (a) STORING (select)
< error
| 1064 42000 …near 'select)' at line 1
> exec
| CREATE INDEX s14 ON t (a) `STORING` (b)
< error
| 1064 42000 …near '`STORING` (b)' at line 1
> exec
| CREATE INDEX s15 ON t (a) `PARALLEL` 2
< error
| 1064 42000 …near '`PARALLEL` 2' at line 1
> exec
| CREATE INDEX s16 ON t (a) `NOPARALLEL`
< error
| 1064 42000 …near '`NOPARALLEL`' at line 1
> exec
| CREATE INDEX w1 ON t (a) WITH ROWID
< ok 0 0
> exec
| CREATE INDEX w2 ON t (a) WITH ROW
< error
| 1064 42000 …near 'ROW' at line 1
> exec
| CREATE INDEX w3 ON t (a) WITH
< error
| 1064 42000 …near 'WITH' at line 1
> exec
| CREATE INDEX i1 ON t (a) DATA_TABLE_ID = 1
< error
| 4179 HY000 create index with index_table_id/data_table_id not allowed
> exec
| CREATE INDEX i2 ON t (a) DATA_TABLE_ID 1
< error
| 4179 HY000 create index with index_table_id/data_table_id not allowed
> exec
| CREATE INDEX i3 ON t (a) DATA_TABLE_ID a
< error
| 1064 42000 …near 'a' at line 1
> exec
| CREATE INDEX i4 ON t (a) INDEX_TABLE_ID = 1
< error
| 4179 HY000 create index with index_table_id/data_table_id not allowed
> exec
| CREATE INDEX i5 ON t (a) VIRTUAL_COLUMN_ID = 1
< error
| 1064 42000 Parse error
> exec
| CREATE INDEX i6 ON t (a) MAX_USED_PART_ID = 1
< ok 0 0
> exec
| CREATE INDEX i7 ON t (a) INDEX_TABLE_ID
< error
| 1064 42000 …near 'INDEX_TABLE_ID' at line 1
> exec
| CREATE INDEX c1 ON t (a) STORAGE_CACHE_POLICY = (GLOBAL = 'hot')
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| CREATE INDEX c2 ON t (a) STORAGE_CACHE_POLICY (GLOBAL = 'hot')
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| CREATE INDEX c3 ON t (a) STORAGE_CACHE_POLICY = (GLOBAL = hot)
< error
| 1064 42000 …near 'hot)' at line 1
> exec
| CREATE INDEX c4 ON t (a) STORAGE_CACHE_POLICY = ()
< error
| 1064 42000 …near ')' at line 1
> exec
| CREATE INDEX c5 ON t (a) STORAGE_CACHE_POLICY = (BOUNDARY_COLUMN = a, HOT_RETENTION = 1 DAY)
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| CREATE INDEX c6 ON t (a) STORAGE_CACHE_POLICY = (BOUNDARY_COLUMN = a, BOUNDARY_COLUMN_UNIT = 's', HOT_RETENTION = 1 DAY)
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| CREATE INDEX c7 ON t (a) STORAGE_CACHE_POLICY = 'hot'
< error
| 1064 42000 …near ''hot'' at line 1
> exec
| CREATE TABLE n1 (a INT, b INT, KEY k (a) PARALLEL 2)
< ok 0 0
> exec
| CREATE TABLE n2 (a INT, b INT, KEY k (a) NOPARALLEL)
< ok 0 0
> exec
| CREATE TABLE n3 (a INT, b INT, KEY k (a) STORING (b))
< ok 0 0
> exec
| CREATE TABLE n4 (a INT, b INT, KEY k (a) WITH ROWID)
< ok 0 0
> exec
| CREATE TABLE n5 (a INT, b INT, UNIQUE KEY k (a) PARALLEL 2)
< ok 0 0
> exec
| CREATE TABLE n6 (a INT, b INT, UNIQUE KEY k (a) STORING (b))
< ok 0 0
> exec
| CREATE TABLE n7 (a INT, c TEXT, FULLTEXT KEY k (c) PARALLEL 2)
< ok 0 0
> exec
| CREATE TABLE n8 (a INT, c TEXT, FULLTEXT KEY k (c) STORING (a))
< ok 0 0
> exec
| CREATE TABLE n9 (a INT, g GEOMETRY NOT NULL, SPATIAL KEY k (g) PARALLEL 2)
< ok 0 0
> exec
| CREATE TABLE n10 (a INT, g GEOMETRY NOT NULL, SPATIAL KEY k (g) STORING (a))
< ok 0 0
> exec
| CREATE TABLE n11 (a INT, b INT, PRIMARY KEY (a) PARALLEL 2)
< error
| 1064 42000 …near 'PARALLEL 2)' at line 1
> exec
| CREATE TABLE n12 (a INT, b INT, PRIMARY KEY (a) STORING (b))
< error
| 1064 42000 …near 'STORING (b))' at line 1
> exec
| CREATE TABLE n13 (a INT, b INT, CONSTRAINT c PRIMARY KEY (a) NOPARALLEL)
< error
| 1064 42000 …near 'NOPARALLEL)' at line 1
> exec
| CREATE TABLE n14 (a INT, b INT, CONSTRAINT c UNIQUE KEY k (a) PARALLEL 2)
< ok 0 0
> exec
| CREATE TABLE n15 (a INT, b INT, CONSTRAINT c UNIQUE KEY k (a) STORING (b))
< ok 0 0
> exec
| CREATE TABLE n16 (a INT PARALLEL 2)
< error
| 1064 42000 …near 'PARALLEL 2)' at line 1
> exec
| ALTER TABLE t ADD INDEX a1 (a) PARALLEL 4
< ok 0 0
> exec
| ALTER TABLE t ADD INDEX a2 (a) NOPARALLEL
< ok 0 0
> exec
| ALTER TABLE t ADD INDEX a3 (a) STORING (b)
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE KEY a4 (a) STORING (b)
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE KEY a5 (a) PARALLEL 4
< ok 0 0
> exec
| ALTER TABLE t ADD FULLTEXT KEY a6 (c) PARALLEL 4
< ok 0 0
> exec
| ALTER TABLE t ADD CONSTRAINT c7 UNIQUE KEY a7 (a) PARALLEL 4
< ok 0 0
> exec
| ALTER TABLE t ADD PRIMARY KEY (a) PARALLEL 4
< ok 0 0
> exec
| ALTER TABLE t ADD CONSTRAINT c8 PRIMARY KEY (a) STORING (b)
< error
| 1068 42000 Multiple primary key defined
> exec
| ALTER TABLE t ADD SPATIAL KEY a9 (g) STORING (a)
< ok 0 0
> exec
| ALTER TABLE t ADD INDEX a11 (a) PARALLEL 2, ADD INDEX a12 (b) STORING (a)
< ok 0 0
> exec
| CREATE UNIQUE INDEX u1 ON t (a) STORING (b)
< ok 0 0
> exec
| CREATE UNIQUE INDEX u2 ON t (a) PARALLEL 2
< ok 0 0
> exec
| CREATE FULLTEXT INDEX f1 ON t (c) PARALLEL 2
< ok 0 0
> exec
| CREATE FULLTEXT INDEX f2 ON t (c) STORING (a)
< ok 0 0
> exec
| CREATE SPATIAL INDEX g1 ON t (g) PARALLEL 2
< error
| 1060 42S21 Duplicate column name
> exec
| CREATE SPATIAL INDEX g2 ON t (g) STORING (a)
< error
| 1060 42S21 Duplicate column name
> exec
| CREATE VECTOR INDEX v1 ON t (a) WITH (distance=l2, type=hnsw) PARALLEL 2
< error
| 7601 HY000 Column 'a' cannot be part of VECTOR index
> exec
| CREATE VECTOR INDEX v2 ON t (a) STORING (b) WITH (distance=l2, type=hnsw)
< error
| 7601 HY000 Column 'a' cannot be part of VECTOR index
> exec
| CREATE INDEX probe.q1 ON t (a)
< ok 0 0
> exec
| CREATE INDEX obo.q2 ON t (a)
< error
| 1235 0A000 Index name including database name is not supported
> exec
| CREATE INDEX nosuch.q3 ON t (a)
< error
| 1049 42000 Unknown database
> exec
| CREATE INDEX x.y.q4 ON t (a)
< error
| 1049 42000 Unknown database
> exec
| CREATE INDEX w.x.y.q5 ON t (a)
< error
| 1064 42000 …near '.q5 ON t (a)' at line 1
> exec
| CREATE INDEX .q6 ON t (a)
< error
| 1064 42000 …near '.q6 ON t (a)' at line 1
> exec
| CREATE INDEX probe.select ON t (a)
< ok 0 0
> exec
| CREATE INDEX select.q7 ON t (a)
< error
| 1049 42000 Unknown database
> exec
| CREATE INDEX `probe`.`q8` ON t (a)
< ok 0 0
> exec
| CREATE INDEX probe. q9 ON t (a)
< ok 0 0
> exec
| CREATE INDEX probe.q10 ON obo.ot (a)
< error
| 1235 0A000 Index name including database name is not supported
> exec
| CREATE INDEX obo.q11 ON obo.ot (a)
< ok 0 0
> exec
| CREATE INDEX IF NOT EXISTS probe.q12 ON t (a)
< ok 0 0
> exec
| CREATE UNIQUE INDEX probe.q13 ON t (a)
< ok 0 0
> exec
| CREATE FULLTEXT INDEX probe.q14 ON t (c)
< ok 0 0
> exec
| CREATE VECTOR INDEX probe.q15 ON t (a) WITH (distance=l2, type=hnsw)
< error
| 7601 HY000 Column 'a' cannot be part of VECTOR index
> exec
| CREATE INDEX probe.q16 USING BTREE ON t (a)
< ok 0 0
> exec
| CREATE INDEX q17@l ON t (a)
< ok 0 0
> exec
| CREATE INDEX probe.q18@l ON t (a)
< ok 0 0
> exec
| CREATE TABLE q19 (a INT, KEY probe.k (a))
< error
| 1064 42000 …near 'probe.k (a))' at line 1
> exec
| ALTER TABLE t ADD INDEX probe.q20 (a)
< error
| 1064 42000 …near 'probe.q20 (a)' at line 1
> exec
| DROP INDEX probe.q1 ON t
< error
| 1064 42000 …near 'probe.q1 ON t' at line 1
> exec
| CREATE TABLE m1 (a INT, CONSTRAINT c INDEX (a))
< error
| 1064 42000 …near 'INDEX (a))' at line 1
> exec
| CREATE TABLE m2 (a INT, CONSTRAINT c KEY (a))
< error
| 1064 42000 …near 'KEY (a))' at line 1
> exec
| CREATE TABLE m3 (a INT, CONSTRAINT INDEX (a))
< error
| 1064 42000 …near 'INDEX (a))' at line 1
> exec
| CREATE TABLE m4 (a INT, CONSTRAINT KEY (a))
< error
| 1064 42000 …near 'KEY (a))' at line 1
> exec
| CREATE TABLE m5 (a INT, CONSTRAINT c KEY k (a))
< error
| 1064 42000 …near 'KEY k (a))' at line 1
> exec
| CREATE TABLE m6 (c TEXT, CONSTRAINT x FULLTEXT (c))
< error
| 1064 42000 …near 'FULLTEXT (c))' at line 1
> exec
| CREATE TABLE m7 (c TEXT, CONSTRAINT FULLTEXT (c))
< error
| 1064 42000 …near 'FULLTEXT (c))' at line 1
> exec
| CREATE TABLE m8 (g GEOMETRY NOT NULL, CONSTRAINT x SPATIAL (g))
< error
| 1064 42000 …near 'SPATIAL (g))' at line 1
> exec
| CREATE TABLE m9 (g GEOMETRY NOT NULL, CONSTRAINT SPATIAL (g))
< error
| 1064 42000 …near 'SPATIAL (g))' at line 1
> exec
| CREATE TABLE m10 (a INT, CONSTRAINT VECTOR INDEX (a))
< error
| 1064 42000 …near 'INDEX (a))' at line 1
> exec
| CREATE TABLE m11 (a INT, CONSTRAINT x VECTOR INDEX (a))
< error
| 1064 42000 …near 'VECTOR INDEX (a))' at line 1
> exec
| CREATE TABLE m12 (a INT, CONSTRAINT x UNIQUE (a))
< ok 0 0
> exec
| CREATE TABLE m13 (a INT, CONSTRAINT UNIQUE (a))
< ok 0 0
> exec
| CREATE TABLE m14 (a INT, CONSTRAINT x PRIMARY KEY (a))
< ok 0 0
> exec
| CREATE TABLE m15 (a INT, CONSTRAINT x FOREIGN KEY (a) REFERENCES obo.ot (a))
< ok 0 0
> exec
| CREATE TABLE m16 (a INT, CONSTRAINT x CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE m17 (a INT, CONSTRAINT)
< error
| 1064 42000 …near ')' at line 1
> exec
| CREATE TABLE m18 (a INT, CONSTRAINT x)
< error
| 1064 42000 …near ')' at line 1
> exec
| ALTER TABLE t ADD CONSTRAINT x INDEX (a)
< error
| 1064 42000 …near 'INDEX (a)' at line 1
> exec
| ALTER TABLE t ADD CONSTRAINT x KEY k (a)
< error
| 1064 42000 …near 'KEY k (a)' at line 1
> exec
| ALTER TABLE t ADD CONSTRAINT INDEX (a)
< error
| 1064 42000 …near 'INDEX (a)' at line 1
> exec
| ALTER TABLE t ADD CONSTRAINT KEY (a)
< error
| 1064 42000 …near 'KEY (a)' at line 1
> exec
| ALTER TABLE t ADD CONSTRAINT x FULLTEXT (c)
< error
| 1064 42000 …near 'FULLTEXT (c)' at line 1
> exec
| ALTER TABLE t ADD CONSTRAINT x SPATIAL (g)
< error
| 1064 42000 …near 'SPATIAL (g)' at line 1
> exec
| ALTER TABLE t ADD CONSTRAINT FULLTEXT (c)
< error
| 1064 42000 …near 'FULLTEXT (c)' at line 1
> exec
| ALTER TABLE t ADD (CONSTRAINT x INDEX (a))
< error
| 1064 42000 …near 'INDEX (a))' at line 1
> exec
| ALTER TABLE t ADD COLUMN (d INT, CONSTRAINT x KEY (d))
< error
| 1064 42000 …near 'CONSTRAINT x KEY (d))' at line 1
> exec
| DROP DATABASE IF EXISTS obo
< ok 2 0
> exec
| DROP DATABASE IF EXISTS probe
< ok 93 0
