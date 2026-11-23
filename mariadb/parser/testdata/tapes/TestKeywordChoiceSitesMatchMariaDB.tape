-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestKeywordChoiceSitesMatchMariaDB default#1
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE DATABASE d CHAR SET utf8mb4
< ok 0 0
> exec ^
: string
| CREATE DATABASE d DEFAULT CHAR SET = utf8mb4
< ok 0 0
> exec ^
: string
| ALTER DATABASE d CHAR SET utf8mb4
< ok 0 0
> exec ^
: string
| ALTER DATABASE d CHAR utf8mb4
< error
| 1064 42000 …near 'utf8mb4' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) CHAR SET utf8mb4
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) DEFAULT CHAR SET = utf8mb4
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) CHAR
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| ALTER TABLE t CHAR SET utf8mb4
< ok 0 0
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t CHAR SET utf8mb4
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t CHAR SET utf8mb4
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t CHARACTER utf8mb4
< error
| 1064 42000 …near 'utf8mb4' at line 1
> exec ^
: string
| SELECT 1 INTO OUTFILE 'f' CHAR SET utf8mb4
< ok 0 0
> exec ^
: string
| SELECT 1 INTO OUTFILE 'f' CHARACTER utf8mb4
< error
| 1064 42000 …near 'utf8mb4' at line 1
> exec ^
: string
| SHOW CHAR SET
< ok 0 0
> exec ^
: string
| DROP TABLES t
< ok 0 0
> exec ^
: string
| DROP TEMPORARY TABLES IF EXISTS t, u
< ok 0 0
> exec ^
: string
| ANALYZE TABLES t
< ok 0 0
> exec ^
: string
| ANALYZE LOCAL TABLES t
< ok 0 0
> exec ^
: string
| CHECK TABLES t
< ok 0 0
> exec ^
: string
| OPTIMIZE TABLES t
< ok 0 0
> exec ^
: string
| REPAIR TABLES t
< ok 0 0
> exec ^
: string
| CREATE LOGFILE GROUP g ADD UNDOFILE 'u' NO_WAIT
< error
| 1064 42000 …near 'LOGFILE GROUP g ADD UNDOFILE 'u' NO_WAIT' at line 1
> exec ^
: string
| CREATE LOGFILE GROUP g ADD UNDOFILE 'u' WAIT NO_WAIT
< error
| 1064 42000 …near 'LOGFILE GROUP g ADD UNDOFILE 'u' WAIT NO_WAIT' at line 1
> exec ^
: string
| ALTER LOGFILE GROUP g ADD UNDOFILE 'u' NO_WAIT
< error
| 1064 42000 …near 'LOGFILE GROUP g ADD UNDOFILE 'u' NO_WAIT' at line 1
