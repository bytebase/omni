-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestCommentValueOracle default#1
> exec
| CREATE TABLE ta (a INT)
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT COMMENT)'
< error
| 1064 42000 …near ')' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT COMMENT, b INT)'
< error
| 1064 42000 …near ', b INT)' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT COMMENT NOT NULL)'
< error
| 1064 42000 …near 'NOT NULL)' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT COMMENT DEFAULT 1)'
< error
| 1064 42000 …near 'DEFAULT 1)' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT COMMENT = ''x'')'
< error
| 1064 42000 …near '= 'x')' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT AS (1) COMMENT)'
< error
| 1064 42000 …near ')' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta ADD COLUMN b INT COMMENT'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta MODIFY a INT COMMENT'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta CHANGE a b INT COMMENT'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE INDEX i ON ta (a) COMMENT'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT, KEY i (a) COMMENT)'
< error
| 1064 42000 …near ')' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT, KEY i (a) COMMENT, b INT)'
< error
| 1064 42000 …near ', b INT)' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT, UNIQUE KEY i (a) COMMENT)'
< error
| 1064 42000 …near ')' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT, PRIMARY KEY (a) COMMENT)'
< error
| 1064 42000 …near ')' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT, CONSTRAINT c UNIQUE (a) COMMENT)'
< error
| 1064 42000 …near ')' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a TEXT, FULLTEXT KEY i (a) COMMENT)'
< error
| 1064 42000 …near ')' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT, KEY i (a) KEY_BLOCK_SIZE =)'
< error
| 1064 42000 …near ')' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT, KEY i (a) USING)'
< error
| 1064 42000 …near ')' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a TEXT, FULLTEXT KEY i (a) WITH PARSER)'
< error
| 1064 42000 …near ')' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT, KEY i (a) ENGINE_ATTRIBUTE)'
< error
| 1064 42000 …near ')' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT, KEY i (a) SECONDARY_ENGINE_ATTRIBUTE =)'
< error
| 1064 42000 …near ')' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta ADD INDEX i (a) COMMENT'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta ADD UNIQUE (a) COMMENT'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta ADD PRIMARY KEY (a) COMMENT'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta ADD INDEX i (a) COMMENT, ADD COLUMN b INT'
< error
| 1064 42000 …near ', ADD COLUMN b INT' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE PROCEDURE p() COMMENT SELECT 1'
< error
| 1064 42000 …near 'SELECT 1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE PROCEDURE p() COMMENT BEGIN END'
< error
| 1064 42000 …near 'BEGIN END' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE FUNCTION f() RETURNS INT COMMENT RETURN 1'
< error
| 1064 42000 …near 'RETURN 1' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE FUNCTION f() RETURNS INT NO SQL COMMENT RETURN 1'
< error
| 1064 42000 …near 'RETURN 1' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER PROCEDURE p COMMENT'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER FUNCTION f COMMENT'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE EVENT e ON SCHEDULE EVERY 1 DAY COMMENT DO SELECT 1'
< error
| 1064 42000 …near 'DO SELECT 1' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER EVENT e COMMENT'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER EVENT e COMMENT DO SELECT 1'
< error
| 1064 42000 …near 'DO SELECT 1' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER EVENT e ENABLE COMMENT'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE USER u COMMENT'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER USER u COMMENT'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE USER u ATTRIBUTE'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER USER u ATTRIBUTE'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLESPACE ts ADD DATAFILE ''ts.ibd'' COMMENT'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLESPACE ts ADD DATAFILE ''ts.ibd'' COMMENT ='
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE LOGFILE GROUP lg ADD UNDOFILE ''u.dat'' COMMENT'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE LOGFILE GROUP lg ADD UNDOFILE ''u.dat'' COMMENT ='
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT COMMENT ''x'', b INT)'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT NOT NULL COMMENT "x" DEFAULT 1)'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT AS (1) COMMENT ''x'')'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta MODIFY a INT COMMENT ''x'''
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta ADD COLUMN b INT COMMENT ''x'' FIRST'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE INDEX i ON ta (a) COMMENT ''x'''
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT, KEY i (a) COMMENT ''x'' USING BTREE KEY_BLOCK_SIZE = 8)'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT, PRIMARY KEY (a) COMMENT ''x'', UNIQUE (a) COMMENT ''y'')'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a TEXT, FULLTEXT KEY i (a) WITH PARSER ngram)'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta ADD INDEX i (a) COMMENT ''x'', ADD COLUMN b INT'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE PROCEDURE p() COMMENT ''x'' SELECT 1'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec
| PREPARE parse_probe FROM 'CREATE FUNCTION f() RETURNS INT COMMENT ''x'' RETURN 1'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec
| PREPARE parse_probe FROM 'ALTER PROCEDURE p COMMENT ''x'''
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec
| PREPARE parse_probe FROM 'CREATE EVENT e ON SCHEDULE EVERY 1 DAY COMMENT ''x'' DO SELECT 1'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec
| PREPARE parse_probe FROM 'ALTER EVENT e COMMENT ''x'' DO SELECT 1'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec
| PREPARE parse_probe FROM 'CREATE USER u COMMENT ''x'''
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER USER u COMMENT ''x'''
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE USER u ATTRIBUTE ''{"a": 1}'''
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TABLESPACE ts ADD DATAFILE ''ts.ibd'' COMMENT = ''x'''
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec
| PREPARE parse_probe FROM 'CREATE LOGFILE GROUP lg ADD UNDOFILE ''u.dat'' COMMENT ''x'''
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
