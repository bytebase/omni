-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestCharTypeSpellingsMatchMySQL default#1
> exec
| CREATE TABLE t (a CHARACTER)
< ok 0 0
> exec
| CREATE TABLE t (a CHARACTER(10))
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a CHARACTER(10))
< ok 0 0
> exec
| CREATE TABLE t (a CHAR VARYING(50))
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a CHAR VARYING(50))
< ok 0 0
> exec
| CREATE TABLE t (a CHARACTER VARYING(50))
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a CHARACTER VARYING(50))
< ok 0 0
> exec
| CREATE TABLE t (a VARCHARACTER(50))
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a VARCHARACTER(50))
< ok 0 0
> exec
| CREATE TABLE t (a CHARACTER VARYING(10) CHARACTER SET utf8)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a CHARACTER VARYING(10) CHARACTER SET utf8)
< ok 0 0
> exec
| CREATE TABLE t (a VARCHARACTER(50) CHARSET utf8)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a VARCHARACTER(50) CHARSET utf8)
< ok 0 0
> exec
| CREATE TABLE t (a NATIONAL CHARACTER)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a NATIONAL CHARACTER)
< ok 0 0
> exec
| CREATE TABLE t (a NATIONAL CHARACTER VARYING(50))
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a NATIONAL CHARACTER VARYING(50))
< ok 0 0
> exec
| CREATE TABLE t (a NATIONAL CHAR VARYING(50))
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a NATIONAL CHAR VARYING(50))
< ok 0 0
> exec
| CREATE TABLE t (a NATIONAL VARCHARACTER(50))
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a NATIONAL VARCHARACTER(50))
< ok 0 0
> exec
| CREATE TABLE t (a NCHAR VARCHAR(50))
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a NCHAR VARCHAR(50))
< ok 0 0
> exec
| CREATE TABLE t (a NCHAR VARYING(50))
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a NCHAR VARYING(50))
< ok 0 0
> exec
| CREATE TABLE t (a NCHAR VARCHARACTER(50))
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a NCHAR VARCHARACTER(50))
< ok 0 0
> exec
| CREATE TABLE t (a LONG VARCHARACTER)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a LONG VARCHARACTER)
< ok 0 0
> exec
| CREATE TABLE t (a LONG CHAR VARYING)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a LONG CHAR VARYING)
< ok 0 0
> exec
| CREATE TABLE t (a LONG CHARACTER VARYING)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a LONG CHARACTER VARYING)
< ok 0 0
> exec
| CREATE TABLE t (a LONG CHARACTER VARYING CHARACTER SET utf8)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a LONG CHARACTER VARYING CHARACTER SET utf8)
< ok 0 0
> exec
| CREATE TABLE t (a LONG CHAR SET utf8mb4)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a LONG CHAR SET utf8mb4)
< ok 0 0
> exec
| CREATE TABLE t (a CHAR CHAR SET utf8mb4)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a CHAR CHAR SET utf8mb4)
< ok 0 0
> exec
| ALTER TABLE t MODIFY a CHARACTER VARYING(10)
< ok 0 0
> exec
| CREATE TABLE t (a VARCHAR)
< error
| 1064 42000 …near ')' at line 1
> exec
| CREATE TABLE t (a VARCHARACTER)
< error
| 1064 42000 …near ')' at line 1
> exec
| CREATE TABLE t (a CHAR VARYING)
< error
| 1064 42000 …near ')' at line 1
> exec
| CREATE TABLE t (a CHARACTER VARYING)
< error
| 1064 42000 …near ')' at line 1
> exec
| CREATE TABLE t (a NVARCHAR)
< error
| 1064 42000 …near ')' at line 1
> exec
| CREATE TABLE t (a NATIONAL VARCHAR)
< error
| 1064 42000 …near ')' at line 1
> exec
| CREATE TABLE t (a NCHAR VARYING)
< error
| 1064 42000 …near ')' at line 1
> exec
| CREATE TABLE t (a VARBINARY)
< error
| 1064 42000 …near ')' at line 1
> exec
| CREATE TABLE t (a VARCHAR())
< error
| 1064 42000 …near '))' at line 1
> exec
| CREATE TABLE t (a VARCHARACTER())
< error
| 1064 42000 …near '))' at line 1
> exec
| CREATE TABLE t (a CHARACTER VARYING())
< error
| 1064 42000 …near '))' at line 1
> exec
| CREATE TABLE t (a CHAR VARYING())
< error
| 1064 42000 …near '))' at line 1
> exec
| CREATE TABLE t (a NCHAR VARYING())
< error
| 1064 42000 …near '))' at line 1
> exec
| CREATE TABLE t (a NVARCHAR())
< error
| 1064 42000 …near '))' at line 1
> exec
| CREATE TABLE t (a VARBINARY())
< error
| 1064 42000 …near '))' at line 1
> exec
| CREATE TABLE t (a LONG CHAR VARYING(10))
< error
| 1064 42000 …near '(10))' at line 1
> exec
| CREATE TABLE t (a LONG VARCHARACTER(10))
< error
| 1064 42000 …near '(10))' at line 1
> exec
| CREATE TABLE t (a LONG CHARACTER)
< error
| 1064 42000 …near ')' at line 1
> exec
| CREATE TABLE t (a TEXT CHARACTER)
< error
| 1064 42000 …near ')' at line 1
> exec
| CREATE TABLE t (a CHARACTER SET utf8)
< error
| 1064 42000 …near 'SET utf8)' at line 1
> exec
| CREATE TABLE t (a CHAR SET utf8)
< error
| 1064 42000 …near 'SET utf8)' at line 1
