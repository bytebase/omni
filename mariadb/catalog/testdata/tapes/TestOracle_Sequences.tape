-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestOracle_Sequences default#1
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_sequences`
< ok 0 0
> exec
| CREATE DATABASE `omni_testoracle_sequences`
< ok 1 0
> exec
| USE `omni_testoracle_sequences`
< ok 0 0
> exec
| CREATE SEQUENCE s
< ok 0 0
> query
| SHOW CREATE SEQUENCE `s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE SEQUENCE `s` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_sequences`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_sequences`
< ok 1 0
> exec
| USE `omni_testoracle_sequences`
< ok 0 0
> exec
| CREATE SEQUENCE s START WITH 100 INCREMENT BY 5 MINVALUE 10 MAXVALUE 1000 CACHE 20 CYCLE COMMENT 'c''x'
< ok 0 0
> query
| SHOW CREATE SEQUENCE `s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE SEQUENCE `s` start with 100 minvalue 10 maxvalue 1000 increment by 5 cache 20 cycle ENGINE=InnoDB COMMENT='c''x'
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_sequences`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_sequences`
< ok 1 0
> exec
| USE `omni_testoracle_sequences`
< ok 0 0
> exec
| CREATE SEQUENCE s INCREMENT BY -1
< ok 0 0
> query
| SHOW CREATE SEQUENCE `s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE SEQUENCE `s` start with -1 minvalue -9223372036854775807 maxvalue -1 increment by -1 cache 1000 nocycle ENGINE=InnoDB
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_sequences`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_sequences`
< ok 1 0
> exec
| USE `omni_testoracle_sequences`
< ok 0 0
> exec
| CREATE SEQUENCE s NOCACHE NOMINVALUE NOMAXVALUE ENGINE=MyISAM
< ok 0 0
> query
| SHOW CREATE SEQUENCE `s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE SEQUENCE `s` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 nocache nocycle ENGINE=MyISAM
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_sequences`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_sequences`
< ok 1 0
> exec
| USE `omni_testoracle_sequences`
< ok 0 0
> exec
| CREATE SEQUENCE s START WITH 3 COMMENT='x' ENGINE=Aria
< ok 0 0
> query
| SHOW CREATE SEQUENCE `s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE SEQUENCE `s` start with 3 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=Aria COMMENT='x'
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_sequences`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_sequences`
< ok 1 0
> exec
| USE `omni_testoracle_sequences`
< ok 0 0
> exec
| CREATE SEQUENCE s AS INT UNSIGNED
< ok 0 0
> exec
| CREATE SEQUENCE u AS TINYINT START WITH -5 MINVALUE -10 MAXVALUE 10
< ok 0 0
> query
| SHOW CREATE SEQUENCE `s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE SEQUENCE `s` as int unsigned start with 1 minvalue 1 maxvalue 4294967294 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `u`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| u
: bytes
| CREATE SEQUENCE `u` as tinyint start with -5 minvalue -10 maxvalue 10 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_sequences`
< ok 2 0
> exec
| CREATE DATABASE `omni_testoracle_sequences`
< ok 1 0
> exec
| USE `omni_testoracle_sequences`
< ok 0 0
> exec
| CREATE SEQUENCE s START WITH 5 MINVALUE 10
< error
| 4085 HY000 Sequence 'omni_testoracle_sequences.s' has out of range value for options
> exec
| CREATE SEQUENCE u MAXVALUE 1 MINVALUE 1
< error
| 4085 HY000 Sequence 'omni_testoracle_sequences.u' has out of range value for options
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_sequences`
< ok 0 0
> exec
| CREATE DATABASE `omni_testoracle_sequences`
< ok 1 0
> exec
| USE `omni_testoracle_sequences`
< ok 0 0
> exec
| CREATE SEQUENCE v AS TINYINT MAXVALUE 200
< ok 0 0
> exec
| CREATE SEQUENCE w AS TINYINT MINVALUE -200
< ok 0 0
> exec
| CREATE SEQUENCE x AS TINYINT UNSIGNED MAXVALUE 300
< ok 0 0
> exec
| CREATE SEQUENCE y MAXVALUE 9223372036854775807
< ok 0 0
> query
| SHOW CREATE SEQUENCE `v`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE SEQUENCE `v` as tinyint start with 1 minvalue 1 maxvalue 126 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `w`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| w
: bytes
| CREATE SEQUENCE `w` as tinyint start with -127 minvalue -127 maxvalue 126 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `x`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| x
: bytes
| CREATE SEQUENCE `x` as tinyint unsigned start with 1 minvalue 1 maxvalue 254 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `y`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| y
: bytes
| CREATE SEQUENCE `y` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_sequences`
< ok 4 0
> exec
| CREATE DATABASE `omni_testoracle_sequences`
< ok 1 0
> exec
| USE `omni_testoracle_sequences`
< ok 0 0
> exec
| CREATE SEQUENCE s START WITH 100 INCREMENT BY 5 MINVALUE 10 MAXVALUE 1000
< ok 0 0
> exec
| ALTER SEQUENCE s START WITH 50
< ok 0 0
> exec
| ALTER SEQUENCE s INCREMENT BY 2 RESTART WITH 3
< error
| 4085 HY000 Sequence 'omni_testoracle_sequences.s' has out of range value for options
> exec
| ALTER SEQUENCE s MAXVALUE 5
< error
| 4085 HY000 Sequence 'omni_testoracle_sequences.s' has out of range value for options
> exec
| ALTER SEQUENCE s NOCACHE CYCLE
< ok 0 0
> query
| SHOW CREATE SEQUENCE `s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE SEQUENCE `s` start with 50 minvalue 10 maxvalue 1000 increment by 5 nocache cycle ENGINE=InnoDB
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_sequences`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_sequences`
< ok 1 0
> exec
| USE `omni_testoracle_sequences`
< ok 0 0
> exec
| CREATE SEQUENCE s
< ok 0 0
> exec ^
< error
| 1050 42S01 Table 's' already exists
> exec
| CREATE TABLE s (a INT)
< error
| 1050 42S01 Table 's' already exists
> exec
| CREATE VIEW s AS SELECT 1
< error
| 1050 42S01 Table 's' already exists
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE SEQUENCE t
< error
| 1050 42S01 Table 't' already exists
> exec
| DROP SEQUENCE t
< error
| 4090 42S02 'omni_testoracle_sequences.t' is not a SEQUENCE
> exec
| ALTER SEQUENCE t CACHE 5
< error
| 4089 42S02 'omni_testoracle_sequences.t' is not a SEQUENCE
> exec
| ALTER SEQUENCE missing CACHE 5
< error
| 4091 42S02 Unknown SEQUENCE: 'missing'
> exec
| DROP SEQUENCE IF EXISTS missing
< ok 0 0
> query
| SHOW CREATE SEQUENCE `s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE SEQUENCE `s` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_sequences`
< ok 2 0
> exec
| CREATE DATABASE `omni_testoracle_sequences`
< ok 1 0
> exec
| USE `omni_testoracle_sequences`
< ok 0 0
> exec
| CREATE SEQUENCE s
< ok 0 0
> exec
| CREATE OR REPLACE SEQUENCE s START WITH 7
< ok 0 0
> exec
| CREATE SEQUENCE IF NOT EXISTS s START WITH 9
< ok 0 0
> query
| SHOW CREATE SEQUENCE `s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE SEQUENCE `s` start with 7 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_sequences`
< ok 1 0
> exec
| CREATE DATABASE `omni_testoracle_sequences`
< ok 1 0
> exec
| USE `omni_testoracle_sequences`
< ok 0 0
> exec
| CREATE SEQUENCE s
< ok 0 0
> exec
| CREATE SEQUENCE u
< ok 0 0
> exec
| DROP TABLE s
< ok 0 0
> exec
| RENAME TABLE u TO w
< ok 0 0
> exec
| CREATE SEQUENCE s START WITH 2
< ok 0 0
> query
| SHOW CREATE SEQUENCE `s`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| s
: bytes
| CREATE SEQUENCE `s` start with 2 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
> query
| SHOW CREATE SEQUENCE `w`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| w
: bytes
| CREATE SEQUENCE `w` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB
