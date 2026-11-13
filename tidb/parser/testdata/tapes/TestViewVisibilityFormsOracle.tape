-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestViewVisibilityFormsOracle default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| CREATE TABLE t (a INT, b INT INVISIBLE)
< error
| 1064 42000 …line 1 column 38 near "INVISIBLE)" 
> exec ^
: string
| CREATE TABLE t (a INT, b INT VISIBLE)
< error
| 1064 42000 …line 1 column 36 near "VISIBLE)" 
> exec ^
: string
| CREATE TABLE t (a INT, b INT NOT NULL INVISIBLE DEFAULT 1)
< error
| 1064 42000 …line 1 column 47 near "INVISIBLE DEFAULT 1)" 
> exec ^
: string
| CREATE TABLE t (a INT PRIMARY KEY INVISIBLE)
< error
| 1064 42000 …line 1 column 43 near "INVISIBLE)" 
> exec ^
: string
| CREATE TABLE t (a INT, b INT AS (a + 1) VIRTUAL INVISIBLE)
< error
| 1064 42000 …line 1 column 57 near "INVISIBLE)" 
> exec ^
: string
| CREATE TABLE t (a INT, b INT GENERATED ALWAYS AS (a + 1) VIRTUAL VISIBLE)
< error
| 1064 42000 …line 1 column 72 near "VISIBLE)" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN b INT INVISIBLE
< error
| 1064 42000 …line 1 column 40 near "INVISIBLE" 
> exec ^
: string
| ALTER TABLE t MODIFY b INT INVISIBLE
< error
| 1064 42000 …line 1 column 36 near "INVISIBLE" 
> exec ^
: string
| ALTER TABLE t CHANGE b c INT VISIBLE
< error
| 1064 42000 …line 1 column 36 near "VISIBLE" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN b SET VISIBLE
< error
| 1064 42000 …line 1 column 40 near "VISIBLE" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN b SET INVISIBLE
< error
| 1064 42000 …line 1 column 42 near "INVISIBLE" 
> exec ^
: string
| ALTER TABLE t ALTER b SET INVISIBLE
< error
| 1064 42000 …line 1 column 35 near "INVISIBLE" 
> exec ^
: string
| ALTER TABLE t ALTER t.b SET INVISIBLE
< error
| 1064 42000 …line 1 column 37 near "INVISIBLE" 
> exec ^
: string
| ALTER TABLE t ALTER b SET DEFAULT 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER b SET DEFAULT 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER b DROP DEFAULT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER b DROP DEFAULT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER b SET
< error
| 1064 42000 …line 1 column 25 near "" 
> exec ^
: string
| ALTER TABLE t ALTER b DROP
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| ALTER TABLE t ALTER b
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| ALTER TABLE t ALTER INDEX i VISIBLE
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER INDEX i VISIBLE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER INDEX i INVISIBLE
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER INDEX i INVISIBLE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER INDEX i
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| CREATE TABLE t (a INT, INDEX i (a) INVISIBLE)
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) INVISIBLE
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i ON t (a) INVISIBLE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE visible (invisible INT)
< ok 0 0
> exec ^
: string
| SELECT 1 FROM t AS visible
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM t AS visible
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE VIEW v AS SELECT 1 WITH CHECK OPTION
< error
| 1064 42000 …line 1 column 36 near "CHECK OPTION" 
> exec ^
: string
| CREATE VIEW v AS SELECT 1 WITH CASCADED CHECK OPTION
< ok 0 0
> exec ^
: string
| CREATE VIEW v AS SELECT 1 WITH LOCAL CHECK OPTION
< ok 0 0
> exec ^
: string
| CREATE VIEW v AS SELECT 1 WITH LOCAL CHECK
< error
| 1064 42000 …line 1 column 42 near "" 
> exec ^
: string
| CREATE VIEW v AS SELECT 1 WITH LOCAL
< error
| 1064 42000 …line 1 column 36 near "" 
> exec ^
: string
| CREATE VIEW v AS SELECT 1 WITH CASCADED
< error
| 1064 42000 …line 1 column 39 near "" 
> exec ^
: string
| CREATE OR REPLACE ALGORITHM = MERGE DEFINER = CURRENT_USER SQL SECURITY INVOKER VIEW v (a) AS SELECT 1 WITH LOCAL CHECK OPTION
< ok 0 0
> exec ^
: string
| ALTER VIEW v AS SELECT 1
< error
| 1064 42000 …line 1 column 10 near "VIEW v AS SELECT 1" 
> exec ^
: string
| ALTER VIEW v (a) AS SELECT 1 WITH CASCADED CHECK OPTION
< error
| 1064 42000 …line 1 column 10 near "VIEW v (a) AS SELECT 1 WITH CASCADED CHECK OPTION" 
> exec ^
: string
| ALTER ALGORITHM = MERGE VIEW v AS SELECT 1
< error
| 1064 42000 …line 1 column 15 near "ALGORITHM = MERGE VIEW v AS SELECT 1" 
> exec ^
: string
| ALTER DEFINER = CURRENT_USER VIEW v AS SELECT 1
< error
| 1064 42000 …line 1 column 13 near "DEFINER = CURRENT_USER VIEW v AS SELECT 1" 
> exec ^
: string
| ALTER SQL SECURITY INVOKER VIEW v AS SELECT 1
< error
| 1064 42000 …line 1 column 9 near "SQL SECURITY INVOKER VIEW v AS SELECT 1" 
> exec ^
: string
| ALTER DEFINER = CURRENT_USER EVENT e ENABLE
< error
| 1064 42000 …line 1 column 13 near "DEFINER = CURRENT_USER EVENT e ENABLE" 
> exec ^
: string
| ALTER ALGORITHM = INPLACE TABLE t ADD c INT
< error
| 1064 42000 …line 1 column 15 near "ALGORITHM = INPLACE TABLE t ADD c INT" 
