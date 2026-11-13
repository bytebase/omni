-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestMariaDBCreateModifierOracle default#1
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
| PREPARE parse_probe FROM 'CREATE TEMPORARY SEQUENCE IF NOT EXISTS s1'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TEMPORARY VIEW v1 AS SELECT 1'
< error
| 1064 42000 …near 'VIEW v1 AS SELECT 1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE TEMPORARY VIEW v1 AS SELECT 1'
< error
| 1064 42000 …near 'VIEW v1 AS SELECT 1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TEMPORARY ALGORITHM=MERGE VIEW v1 AS SELECT 1'
< error
| 1064 42000 …near 'ALGORITHM=MERGE VIEW v1 AS SELECT 1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TEMPORARY DEFINER=CURRENT_USER VIEW v1 AS SELECT 1'
< error
| 1064 42000 …near 'DEFINER=CURRENT_USER VIEW v1 AS SELECT 1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TEMPORARY SQL SECURITY INVOKER VIEW v1 AS SELECT 1'
< error
| 1064 42000 …near 'SQL SECURITY INVOKER VIEW v1 AS SELECT 1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TEMPORARY INDEX ix ON t (n)'
< error
| 1064 42000 …near 'INDEX ix ON t (n)' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TEMPORARY UNIQUE INDEX ux ON t (n)'
< error
| 1064 42000 …near 'UNIQUE INDEX ux ON t (n)' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TEMPORARY FULLTEXT INDEX fx ON t (n)'
< error
| 1064 42000 …near 'FULLTEXT INDEX fx ON t (n)' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TEMPORARY SPATIAL INDEX sx ON t (n)'
< error
| 1064 42000 …near 'SPATIAL INDEX sx ON t (n)' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TEMPORARY DATABASE d1'
< error
| 1064 42000 …near 'DATABASE d1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TEMPORARY SCHEMA d1'
< error
| 1064 42000 …near 'SCHEMA d1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TEMPORARY FUNCTION f1() RETURNS INT RETURN 1'
< error
| 1064 42000 …near 'FUNCTION f1() RETURNS INT RETURN 1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TEMPORARY AGGREGATE FUNCTION agg1 RETURNS INTEGER SONAME ''x.so'''
< error
| 1064 42000 …near 'AGGREGATE FUNCTION agg1 RETURNS INTEGER SONAME 'x.so'' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TEMPORARY PROCEDURE p1() SELECT 1'
< error
| 1064 42000 …near 'PROCEDURE p1() SELECT 1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TEMPORARY TRIGGER tr1 BEFORE INSERT ON t FOR EACH ROW SET @a = 1'
< error
| 1064 42000 …near 'TRIGGER tr1 BEFORE INSERT ON t FOR EACH ROW SET @a = 1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TEMPORARY EVENT e1 ON SCHEDULE EVERY 1 DAY DO SELECT 1'
< error
| 1064 42000 …near 'EVENT e1 ON SCHEDULE EVERY 1 DAY DO SELECT 1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TEMPORARY USER u1'
< error
| 1064 42000 …near 'USER u1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TEMPORARY ROLE ro1'
< error
| 1064 42000 …near 'ROLE ro1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TEMPORARY SERVER sv1 FOREIGN DATA WRAPPER mysql OPTIONS (USER ''root'')'
< error
| 1064 42000 …near 'SERVER sv1 FOREIGN DATA WRAPPER mysql OPTIONS (USER 'root')' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TEMPORARY TABLESPACE ts1 ADD DATAFILE ''x.ibd'''
< error
| 1064 42000 …near 'TABLESPACE ts1 ADD DATAFILE 'x.ibd'' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE FUNCTION IF NOT EXISTS f1() RETURNS INT RETURN 1'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE FUNCTION IF f1() RETURNS INT RETURN 1'
< error
| 1064 42000 …near 'f1() RETURNS INT RETURN 1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE FUNCTION IF NOT f1() RETURNS INT RETURN 1'
< error
| 1064 42000 …near 'f1() RETURNS INT RETURN 1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE FUNCTION IF EXISTS f1() RETURNS INT RETURN 1'
< error
| 1064 42000 …near 'EXISTS f1() RETURNS INT RETURN 1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE PROCEDURE IF NOT EXISTS p1() SELECT 1'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE PROCEDURE IF EXISTS p1() SELECT 1'
< error
| 1064 42000 …near 'EXISTS p1() SELECT 1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE AGGREGATE FUNCTION IF NOT EXISTS agg1 RETURNS INTEGER SONAME ''x.so'''
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE AGGREGATE FUNCTION IF agg1 RETURNS INTEGER SONAME ''x.so'''
< error
| 1064 42000 …near 'agg1 RETURNS INTEGER SONAME 'x.so'' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE FUNCTION IF EXISTS udf1 RETURNS STRING SONAME ''x.so'''
< error
| 1064 42000 …near 'EXISTS udf1 RETURNS STRING SONAME 'x.so'' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TRIGGER IF NOT tr1 BEFORE INSERT ON t FOR EACH ROW SET @a = 1'
< error
| 1064 42000 …near 'tr1 BEFORE INSERT ON t FOR EACH ROW SET @a = 1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TRIGGER IF EXISTS tr1 BEFORE INSERT ON t FOR EACH ROW SET @a = 1'
< error
| 1064 42000 …near 'EXISTS tr1 BEFORE INSERT ON t FOR EACH ROW SET @a = 1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE EVENT IF e1 ON SCHEDULE EVERY 1 DAY DO SELECT 1'
< error
| 1064 42000 …near 'e1 ON SCHEDULE EVERY 1 DAY DO SELECT 1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE EVENT IF EXISTS e1 ON SCHEDULE EVERY 1 DAY DO SELECT 1'
< error
| 1064 42000 …near 'EXISTS e1 ON SCHEDULE EVERY 1 DAY DO SELECT 1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE USER IF u1'
< error
| 1064 42000 …near 'u1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE USER IF NOT u1'
< error
| 1064 42000 …near 'u1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE USER IF EXISTS u1'
< error
| 1064 42000 …near 'EXISTS u1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE ROLE IF ro1'
< error
| 1064 42000 …near 'ro1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE ROLE IF EXISTS ro1'
< error
| 1064 42000 …near 'EXISTS ro1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE SERVER IF sv1 FOREIGN DATA WRAPPER mysql OPTIONS (USER ''root'')'
< error
| 1064 42000 …near 'sv1 FOREIGN DATA WRAPPER mysql OPTIONS (USER 'root')' at line 1
