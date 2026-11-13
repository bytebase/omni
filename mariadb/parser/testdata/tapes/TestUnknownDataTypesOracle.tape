-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestUnknownDataTypesOracle default#1
> exec
| CREATE TABLE p (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, n INT, pid INT, CONSTRAINT ck_n CHECK (n > 0), INDEX i (n), CONSTRAINT f FOREIGN KEY (pid) REFERENCES p (id))
< ok 0 0
> exec
| CREATE TABLE pt (id INT) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20), PARTITION p2 VALUES LESS THAN MAXVALUE)
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TABLE z1 (a foo)'
< error
| 4161 HY000 Unknown data type: 'foo'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE z2 (a foo(3))'
< error
| 4161 HY000 Unknown data type: 'foo'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE z3 (a foo(3,2))'
< error
| 4161 HY000 Unknown data type: 'foo'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE z5 (a `INT`)'
< error
| 4161 HY000 Unknown data type: 'INT'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE z6 (a `foo`)'
< error
| 4161 HY000 Unknown data type: 'foo'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE z7 (a `foo`(3))'
< error
| 4161 HY000 Unknown data type: 'foo'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE z8 (a foo UNSIGNED)'
< error
| 4161 HY000 Unknown data type: 'foo'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE z9 (a foo NOT NULL)'
< error
| 4161 HY000 Unknown data type: 'foo'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE z10 (a foo DEFAULT 1)'
< error
| 4161 HY000 Unknown data type: 'foo'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE z11 (a foo CHARACTER SET utf8mb4)'
< error
| 4161 HY000 Unknown data type: 'foo'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE z12 (a foo COLLATE utf8mb4_bin)'
< error
| 4161 HY000 Unknown data type: 'foo'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE z13 (a foo BINARY)'
< error
| 4161 HY000 Unknown data type: 'foo'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE z14 (a foo.bar)'
< error
| 4161 HY000 Unknown data type: 'bar'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE z15 (a status)'
< error
| 4161 HY000 Unknown data type: 'status'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE z16 (a select)'
< error
| 4161 HY000 Unknown data type: 'select'
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t ADD COLUMN z25 foo'
< error
| 4161 HY000 Unknown data type: 'foo'
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t MODIFY a foo'
< error
| 4161 HY000 Unknown data type: 'foo'
> exec
| PREPARE parse_probe FROM 'SELECT CAST(''x'' AS foo)'
< error
| 4161 HY000 Unknown data type: 'foo'
> exec
| PREPARE parse_probe FROM 'SELECT CAST(''x'' AS foo(3))'
< error
| 4161 HY000 Unknown data type: 'foo'
> exec
| PREPARE parse_probe FROM 'SELECT CAST(''x'' AS `foo`)'
< error
| 4161 HY000 Unknown data type: 'foo'
> exec
| PREPARE parse_probe FROM 'SELECT CAST(''x'' AS `INT`)'
< error
| 4161 HY000 Unknown data type: 'INT'
> exec
| PREPARE parse_probe FROM 'SELECT CONVERT(''x'', foo)'
< error
| 4161 HY000 Unknown data type: 'foo'
> exec
| PREPARE parse_probe FROM 'CREATE FUNCTION f1(a foo) RETURNS INT RETURN 1'
< error
| 4161 HY000 Unknown data type: 'foo'
> exec
| PREPARE parse_probe FROM 'CREATE FUNCTION f2() RETURNS foo RETURN 1'
< error
| 4161 HY000 Unknown data type: 'foo'
> exec
| PREPARE parse_probe FROM 'CREATE PROCEDURE p1() BEGIN DECLARE b foo; END'
< error
| 4161 HY000 Unknown data type: 'foo'
