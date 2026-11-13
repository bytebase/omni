-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestCreateOrReplaceTiDBOracle default#1
> exec
| DROP DATABASE IF EXISTS omni_or_replace
< ok 0 0
> exec
| CREATE DATABASE omni_or_replace
< ok 0 0
> exec
| USE omni_or_replace
< ok 0 0
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
| PREPARE parse_probe FROM 'CREATE OR REPLACE DEFINER = ''root''@''%'' SQL SECURITY DEFINER VIEW v AS SELECT 1'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE SQL SECURITY INVOKER VIEW v AS SELECT 1'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE VIEW v (a) AS SELECT 1 WITH LOCAL CHECK OPTION'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE PLACEMENT POLICY pp FOLLOWERS = 2'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE PLACEMENT POLICY IF NOT EXISTS pp FOLLOWERS = 2'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE ALGORITHM = MERGE OR REPLACE VIEW v AS SELECT 1'
< error
| 1064 42000 …line 1 column 27 near "OR REPLACE VIEW v AS SELECT 1" 
> exec
| PREPARE parse_probe FROM 'CREATE DEFINER = CURRENT_USER OR REPLACE VIEW v AS SELECT 1'
< error
| 1064 42000 …line 1 column 32 near "OR REPLACE VIEW v AS SELECT 1" 
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE OR REPLACE VIEW v AS SELECT 1'
< error
| 1064 42000 …line 1 column 20 near "OR REPLACE VIEW v AS SELECT 1" 
> exec
| PREPARE parse_probe FROM 'CREATE OR VIEW v AS SELECT 1'
< error
| 1064 42000 …line 1 column 14 near "VIEW v AS SELECT 1" 
> exec
| PREPARE parse_probe FROM 'CREATE OR PLACEMENT POLICY pp FOLLOWERS = 2'
< error
| 1064 42000 …line 1 column 19 near "PLACEMENT POLICY pp FOLLOWERS = 2" 
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE'
< error
| 1064 42000 …line 1 column 17 near "" 
> exec
| PREPARE parse_probe FROM 'CREATE OR'
< error
| 1064 42000 …line 1 column 9 near "" 
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE TEMPORARY VIEW v AS SELECT 1'
< error
| 1064 42000 …line 1 column 27 near "TEMPORARY VIEW v AS SELECT 1" 
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE TEMPORARY TABLE r (id INT)'
< error
| 1064 42000 …line 1 column 27 near "TEMPORARY TABLE r (id INT)" 
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE GLOBAL TEMPORARY TABLE r (id INT) ON COMMIT DELETE ROWS'
< error
| 1064 42000 …line 1 column 24 near "GLOBAL TEMPORARY TABLE r (id INT) ON COMMIT DELETE ROWS" 
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE SPATIAL REFERENCE SYSTEM 4120 NAME ''Greek'' DEFINITION ''GEOGCS[]'''
< error
| 1064 42000 …line 1 column 25 near "SPATIAL REFERENCE SYSTEM 4120 NAME 'Greek' DEFINITION 'GEOGCS[]'" 
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE SPATIAL INDEX ix ON t (g)'
< error
| 1064 42000 …line 1 column 25 near "SPATIAL INDEX ix ON t (g)" 
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE TABLE r (id INT)'
< error
| 1064 42000 …line 1 column 23 near "TABLE r (id INT)" 
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE TABLE r LIKE t'
< error
| 1064 42000 …line 1 column 23 near "TABLE r LIKE t" 
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE DATABASE d1'
< error
| 1064 42000 …line 1 column 26 near "DATABASE d1" 
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE SCHEMA d1'
< error
| 1064 42000 …line 1 column 24 near "SCHEMA d1" 
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE INDEX ix ON t (n)'
< error
| 1064 42000 …line 1 column 23 near "INDEX ix ON t (n)" 
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE UNIQUE INDEX ix ON t (n)'
< error
| 1064 42000 …line 1 column 24 near "UNIQUE INDEX ix ON t (n)" 
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE FULLTEXT INDEX ix ON t (n)'
< error
| 1064 42000 …line 1 column 26 near "FULLTEXT INDEX ix ON t (n)" 
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE USER u1'
< error
| 1064 42000 …line 1 column 22 near "USER u1" 
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE ROLE r1'
< error
| 1064 42000 …line 1 column 22 near "ROLE r1" 
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE FUNCTION f1() RETURNS INT DETERMINISTIC RETURN 1'
< error
| 1064 42000 …line 1 column 26 near "FUNCTION f1() RETURNS INT DETERMINISTIC RETURN 1" 
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE AGGREGATE FUNCTION f1 RETURNS INTEGER SONAME ''x.so'''
< error
| 1064 42000 …line 1 column 27 near "AGGREGATE FUNCTION f1 RETURNS INTEGER SONAME 'x.so'" 
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE PROCEDURE p1() SELECT 1'
< error
| 1064 42000 …line 1 column 27 near "PROCEDURE p1() SELECT 1" 
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE TRIGGER tr1 BEFORE INSERT ON t FOR EACH ROW SET @a = 1'
< error
| 1064 42000 …line 1 column 25 near "TRIGGER tr1 BEFORE INSERT ON t FOR EACH ROW SET @a = 1" 
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE EVENT e1 ON SCHEDULE EVERY 1 DAY DO SELECT 1'
< error
| 1064 42000 …line 1 column 23 near "EVENT e1 ON SCHEDULE EVERY 1 DAY DO SELECT 1" 
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE TABLESPACE ts1 ADD DATAFILE ''x.ibd'''
< error
| 1064 42000 …line 1 column 28 near "TABLESPACE ts1 ADD DATAFILE 'x.ibd'" 
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE UNDO TABLESPACE ts1 ADD DATAFILE ''x.ibu'''
< error
| 1064 42000 …line 1 column 22 near "UNDO TABLESPACE ts1 ADD DATAFILE 'x.ibu'" 
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE SERVER s1 FOREIGN DATA WRAPPER mysql OPTIONS (USER ''x'')'
< error
| 1064 42000 …line 1 column 24 near "SERVER s1 FOREIGN DATA WRAPPER mysql OPTIONS (USER 'x')" 
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE LOGFILE GROUP lg ADD UNDOFILE ''u.dat'' ENGINE = NDB'
< error
| 1064 42000 …line 1 column 25 near "LOGFILE GROUP lg ADD UNDOFILE 'u.dat' ENGINE = NDB" 
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE RESOURCE GROUP rg RU_PER_SEC = 100'
< error
| 1064 42000 …line 1 column 26 near "RESOURCE GROUP rg RU_PER_SEC = 100" 
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE SEQUENCE s1'
< error
| 1064 42000 …line 1 column 26 near "SEQUENCE s1" 
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE GLOBAL BINDING FOR SELECT 1 USING SELECT 1'
< error
| 1064 42000 …line 1 column 24 near "GLOBAL BINDING FOR SELECT 1 USING SELECT 1" 
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE DEFINER = CURRENT_USER PROCEDURE p1() SELECT 1'
< error
| 1064 42000 …line 1 column 50 near "PROCEDURE p1() SELECT 1" 
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE DEFINER = CURRENT_USER FUNCTION f1() RETURNS INT DETERMINISTIC RETURN 1'
< error
| 1064 42000 …line 1 column 49 near "FUNCTION f1() RETURNS INT DETERMINISTIC RETURN 1" 
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE ALGORITHM = MERGE TABLE r (id INT)'
< error
| 1064 42000 …line 1 column 41 near "TABLE r (id INT)" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE r (id INT)'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TEMPORARY TABLE r (id INT)'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE DATABASE d1'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE INDEX ix ON t (n)'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE USER u1'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE ROLE r1'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE PLACEMENT POLICY pp FOLLOWERS = 2'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE SEQUENCE s1'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE ALGORITHM = MERGE VIEW v AS SELECT 1'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE DEFINER = CURRENT_USER SQL SECURITY INVOKER VIEW v AS SELECT 1'
< ok 0 0
> exec
| DROP DATABASE IF EXISTS omni_or_replace
< ok 0 0
