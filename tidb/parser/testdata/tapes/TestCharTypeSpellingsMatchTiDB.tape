-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestCharTypeSpellingsMatchTiDB default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| CREATE TABLE t (a CHARACTER)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHARACTER(10))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR VARYING(50))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHARACTER VARYING(50))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VARCHARACTER(50))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHARACTER VARYING(0))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHARACTER VARYING(10) CHARACTER SET utf8)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VARCHARACTER(50) CHARSET utf8)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NATIONAL CHARACTER)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NATIONAL CHARACTER VARYING(50))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NATIONAL CHAR VARYING(50))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHARACTER(50))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NCHAR VARCHAR(50))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NCHAR VARYING(50))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NCHAR VARCHARACTER(50))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a LONG VARCHARACTER)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a LONG CHAR VARYING)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a LONG CHARACTER VARYING)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a LONG CHARACTER VARYING CHARACTER SET utf8)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a LONG CHAR SET utf8mb4)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR CHAR SET utf8mb4)
< ok 0 0
> exec ^
: string
| ALTER TABLE t MODIFY a CHARACTER VARYING(10)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t MODIFY a CHARACTER VARYING(10)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t CHANGE COLUMN a a LONG CHAR VARYING BINARY CHARSET utf8 FIRST
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t CHANGE COLUMN a a LONG CHAR VARYING BINARY CHARSET utf8 FIRST
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a VARCHAR)
< error
| 1064 42000 …line 1 column 26 near ")" 
> exec ^
: string
| CREATE TABLE t (a VARCHARACTER)
< error
| 1064 42000 …line 1 column 31 near ")" 
> exec ^
: string
| CREATE TABLE t (a CHAR VARYING)
< error
| 1064 42000 …line 1 column 31 near ")" 
> exec ^
: string
| CREATE TABLE t (a CHARACTER VARYING)
< error
| 1064 42000 …line 1 column 36 near ")" 
> exec ^
: string
| CREATE TABLE t (a NVARCHAR)
< error
| 1064 42000 …line 1 column 27 near ")" 
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHAR)
< error
| 1064 42000 …line 1 column 35 near ")" 
> exec ^
: string
| CREATE TABLE t (a NCHAR VARYING)
< error
| 1064 42000 …line 1 column 32 near ")" 
> exec ^
: string
| CREATE TABLE t (a VARBINARY)
< error
| 1064 42000 …line 1 column 28 near ")" 
> exec ^
: string
| CREATE TABLE t (a VARCHAR())
< error
| 1064 42000 …line 1 column 27 near "))" 
> exec ^
: string
| CREATE TABLE t (a VARCHARACTER())
< error
| 1064 42000 …line 1 column 32 near "))" 
> exec ^
: string
| CREATE TABLE t (a CHARACTER VARYING())
< error
| 1064 42000 …line 1 column 37 near "))" 
> exec ^
: string
| CREATE TABLE t (a CHAR VARYING())
< error
| 1064 42000 …line 1 column 32 near "))" 
> exec ^
: string
| CREATE TABLE t (a NCHAR VARYING())
< error
| 1064 42000 …line 1 column 33 near "))" 
> exec ^
: string
| CREATE TABLE t (a NVARCHAR())
< error
| 1064 42000 …line 1 column 28 near "))" 
> exec ^
: string
| CREATE TABLE t (a VARBINARY())
< error
| 1064 42000 …line 1 column 29 near "))" 
> exec ^
: string
| CREATE TABLE t (a LONG CHAR VARYING(10))
< error
| 1064 42000 …line 1 column 36 near "(10))" 
> exec ^
: string
| CREATE TABLE t (a LONG VARCHARACTER(10))
< error
| 1064 42000 …line 1 column 36 near "(10))" 
> exec ^
: string
| CREATE TABLE t (a LONG CHARACTER)
< error
| 1064 42000 …line 1 column 33 near ")" 
> exec ^
: string
| CREATE TABLE t (a TEXT CHARACTER)
< error
| 1064 42000 …line 1 column 33 near ")" 
> exec ^
: string
| CREATE TABLE t (a CHARACTER SET utf8)
< error
| 1064 42000 …line 1 column 31 near "SET utf8)" 
> exec ^
: string
| CREATE TABLE t (a CHAR SET utf8)
< error
| 1064 42000 …line 1 column 26 near "SET utf8)" 
