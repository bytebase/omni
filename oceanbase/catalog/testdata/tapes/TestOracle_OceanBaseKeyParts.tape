-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_OceanBaseKeyParts/ALTER_TABLE_and_CREATE_INDEX default#1
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, v VARCHAR(10), KEY idx_c (c))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY idx_c (c, c)
< error
| 1060 42S21 Duplicate column name 'c'
> exec
| ALTER TABLE a ADD KEY (zz, zz)
< error
| 1060 42S21 Duplicate column name 'zz'
> exec
| ALTER TABLE a ADD KEY ((c))
< error
| 3762 HY000 Functional index on a column is not supported. Consider using a regular index instead.
> exec
| ALTER TABLE a ADD KEY (c(1))
< error
| 1089 HY000 Incorrect prefix key; the used key part isn't a string, the used length is longer than the key part, or the storage engine doesn't support unique prefix keys
> exec
| ALTER TABLE a ADD KEY (c) BLOCK_SIZE 1
< error
| 5258 HY000 Invalid block size, block size should between 1024 and 1048576
> exec
| ALTER TABLE a ADD SPATIAL (c)
< error
| 1687 42000 A SPATIAL index may only contain a geometrical type column.
> exec
| CREATE SPATIAL INDEX idx_c ON a (c)
< error
| 1687 42000 A SPATIAL index may only contain a geometrical type column.
> exec
| CREATE VECTOR INDEX idx_c ON a (c)
< error
| 7601 HY000 Column 'c' cannot be part of VECTOR index
> exec
| CREATE FULLTEXT INDEX idx_c ON a (c)
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| CREATE INDEX i ON a (c, c)
< error
| 1060 42S21 Duplicate column name 'c'
> exec
| ALTER TABLE a ADD KEY ((c + 1)), ADD KEY (v(1)) BLOCK_SIZE 16384
< ok 0 0

== TestOracle_OceanBaseKeyParts/CREATE_TABLE default#1
> exec
| CREATE TABLE t2 (b INT PRIMARY KEY, KEY (b, b))
< error
| 1060 42S21 Duplicate column name 'b'
> exec
| CREATE TABLE t2 (b INT PRIMARY KEY, KEY ((b)))
< error
| 3762 HY000 Functional index on a column is not supported. Consider using a regular index instead.
> exec
| CREATE TABLE t2 (b INT PRIMARY KEY, KEY (b(1)))
< error
| 1089 HY000 Incorrect prefix key; the used key part isn't a string, the used length is longer than the key part, or the storage engine doesn't support unique prefix keys
> exec
| CREATE TABLE t2 (b INT PRIMARY KEY, KEY (b) BLOCK_SIZE 1)
< error
| 5258 HY000 Invalid block size, block size should between 1024 and 1048576
