-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestContainerTiDBRefusedForms default#1
> exec
| CREATE TABLE t_vinv (a INT, b VARCHAR(255) INVISIBLE)
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 52 near "INVISIBLE)" 
> exec
| CREATE TABLE t_ivis (a INT VISIBLE)
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 34 near "VISIBLE)" 
> exec
| CREATE TABLE t_gen6 (col1 INT, col2 INT, col3 INT GENERATED ALWAYS AS (col1 + col2) VIRTUAL INVISIBLE)
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 101 near "INVISIBLE)" 
> exec
| CREATE TABLE t_vis (id INT PRIMARY KEY, val INT /*!80023 INVISIBLE */)
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 66 near "INVISIBLE */)" 
> exec
| ALTER TABLE t ALTER COLUMN b SET VISIBLE
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 40 near "VISIBLE" 
> exec
| ALTER TABLE t ALTER COLUMN b SET INVISIBLE
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 42 near "INVISIBLE" 
> exec
| ALTER VIEW v AS SELECT 2 AS b, 3 AS c
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 10 near "VIEW v AS SELECT 2 AS b, 3 AS c" 
> exec
| ALTER ALGORITHM = MERGE VIEW v AS SELECT 2 AS b
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 15 near "ALGORITHM = MERGE VIEW v AS SELECT 2 AS b" 
> exec
| CREATE VIEW v1 AS SELECT a FROM t WHERE a > 0 WITH CHECK OPTION
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 56 near "CHECK OPTION" 
