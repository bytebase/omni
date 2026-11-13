-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestMariaDBSequenceOracle default#1
> exec
| CREATE TABLE IF NOT EXISTS pt (id INT, v INT)
< ok 0 0
> exec
| CREATE SEQUENCE IF NOT EXISTS sq
< ok 0 0
> exec
| CREATE SEQUENCE cs_plain
< ok 0 0
> exec
| CREATE SEQUENCE IF NOT EXISTS cs_ine
< ok 0 0
> exec
| CREATE OR REPLACE SEQUENCE cs_cor
< ok 0 0
> exec
| CREATE TEMPORARY SEQUENCE cs_tmp
< ok 0 0
> exec
| CREATE OR REPLACE SEQUENCE cs1 START WITH 1000
< ok 0 0
> exec
| CREATE OR REPLACE SEQUENCE cs2 START = 5
< ok 0 0
> exec
| CREATE OR REPLACE SEQUENCE cs3 START 5
< ok 0 0
> exec
| CREATE OR REPLACE SEQUENCE cs4 INCREMENT BY 5
< ok 0 0
> exec
| CREATE OR REPLACE SEQUENCE cs5 INCREMENT = 5
< ok 0 0
> exec
| CREATE OR REPLACE SEQUENCE cs6 INCREMENT 5
< ok 0 0
> exec
| CREATE OR REPLACE SEQUENCE cs7 INCREMENT BY -1 MAXVALUE -1 MINVALUE -100
< ok 0 0
> exec
| CREATE OR REPLACE SEQUENCE cs8 INCREMENT BY +5 START WITH +5
< ok 0 0
> exec
| CREATE OR REPLACE SEQUENCE cs9 MINVALUE 10 MAXVALUE 9999
< ok 0 0
> exec
| CREATE OR REPLACE SEQUENCE cs10 MINVALUE = 1 MAXVALUE = 99
< ok 0 0
> exec
| CREATE OR REPLACE SEQUENCE cs11 NOMINVALUE NOMAXVALUE
< ok 0 0
> exec
| CREATE OR REPLACE SEQUENCE cs12 NO MINVALUE NO MAXVALUE
< ok 0 0
> exec
| CREATE OR REPLACE SEQUENCE cs13 CACHE 50
< ok 0 0
> exec
| CREATE OR REPLACE SEQUENCE cs14 CACHE = 50
< ok 0 0
> exec
| CREATE OR REPLACE SEQUENCE cs15 NOCACHE
< ok 0 0
> exec
| CREATE OR REPLACE SEQUENCE cs16 MINVALUE 1 MAXVALUE 5 CYCLE
< ok 0 0
> exec
| CREATE OR REPLACE SEQUENCE cs17 MINVALUE 1 MAXVALUE 5 NOCYCLE
< ok 0 0
> exec
| CREATE OR REPLACE SEQUENCE cs18 START WITH 100 MINVALUE 10 MAXVALUE 100000 INCREMENT BY 5 CACHE 20 CYCLE
< ok 0 0
> exec
| CREATE OR REPLACE SEQUENCE cs19 AS BIGINT
< ok 0 0
> exec
| CREATE OR REPLACE SEQUENCE cs20 AS INT UNSIGNED
< ok 0 0
> exec
| CREATE OR REPLACE SEQUENCE cs21 START WITH 1 AS BIGINT
< ok 0 0
> exec
| CREATE OR REPLACE SEQUENCE cs22 AS TINYINT START WITH 1
< ok 0 0
> exec
| CREATE SEQUENCE cr1 RESTART
< error
| 1064 42000 …near 'RESTART' at line 1
> exec
| CREATE SEQUENCE cr2 RESTART WITH 5
< error
| 1064 42000 …near 'RESTART' at line 1
> exec
| CREATE SEQUENCE cr3 CACHE
< error
| 1064 42000 …near '' at line 1
> exec
| CREATE SEQUENCE cr4 INCREMENT
< error
| 1064 42000 …near '' at line 1
> exec
| CREATE SEQUENCE cr5 START
< error
| 1064 42000 …near '' at line 1
> exec
| CREATE SEQUENCE cr6 INCREMENT BY
< error
| 1064 42000 …near '' at line 1
> exec
| CREATE SEQUENCE cr7 NO CACHE
< error
| 1064 42000 …near 'CACHE' at line 1
> exec
| CREATE SEQUENCE cr8 NO CYCLE
< error
| 1064 42000 …near 'CYCLE' at line 1
> exec
| CREATE SEQUENCE cr9 AS VARCHAR(10)
< error
| 1064 42000 …near 'VARCHAR(10)' at line 1
> exec
| CREATE SEQUENCE cr10 AS DECIMAL(10)
< error
| 1064 42000 …near 'DECIMAL(10)' at line 1
> exec
| CREATE SEQUENCE cr11 AS INT(11)
< error
| 1064 42000 …near '(11)' at line 1
> exec
| CREATE SEQUENCE cr12 AS DATE
< error
| 1064 42000 …near 'DATE' at line 1
> exec
| CREATE SEQUENCE cr13 START WITH 1+1
< error
| 1064 42000 …near '+1' at line 1
> exec
| CREATE SEQUENCE cr14 START WITH (5)
< error
| 1064 42000 …near '(5)' at line 1
> exec
| CREATE SEQUENCE cr15 CACHE CYCLE
< error
| 1064 42000 …near 'CYCLE' at line 1
> exec
| CREATE SEQUENCE
< error
| 1064 42000 …near '' at line 1
> exec
| CREATE OR SEQUENCE cr16
< error
| 1064 42000 …near 'SEQUENCE cr16' at line 1
> exec
| CREATE OR TEMPORARY SEQUENCE cr17
< error
| 1064 42000 …near 'TEMPORARY SEQUENCE cr17' at line 1
> exec
| ALTER SEQUENCE sq RESTART
< ok 0 0
> exec
| ALTER SEQUENCE sq RESTART WITH 500
< ok 0 0
> exec
| ALTER SEQUENCE sq RESTART = 5
< ok 0 0
> exec
| ALTER SEQUENCE sq RESTART 5
< ok 0 0
> exec
| ALTER SEQUENCE sq INCREMENT BY 10
< ok 0 0
> exec
| ALTER SEQUENCE sq MINVALUE 1 NOCACHE
< ok 0 0
> exec
| ALTER SEQUENCE IF EXISTS sq RESTART WITH 1
< ok 0 0
> exec
| ALTER SEQUENCE sq START WITH 5
< ok 0 0
> exec
| ALTER SEQUENCE sq AS BIGINT
< ok 0 0
> exec
| ALTER SEQUENCE sq RESTART INCREMENT BY 5
< ok 0 0
> exec
| ALTER SEQUENCE sq INCREMENT BY 5 RESTART
< ok 0 0
> exec
| ALTER SEQUENCE sq
< error
| 1064 42000 …near '' at line 1
> exec
| DROP SEQUENCE IF EXISTS d_a
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS d_a, d_b
< ok 0 0
> exec
| DROP TEMPORARY SEQUENCE IF EXISTS d_c
< ok 0 0
> exec
| DROP SEQUENCE d_a,
< error
| 1064 42000 …near '' at line 1
> exec
| SELECT NEXT VALUE FOR sq
< ok 0 0
> exec
| SELECT PREVIOUS VALUE FOR sq
< ok 0 0
> exec
| SELECT NEXT VALUE FOR sq + 1
< ok 0 0
> exec
| INSERT INTO pt (id, v) VALUES (NEXT VALUE FOR sq, 1)
< ok 1 0
> exec
| INSERT INTO pt (id) VALUES (PREVIOUS VALUE FOR sq)
< ok 1 0
> exec
| UPDATE pt SET id = NEXT VALUE FOR sq WHERE v = 1
< ok 1 0
> exec
| SELECT * FROM pt WHERE id = NEXT VALUE FOR sq
< ok 0 0
> exec
| SELECT id FROM pt HAVING id = NEXT VALUE FOR sq
< ok 0 0
> exec
| SELECT id FROM pt ORDER BY NEXT VALUE FOR sq
< ok 0 0
> exec
| CREATE OR REPLACE TABLE dseq1 (id INT DEFAULT NEXT VALUE FOR sq, v INT)
< ok 0 0
> exec
| CREATE OR REPLACE TABLE dseq2 (id INT DEFAULT (NEXT VALUE FOR sq), v INT)
< ok 0 0
> exec
| CREATE TABLE cseq1 (id INT CHECK (id < NEXT VALUE FOR sq))
< error
| 1970 42000 CHECK does not support subqueries or stored functions
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE cseq1 (id INT CHECK (id < NEXT VALUE FOR sq))
< error
| 1970 42000 CHECK does not support subqueries or stored functions
> exec ^
: string
| CREATE TABLE cseq1 (id INT CHECK (id < NEXT VALUE FOR sq))
| ; SELECT FROM
< error
| 1970 42000 CHECK does not support subqueries or stored functions
> exec
| CREATE TABLE gseq1 (id INT, g INT AS (NEXT VALUE FOR sq))
< error
| 1901 HY000 Function or expression 'nextval()' cannot be used in the GENERATED ALWAYS AS clause of `g`
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE gseq1 (id INT, g INT AS (NEXT VALUE FOR sq))
< error
| 1901 HY000 Function or expression 'nextval()' cannot be used in the GENERATED ALWAYS AS clause of `g`
> exec ^
: string
| CREATE TABLE gseq1 (id INT, g INT AS (NEXT VALUE FOR sq))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression 'nextval()' cannot be used in the GENERATED ALWAYS AS clause of `g`
> exec
| SELECT ABS(NEXT VALUE FOR sq)
< ok 0 0
> exec
| SELECT NEXT VALUE FOR sq OVER (ORDER BY 1)
< error
| 1064 42000 …near 'OVER (ORDER BY 1)' at line 1
> exec
| SELECT next value
< error
| 1064 42000 …near '' at line 1
> exec
| SELECT previous value
< error
| 1064 42000 …near '' at line 1
> exec
| SELECT next FROM pt
< error
| 1054 42S22 Unknown column 'next' in 'SELECT'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT next FROM pt
< error
| 1054 42S22 Unknown column 'next' in 'SELECT'
> exec ^
: string
| SELECT next FROM pt
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT previous FROM pt
< error
| 1054 42S22 Unknown column 'previous' in 'SELECT'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT previous FROM pt
< error
| 1054 42S22 Unknown column 'previous' in 'SELECT'
> exec ^
: string
| SELECT previous FROM pt
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT id AS next FROM pt
< ok 0 0
> exec
| SELECT id AS previous FROM pt
< ok 0 0
> exec
| SELECT next, value FROM pt
< error
| 1054 42S22 Unknown column 'next' in 'SELECT'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT next, value FROM pt
< error
| 1054 42S22 Unknown column 'next' in 'SELECT'
> exec ^
: string
| SELECT next, value FROM pt
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| CREATE OR REPLACE TABLE icols (sequence INT, previous INT, increment INT, cycle INT, minvalue INT, nocache INT, nocycle INT, nominvalue INT, nomaxvalue INT)
< ok 0 0
> exec
| select nExT vAlUe fOr sq
< ok 0 0
> exec
| CREATE OR REPLACE SEQUENCE `seq two` START WITH 1
< ok 0 0
> exec
| SELECT NEXT /* c */ VALUE FOR sq
< ok 0 0
