-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_OceanBaseIndexOptionRefusals/ALTER_TABLE default#1
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, v VARCHAR(10), KEY idx_c (c))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY (c) WITH (a = 1)
< error
| 1235 0A000 set index params in not vector index is not supported
> exec
| ALTER TABLE a ADD KEY (zz) DATA_TABLE_ID 1
< error
| 1235 0A000 specify data table id not for restore purpose in create index not supported
> exec
| ALTER TABLE a ADD KEY idx_c (c) INDEX_TABLE_ID 1
< error
| 1235 0A000 specify index table id not for restore purpose in create index not supported
> exec
| ALTER TABLE a ADD UNIQUE (c) STORAGE_CACHE_POLICY (GLOBAL = 'hot')
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
> exec
| ALTER TABLE a ADD b INT, ADD KEY (c) WITH (a = 1)
< error
| 1235 0A000 set index params in not vector index is not supported
> exec
| ALTER TABLE a ADD FULLTEXT (v) WITH (a = 1)
< error
| 1235 0A000 set index params in not vector index is not supported

== TestOracle_OceanBaseIndexOptionRefusals/CREATE_TABLE_and_CREATE_INDEX default#1
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, v VARCHAR(10), KEY idx_c (c))
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, KEY (b) WITH (a = 1))
< error
| 1235 0A000 set index params in not vector index is not supported
> exec
| CREATE TABLE t2 (b INT PRIMARY KEY, KEY (b) DATA_TABLE_ID 1)
< error
| 1235 0A000 specify data table id not for restore purpose in create index not supported
> exec
| CREATE INDEX i ON a (c) WITH (a = 1)
< error
| 1235 0A000 set index params in not vector index is not supported
> exec
| CREATE INDEX i ON a (c) DATA_TABLE_ID 1
< error
| 4179 HY000 create index with index_table_id/data_table_id not allowed
> exec
| CREATE INDEX i ON a (c) STORAGE_CACHE_POLICY (GLOBAL = 'hot')
< error
| 1235 0A000 storage cache policy is not supported in shared storage mode not supported
