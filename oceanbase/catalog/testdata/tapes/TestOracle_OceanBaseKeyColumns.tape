-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_OceanBaseKeyColumns/FULLTEXT default#1
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, v VARCHAR(10), t TEXT, vb VARBINARY(10), j JSON)
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT (c, zz)
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| ALTER TABLE a ADD FULLTEXT (vb)
< error
| 1283 HY000 Column cannot be part of FULLTEXT index
> exec
| ALTER TABLE a ADD FULLTEXT (j)
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| CREATE FULLTEXT INDEX i ON a (c)
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| CREATE TABLE f (b INT PRIMARY KEY, c INT, FULLTEXT (c))
< error
| 1283 42000 A FTS index may only contain a text type column.
> exec
| ALTER TABLE a ADD FULLTEXT (v, t)
< ok 0 0

== TestOracle_OceanBaseKeyColumns/a_missing_column default#1
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, v VARCHAR(10), t TEXT, vb VARBINARY(10), j JSON)
< ok 0 0
> exec
| ALTER TABLE a ADD KEY (zz)
< error
| 1072 42000 Key column 'zz' doesn't exist in table
> exec
| ALTER TABLE a ADD KEY ((c + zz))
< error
| 1072 42000 Key column 'zz' doesn't exist in table
> exec
| ALTER TABLE a ADD FULLTEXT (zz, c)
< error
| 1072 42000 Key column 'zz' doesn't exist in table
> exec
| CREATE TABLE f (b INT PRIMARY KEY, FULLTEXT (zz))
< error
| 1072 42000 Key column 'zz' doesn't exist in table
> exec
| ALTER TABLE a ADD KEY k (zz), ADD zz INT
< ok 0 0
