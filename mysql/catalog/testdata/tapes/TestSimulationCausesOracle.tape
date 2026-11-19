-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestSimulationCausesOracle default#1
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW a AS SELECT b FROM a
< error
| 1347 HY000 'test.a' is not VIEW
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| ALTER VIEW a AS SELECT b FROM a
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE VIEW a AS SELECT 1 AS b
< ok 0 0
> exec
| ALTER VIEW a AS SELECT 1 AS b
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP VIEW a
< error
| 1347 HY000 'test.a' is not VIEW
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP VIEW a
< error
| 1051 42S02 Unknown table 'test.a'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(b)
< error
| 1283 HY000 Column 'b' cannot be part of FULLTEXT index
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT idx_d(c)
< error
| 1061 42000 Duplicate key name 'idx_d'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(x), ADD x INT
< error
| 1283 HY000 Column 'x' cannot be part of FULLTEXT index
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(x), ADD x TEXT
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE FULLTEXT INDEX ft ON a(b)
< error
| 1283 HY000 Column 'b' cannot be part of FULLTEXT index
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE s (b INT, FULLTEXT KEY ft (b))
< error
| 1283 HY000 Column 'b' cannot be part of FULLTEXT index
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(a)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD UNIQUE(a)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY idx_d(missing)
< error
| 1061 42000 Duplicate key name 'idx_d'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PRIMARY KEY(missing)
< error
| 1068 42000 Multiple primary key defined
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) WITH PARSER a
< error
| 1128 HY000 Function 'a' is not defined
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b TEXT)
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(b) WITH PARSER ngram
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX i ON a(missing)
< error
| 1072 42000 Key column 'missing' doesn't exist in table
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX idx_d ON a(missing)
< error
| 1061 42000 Duplicate key name 'idx_d'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE FULLTEXT INDEX i ON a(missing) WITH PARSER a
< error
| 1128 HY000 Function 'a' is not defined
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE s (b INT, KEY(missing))
< error
| 1072 42000 Key column 'missing' doesn't exist in table
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(c2), ADD c2 INT
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(c), DROP COLUMN c
< error
| 1072 42000 Key column 'c' doesn't exist in table
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(c), DROP COLUMN c, ADD COLUMN c INT
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(y), CHANGE c y INT
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME COLUMN c TO y, ADD KEY(c)
< error
| 1072 42000 Key column 'c' doesn't exist in table
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(x), ADD COLUMN x INT, RENAME COLUMN c TO y
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD COLUMN x INT, CHANGE x y INT
< error
| 1054 42S22 Unknown column 'x' in 'a'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(x(11)), ADD x VARCHAR(10)
< error
| 1089 HY000 Incorrect prefix key; the used key part isn't a string, the used length is longer than the key part, or the storage engine doesn't support unique prefix keys
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(x), ADD x TEXT
< error
| 1170 42000 BLOB/TEXT column 'x' used in key specification without a key length
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE s (id INT PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE a ADD FOREIGN KEY (x2) REFERENCES s(id), ADD x2 VARCHAR(10)
< error
| 3780 HY000 Referencing column 'x2' and referenced column 'id' in foreign key constraint 'a_ibfk_1' are incompatible.
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE s (id INT PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE a ADD FOREIGN KEY (x2) REFERENCES s(id) ON DELETE SET NULL, ADD x2 INT NOT NULL
< error
| 1830 HY000 Column 'x2' cannot be NOT NULL: needed in a foreign key constraint 'a_ibfk_1' SET NULL
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE s (id INT PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE a ADD FOREIGN KEY (x2) REFERENCES s(id), ADD x2 INT
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL KEY(x), ADD x GEOMETRY
< error
| 1252 42000 All parts of a SPATIAL index must be NOT NULL
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PRIMARY KEY(x2), ADD x2 INT
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(x), ADD COLUMN x INT, RENAME COLUMN x TO y
< error
| 1054 42S22 Unknown column 'x' in 'a'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME COLUMN x TO y
< error
| 1054 42S22 Unknown column 'x' in 'a'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD COLUMN x INT, ADD INDEX idx_x (x), DROP COLUMN x
< error
| 1091 42000 Can't DROP 'x'; check that column/key exists
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE s (b INT, FULLTEXT(b) WITH PARSER a)
< error
| 1128 HY000 Function 'a' is not defined
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE VIEW a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE VIEW a AS SELECT ABS(a)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE VIEW a AS SELECT 1 UNION SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE VIEW a AS SELECT a HAVING 1
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE VIEW a AS SELECT 1 WHERE a
< error
| 1054 42S22 Unknown column 'a' in 'where clause'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE VIEW a AS SELECT (SELECT 1 AS x HAVING x)
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW a AS SELECT b FROM a
< error
| 1050 42S01 Table 'a' already exists
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SERIAL
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `a` bigint unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `a` (`a`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR BYTE
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `a` binary(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR COLLATE BINARY
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `a` binary(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SERIAL NULL
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `a` bigint unsigned AUTO_INCREMENT,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `a` (`a`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT)
< ok 0 0
> exec
| ALTER TABLE a MODIFY c SERIAL
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int NOT NULL,
|   `c` bigint unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `c` (`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT)
< ok 0 0
> exec
| ALTER TABLE a CHANGE c y SERIAL
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int NOT NULL,
|   `y` bigint unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `y` (`y`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REMOVE PARTITIONING
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION (PARTITION p0)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION PARTITIONS 2
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION p0
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TRUNCATE PARTITION ALL
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a COALESCE PARTITION 1
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REORGANIZE PARTITION p0 INTO (PARTITION p1)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REORGANIZE PARTITION p0 INTO (PARTITION p1 VALUES LESS THAN (b))
< error
| 1054 42S22 Unknown column 'b' in 'partition function'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REORGANIZE PARTITION p0 INTO (PARTITION p1 ENGINE nosuch)
< error
| 1286 42000 Unknown storage engine 'nosuch'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REORGANIZE PARTITION p0 INTO (PARTITION p1 VALUES LESS THAN (10))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT) PARTITION BY RANGE(b) SUBPARTITION BY HASH(b) (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> exec
| ALTER TABLE a REORGANIZE PARTITION p0 INTO (PARTITION p1 VALUES LESS THAN (10) ENGINE InnoDB (SUBPARTITION s0 ENGINE nosuch))
< error
| 1286 42000 Unknown storage engine 'nosuch'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION (PARTITION p0 VALUES LESS THAN (b))
< error
| 1054 42S22 Unknown column 'b' in 'partition function'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| ALTER TABLE a REORGANIZE PARTITION p0 INTO (PARTITION p1 ENGINE nosuch)
< error
| 1286 42000 Unknown storage engine 'nosuch'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REBUILD PARTITION ALL
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DISCARD PARTITION ALL TABLESPACE
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a IMPORT PARTITION ALL TABLESPACE
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE s (b INT)
< ok 0 0
> exec
| ALTER TABLE a EXCHANGE PARTITION p0 WITH TABLE s
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a EXCHANGE PARTITION p0 WITH TABLE a
< error
| 1066 42000 Not unique table/alias: 'a'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ANALYZE PARTITION ALL
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CHECK PARTITION ALL
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT) PARTITION BY HASH(b) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE a REMOVE PARTITIONING
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD e BIT(65)
< error
| 1439 42000 Display width out of range for column 'e' (max = 64)
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD e INT(256)
< error
| 1439 42000 Display width out of range for column 'e' (max = 255)
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD e INT(4294967296)
< error
| 1439 42000 Display width out of range for column 'e' (max = 4294967295)
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD e CHAR(4294967296)
< error
| 1439 42000 Display width out of range for column 'e' (max = 4294967295)
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD e DECIMAL(4294967296)
< error
| 1439 42000 Display width out of range for column 'e' (max = 4294967295)
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD b INT(256)
< error
| 1439 42000 Display width out of range for column 'b' (max = 255)
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| ALTER TABLE a ADD e INT(256)
< error
| 1439 42000 Display width out of range for column 'e' (max = 255)
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MODIFY b INT(256)
< error
| 1439 42000 Display width out of range for column 'b' (max = 255)
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE s (b INT(256))
< error
| 1439 42000 Display width out of range for column 'b' (max = 255)
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION (PARTITION p1 VALUES LESS THAN (b))
< error
| 1054 42S22 Unknown column 'b' in 'partition function'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION (PARTITION p1 ENGINE nosuch)
< error
| 1286 42000 Unknown storage engine 'nosuch'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION (PARTITION p1 VALUES LESS THAN (b) ENGINE nosuch)
< error
| 1054 42S22 Unknown column 'b' in 'partition function'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENGINE = nosuch
< error
| 1286 42000 Unknown storage engine 'nosuch'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| ALTER TABLE a ENGINE = nosuch
< error
| 1286 42000 Unknown storage engine 'nosuch'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE s (b INT) ENGINE = nosuch
< error
| 1286 42000 Unknown storage engine 'nosuch'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD e INT(256), ENGINE = nosuch
< error
| 1439 42000 Display width out of range for column 'e' (max = 255)
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENGINE = nosuch, ADD e INT(256)
< error
| 1286 42000 Unknown storage engine 'nosuch'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE s (b INT(256)) ENGINE = nosuch
< error
| 1439 42000 Display width out of range for column 'b' (max = 255)
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE s (b INT) ENGINE = PERFORMANCE_SCHEMA
< error
| 1683 HY000 Invalid performance_schema usage.
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENGINE = PERFORMANCE_SCHEMA
< error
| 1031 HY000 Table storage engine for 'a' doesn't have this option
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| ALTER TABLE a ENGINE = PERFORMANCE_SCHEMA
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| CREATE TABLE s (b INT) ENGINE = nosuch
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| ALTER TABLE a ENGINE = nosuch
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| SET SESSION sql_mode = 'TRADITIONAL'
< ok 0 0
> exec
| CREATE TABLE s (b INT) ENGINE = nosuch
< error
| 1286 42000 Unknown storage engine 'nosuch'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| SET SESSION sql_mode = 'ANSI'
< ok 0 0
> exec
| CREATE TABLE s (b INT) ENGINE = nosuch
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| SET SESSION sql_mode = 'NO_ENGINE_SUBSTITUTION'
< ok 0 0
> exec
| ALTER TABLE a ENGINE = nosuch
< error
| 1286 42000 Unknown storage engine 'nosuch'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE s (b INT) ENGINE = MERGE
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENGINE = MERGE
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE s (b INT) ENGINE = HEAP
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE s (b INT) PARTITION BY RANGE(b) (PARTITION p0 VALUES LESS THAN (b))
< error
| 1054 42S22 Unknown column 'b' in 'partition function'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE s (b INT) PARTITION BY RANGE(b) SUBPARTITION BY HASH(b) (PARTITION p0 VALUES LESS THAN (10) ENGINE InnoDB (SUBPARTITION s0 ENGINE nosuch))
< error
| 1286 42000 Unknown storage engine 'nosuch'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE s (b INT) ENGINE = 'DEFAULT'
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENGINE = 'DEFAULT'
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD e FLOAT(256,2)
< error
| 1439 42000 Display width out of range for column 'e' (max = 255)
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD e DOUBLE(256,2)
< error
| 1439 42000 Display width out of range for column 'e' (max = 255)
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD e FLOAT(255,2)
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD e INT(255)
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD e BIT(64)
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW a AS SELECT *
< error
| 1096 HY000 No tables used
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT *
< error
| 1096 HY000 No tables used
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT t.*
< error
| 1051 42S02 Unknown table 't'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT t.* FROM a t
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT count(*)
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a;
< error
| 1050 42S01 Table 'a' already exists
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE z;
< error
| 4028 HY000 A table must have at least one visible column.
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE z (KEY (b));
< error
| 1072 42000 Key column 'b' doesn't exist in table
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE z PARTITION BY KEY();
< error
| 4028 HY000 A table must have at least one visible column.
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT * FROM (SELECT *) d
< error
| 1096 HY000 No tables used
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS WITH c AS (SELECT *) SELECT * FROM c
< error
| 1096 HY000 No tables used
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS WITH c AS (SELECT *) SELECT 1
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT * FROM DUAL
< error
| 1096 HY000 No tables used
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT 1 FROM DUAL
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE DATABASE cteother
< ok 1 0
> exec
| CREATE TABLE cteother.c (x INT)
< ok 0 0
> exec
| CREATE VIEW v AS WITH c AS (SELECT *) SELECT * FROM cteother.c
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 1 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE z (LIKE a);
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| ALTER TABLE nope ADD e INT DEFAULT NOW()
< error
| 1067 42000 Invalid default value for 'e'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD b INT DEFAULT CURRENT_TIMESTAMP
< error
| 1067 42000 Invalid default value for 'b'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MODIFY b INT DEFAULT NOW()
< error
| 1067 42000 Invalid default value for 'b'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD e INT DEFAULT (NOW())
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD e DATETIME DEFAULT NOW()
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE z (a INT DEFAULT NOW());
< error
| 1067 42000 Invalid default value for 'a'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE z (a DATETIME DEFAULT NOW());
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE z (a INT INVISIBLE);
< error
| 4028 HY000 A table must have at least one visible column.
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = ON
< ok 0 0
> exec
| CREATE TABLE z (a INT INVISIBLE);
< error
| 4028 HY000 A table must have at least one visible column.
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE z (a INT INVISIBLE, b INT);
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| CREATE TABLE z (a INT INVISIBLE) AS SELECT 1 AS b;
< ok 1 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
