-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestSyntaxProbeControls default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| CREATE VECTOR INDEX i ON nosuch (a)
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE VECTOR INDEX i ON nosuch (a)
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE VECTOR INDEX i ON nosuch ((a))
< error
| 1146 42S02 Table 'test.nosuch' doesn't exist
> exec ^
: string
| CREATE VECTOR INDEX i ON nosuch ((a))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT nosuch
< error
| 1054 42S22 Unknown column 'nosuch' in 'field list'
> exec ^
: string
| SELECT nosuch
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT 1 FROM nosuch -- a comment
< error
| 1146 42S02 Table 'test.nosuch' doesn't exist
> exec ^
: string
| SELECT 1 FROM nosuch -- a comment
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a VECTOR(16384))
< error
| 1105 HY000 vector cannot have more than 16383 dimensions
> exec ^
: string
| CREATE TABLE t (a VECTOR(16384))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT @@nosuch
< error
| 1193 HY000 Unknown system variable 'nosuch'
> exec ^
: string
| SELEKT 1
< error
| 1064 42000 …line 1 column 6 near "SELEKT 1" 
> exec ^
: string
| SELECT DISTINCT ALL * FROM nosuch
< error
| 1221 HY000 Incorrect usage of ALL and DISTINCT
> exec ^
: string
| SELECT DISTINCT ALL * FROM nosuch
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of ALL and DISTINCT
> exec ^
: string
| ALTER TABLE t ALGORITHM = nosuch
< error
| 1800 HY000 Unknown ALGORITHM 'ALGORITHM'
> exec ^
: string
| ALTER TABLE t ALGORITHM = nosuch
| ; SELECT FROM
< error
| 1800 HY000 Unknown ALGORITHM 'ALGORITHM'
