-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestCreateOrReplaceOracle default#1
> exec
| CREATE TABLE t (id INT PRIMARY KEY, n INT)
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE VIEW v AS SELECT 1'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE ALGORITHM = MERGE VIEW v AS SELECT 1'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE ALGORITHM = UNDEFINED DEFINER = CURRENT_USER SQL SECURITY INVOKER VIEW v AS SELECT 1'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE DEFINER = CURRENT_USER VIEW v AS SELECT 1'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE DEFINER = ''root''@''localhost'' SQL SECURITY DEFINER VIEW v AS SELECT 1'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE SQL SECURITY INVOKER VIEW v AS SELECT 1'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE VIEW v (a) AS SELECT 1 WITH CHECK OPTION'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE SPATIAL REFERENCE SYSTEM 4120 NAME ''Greek'' DEFINITION ''GEOGCS["Greek",DATUM["Greek",SPHEROID["Bessel 1841",6377397.155,299.1528128]],PRIMEM["Greenwich",0],UNIT["degree",0.017453292519943278],AXIS["Lat",NORTH],AXIS["Lon",EAST]]'''
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec
| PREPARE parse_probe FROM 'CREATE ALGORITHM = MERGE OR REPLACE VIEW v AS SELECT 1'
< error
| 1064 42000 …near 'OR REPLACE VIEW v AS SELECT 1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE DEFINER = CURRENT_USER OR REPLACE VIEW v AS SELECT 1'
< error
| 1064 42000 …near 'OR REPLACE VIEW v AS SELECT 1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE OR REPLACE VIEW v AS SELECT 1'
< error
| 1064 42000 …near 'OR REPLACE VIEW v AS SELECT 1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE OR VIEW v AS SELECT 1'
< error
| 1064 42000 …near 'VIEW v AS SELECT 1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE OR'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE TEMPORARY VIEW v AS SELECT 1'
< error
| 1064 42000 …near 'TEMPORARY VIEW v AS SELECT 1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE TEMPORARY TABLE r (id INT)'
< error
| 1064 42000 …near 'TEMPORARY TABLE r (id INT)' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE SPATIAL REFERENCE SYSTEM IF NOT EXISTS 4120 NAME ''Greek'''
< error
| 1064 42000 …near 'IF NOT EXISTS 4120 NAME 'Greek'' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE SPATIAL INDEX ix ON t (g)'
< error
| 1064 42000 …near 'INDEX ix ON t (g)' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE TABLE r (id INT)'
< error
| 1064 42000 …near 'TABLE r (id INT)' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE DATABASE d1'
< error
| 1064 42000 …near 'DATABASE d1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE SCHEMA d1'
< error
| 1064 42000 …near 'SCHEMA d1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE INDEX ix ON t (n)'
< error
| 1064 42000 …near 'INDEX ix ON t (n)' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE UNIQUE INDEX ix ON t (n)'
< error
| 1064 42000 …near 'UNIQUE INDEX ix ON t (n)' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE FULLTEXT INDEX ix ON t (n)'
< error
| 1064 42000 …near 'FULLTEXT INDEX ix ON t (n)' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE USER u1'
< error
| 1064 42000 …near 'USER u1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE ROLE r1'
< error
| 1064 42000 …near 'ROLE r1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE FUNCTION f1() RETURNS INT DETERMINISTIC RETURN 1'
< error
| 1064 42000 …near 'FUNCTION f1() RETURNS INT DETERMINISTIC RETURN 1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE FUNCTION f1 RETURNS INTEGER SONAME ''x.so'''
< error
| 1064 42000 …near 'FUNCTION f1 RETURNS INTEGER SONAME 'x.so'' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE AGGREGATE FUNCTION f1 RETURNS INTEGER SONAME ''x.so'''
< error
| 1064 42000 …near 'AGGREGATE FUNCTION f1 RETURNS INTEGER SONAME 'x.so'' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE PROCEDURE p1() SELECT 1'
< error
| 1064 42000 …near 'PROCEDURE p1() SELECT 1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE TRIGGER tr1 BEFORE INSERT ON t FOR EACH ROW SET @a = 1'
< error
| 1064 42000 …near 'TRIGGER tr1 BEFORE INSERT ON t FOR EACH ROW SET @a = 1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE EVENT e1 ON SCHEDULE EVERY 1 DAY DO SELECT 1'
< error
| 1064 42000 …near 'EVENT e1 ON SCHEDULE EVERY 1 DAY DO SELECT 1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE TABLESPACE ts1 ADD DATAFILE ''x.ibd'''
< error
| 1064 42000 …near 'TABLESPACE ts1 ADD DATAFILE 'x.ibd'' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE UNDO TABLESPACE ts1 ADD DATAFILE ''x.ibu'''
< error
| 1064 42000 …near 'UNDO TABLESPACE ts1 ADD DATAFILE 'x.ibu'' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE SERVER s1 FOREIGN DATA WRAPPER mysql OPTIONS (USER ''x'')'
< error
| 1064 42000 …near 'SERVER s1 FOREIGN DATA WRAPPER mysql OPTIONS (USER 'x')' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE LOGFILE GROUP lg ADD UNDOFILE ''u.dat'' ENGINE = NDB'
< error
| 1064 42000 …near 'LOGFILE GROUP lg ADD UNDOFILE 'u.dat' ENGINE = NDB' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE RESOURCE GROUP rg TYPE = USER'
< error
| 1064 42000 …near 'RESOURCE GROUP rg TYPE = USER' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE DEFINER = CURRENT_USER FUNCTION f1() RETURNS INT DETERMINISTIC RETURN 1'
< error
| 1064 42000 …near 'FUNCTION f1() RETURNS INT DETERMINISTIC RETURN 1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE DEFINER = CURRENT_USER PROCEDURE p1() SELECT 1'
< error
| 1064 42000 …near 'PROCEDURE p1() SELECT 1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE DEFINER = CURRENT_USER TRIGGER tr1 BEFORE INSERT ON t FOR EACH ROW SET @a = 1'
< error
| 1064 42000 …near 'TRIGGER tr1 BEFORE INSERT ON t FOR EACH ROW SET @a = 1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE DEFINER = CURRENT_USER EVENT e1 ON SCHEDULE EVERY 1 DAY DO SELECT 1'
< error
| 1064 42000 …near 'EVENT e1 ON SCHEDULE EVERY 1 DAY DO SELECT 1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE ALGORITHM = MERGE TABLE r (id INT)'
< error
| 1064 42000 …near 'TABLE r (id INT)' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE r (id INT)'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TEMPORARY TABLE r (id INT)'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE INDEX ix ON t (n)'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE USER u1'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE FUNCTION f1() RETURNS INT DETERMINISTIC RETURN 1'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec
| PREPARE parse_probe FROM 'CREATE ALGORITHM = MERGE VIEW v AS SELECT 1'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE DEFINER = CURRENT_USER SQL SECURITY INVOKER VIEW v AS SELECT 1'
< ok 0 0
