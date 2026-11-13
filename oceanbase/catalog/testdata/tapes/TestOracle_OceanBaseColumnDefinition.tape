-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_OceanBaseColumnDefinition/before_the_name default#1
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT)
< ok 0 0
> exec
| ALTER TABLE a ADD c TINYINT DEFAULT 's'
< error
| 1067 42000 Invalid default value for 'c'
> exec
| ALTER TABLE a ADD c TINYINT DEFAULT CURRENT_TIMESTAMP
< error
| 1067 42000 Invalid default value for 'c'
> exec
| ALTER TABLE a ADD c TINYINT DEFAULT (b)
< error
| 1235 0A000 Column reference in default value expression not supported
> exec
| ALTER TABLE a ADD c TINYINT AS (zz)
< error
| 1054 42S22 Unknown column 'zz' in 'generated column function'
> exec
| ALTER TABLE a ADD c TINYINT AS (c) CHECK (c)
< error
| 1235 0A000 Add column-level check constraint while altering table not supported
> exec
| ALTER TABLE a ADD c TINYINT SRID 1
< error
| 1221 HY000 Incorrect usage of srid.
> exec
| ALTER TABLE a ADD c TINYINT PRIMARY KEY
< error
| 1068 42000 Multiple primary key defined
> exec
| ALTER TABLE a ADD c TINYINT CHECK (c)
< error
| 3812 HY000 An expression of non-boolean type specified to a check constraint 'a_OBCHECK_0000000000000001'
> exec
| ALTER TABLE a ADD c TINYINT CHECK (c > 0)
< error
| 1060 42S21 Duplicate column name 'c'
> exec
| ALTER TABLE a ADD e TINYINT AS (e)
< error
| 1054 42S22 Unknown column 'e' in 'generated column function'
> exec
| ALTER TABLE a ADD e TINYINT REFERENCES a (b) CHECK (e)
< ok 0 0

== TestOracle_OceanBaseColumnDefinition/literal_defaults default#1
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT)
< ok 0 0
> exec
| ALTER TABLE a ADD e1 TINYINT DEFAULT '1.5'
< ok 0 0
> exec
| ALTER TABLE a ADD e2 TINYINT DEFAULT 128
< error
| 1067 42000 Invalid default value for 'e2'
> exec
| ALTER TABLE a ADD e2 TINYINT UNSIGNED DEFAULT 255
< ok 0 0
> exec
| ALTER TABLE a ADD e3 INT DEFAULT ''
< error
| 1067 42000 Invalid default value for 'e3'
> exec
| ALTER TABLE a ADD e3 TINYINT DEFAULT 0x10
< ok 0 0
> exec
| ALTER TABLE a ADD e4 DECIMAL(3,1) DEFAULT 100
< error
| 1067 42000 Invalid default value for 'e4'
> exec
| ALTER TABLE a ADD e4 BIT(1) DEFAULT 2
< error
| 1067 42000 Invalid default value for 'e4'
> exec
| ALTER TABLE a ADD e5 CHAR(2) DEFAULT 'abc'
< error
| 1067 42000 Invalid default value for 'e5'
> exec
| ALTER TABLE a ADD e5 VARCHAR(2) DEFAULT 'éé'
< ok 0 0
> exec
| ALTER TABLE a ADD e6 ENUM('a') DEFAULT 'x'
< error
| 1067 42000 Invalid default value for 'e6'
> exec
| ALTER TABLE a ADD e6 SET('a','b') DEFAULT 'a,b'
< ok 0 0
> exec
| ALTER TABLE a ADD e7 DATETIME DEFAULT '2020-13-01'
< error
| 1067 42000 Invalid default value for 'e7'
> exec
| ALTER TABLE a ADD e7 DATE DEFAULT '20200101'
< ok 0 0
> exec
| ALTER TABLE a ADD e8 TIME DEFAULT 's'
< error
| 1067 42000 Invalid default value for 'e8'
> exec
| ALTER TABLE a ADD e8 YEAR DEFAULT 's'
< ok 0 0
> exec
| ALTER TABLE a ADD e9 DATE DEFAULT CURRENT_TIMESTAMP
< error
| 1067 42000 Invalid default value for 'e9'
> exec
| ALTER TABLE a ALTER c SET DEFAULT 's'
< error
| 1067 42000 Invalid default value for 'c'
> exec
| ALTER TABLE a ALTER zz SET DEFAULT (c)
< error
| 1235 0A000 Column reference in default value expression not supported
> exec
| ALTER TABLE a MODIFY c TINYINT DEFAULT 's'
< error
| 1067 42000 Invalid default value for 'c'
> exec
| CREATE TABLE t2 (b INT PRIMARY KEY, c TINYINT DEFAULT 's')
< error
| 1067 42000 Invalid default value for 'c'
