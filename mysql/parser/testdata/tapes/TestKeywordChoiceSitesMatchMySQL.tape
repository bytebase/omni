-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestKeywordChoiceSitesMatchMySQL default#1
> exec
| CREATE DATABASE d CHAR SET utf8mb4
< error
| 1044 42000 Access denied for user 'omni_test'@'%' to database 'd'
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE DATABASE d CHAR SET utf8mb4
< ok 0 0
> exec
| CREATE DATABASE d DEFAULT CHAR SET = utf8mb4
< error
| 1044 42000 Access denied for user 'omni_test'@'%' to database 'd'
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE DATABASE d DEFAULT CHAR SET = utf8mb4
< ok 0 0
> exec
| ALTER DATABASE d CHAR SET utf8mb4
< error
| 1044 42000 Access denied for user 'omni_test'@'%' to database 'd'
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| ALTER DATABASE d CHAR SET utf8mb4
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER DATABASE d CHAR SET utf8mb4
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| ALTER DATABASE d CHAR utf8mb4
< error
| 1064 42000 …near 'utf8mb4' at line 1
> exec
| CREATE TABLE t (a INT) CHAR SET utf8mb4
< ok 0 0
> exec
| CREATE TABLE t (a INT) DEFAULT CHAR SET = utf8mb4
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a INT) DEFAULT CHAR SET = utf8mb4
< ok 0 0
> exec
| CREATE TABLE t (a INT) CHAR
< error
| 1064 42000 …near '' at line 1
> exec
| ALTER TABLE t CHAR SET utf8mb4
< ok 0 0
> exec
| LOAD DATA INFILE 'f' INTO TABLE t CHAR SET utf8mb4
< error
| 1045 28000 Access denied for user 'omni_test'@'%' (using password: YES)
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| LOAD DATA INFILE 'f' INTO TABLE t CHAR SET utf8mb4
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t CHAR SET utf8mb4
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| LOAD DATA INFILE 'f' INTO TABLE t CHARACTER utf8mb4
< error
| 1064 42000 …near 'utf8mb4' at line 1
> exec
| SELECT 1 INTO OUTFILE 'f' CHAR SET utf8mb4
< error
| 1227 42000 Access denied; you need (at least one of) the FILE privilege(s) for this operation
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| SELECT 1 INTO OUTFILE 'f' CHAR SET utf8mb4
< error
| 1227 42000 Access denied; you need (at least one of) the FILE privilege(s) for this operation
> exec ^
: string
| SELECT 1 INTO OUTFILE 'f' CHAR SET utf8mb4
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT 1 INTO OUTFILE 'f' CHARACTER utf8mb4
< error
| 1064 42000 …near 'utf8mb4' at line 1
> exec
| SHOW CHAR SET
< ok 0 0
> exec
| DROP TABLES t
< ok 0 0
> exec
| DROP TEMPORARY TABLES IF EXISTS t, u
< ok 0 0
> exec
| ANALYZE TABLES t
< ok 0 0
> exec
| ANALYZE LOCAL TABLES t
< ok 0 0
> exec
| CHECK TABLES t
< ok 0 0
> exec
| OPTIMIZE TABLES t
< ok 0 0
> exec
| REPAIR TABLES t
< ok 0 0
> exec
| CREATE LOGFILE GROUP g ADD UNDOFILE 'u' NO_WAIT
< error
| 1227 42000 Access denied; you need (at least one of) the CREATE TABLESPACE privilege(s) for this operation
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE LOGFILE GROUP g ADD UNDOFILE 'u' NO_WAIT
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE LOGFILE GROUP g ADD UNDOFILE 'u' NO_WAIT
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| CREATE LOGFILE GROUP g ADD UNDOFILE 'u' WAIT NO_WAIT
< error
| 1227 42000 Access denied; you need (at least one of) the CREATE TABLESPACE privilege(s) for this operation
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE LOGFILE GROUP g ADD UNDOFILE 'u' WAIT NO_WAIT
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE LOGFILE GROUP g ADD UNDOFILE 'u' WAIT NO_WAIT
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| ALTER LOGFILE GROUP g ADD UNDOFILE 'u' NO_WAIT
< error
| 1227 42000 Access denied; you need (at least one of) the CREATE TABLESPACE privilege(s) for this operation
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| ALTER LOGFILE GROUP g ADD UNDOFILE 'u' NO_WAIT
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER LOGFILE GROUP g ADD UNDOFILE 'u' NO_WAIT
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
