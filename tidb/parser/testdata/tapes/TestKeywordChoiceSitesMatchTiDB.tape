-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestKeywordChoiceSitesMatchTiDB default#1
> exec
| SET @omni_probe = ?
: string
| CREATE DATABASE d CHAR SET utf8mb4
< ok 0 0
> exec
| PREPARE omni_probe_stmt FROM @omni_probe
< ok 0 0
> exec
| SET @omni_probe = ?
: string
| CREATE DATABASE d DEFAULT CHAR SET = utf8mb4
< ok 0 0
> exec
| PREPARE omni_probe_stmt FROM @omni_probe
< ok 0 0
> exec
| SET @omni_probe = ?
: string
| ALTER DATABASE d CHAR SET utf8mb4
< ok 0 0
> exec
| PREPARE omni_probe_stmt FROM @omni_probe
< ok 0 0
> exec
| SET @omni_probe = ?
: string
| ALTER DATABASE d CHAR utf8mb4
< ok 0 0
> exec
| PREPARE omni_probe_stmt FROM @omni_probe
< error
| 1064 42000 …line 1 column 29 near "utf8mb4" 
> exec
| SET @omni_probe = ?
: string
| CREATE TABLE t (a INT) CHAR SET utf8mb4
< ok 0 0
> exec
| PREPARE omni_probe_stmt FROM @omni_probe
< ok 0 0
> exec
| SET @omni_probe = ?
: string
| CREATE TABLE t (a INT) DEFAULT CHAR SET = utf8mb4
< ok 0 0
> exec
| PREPARE omni_probe_stmt FROM @omni_probe
< ok 0 0
> exec
| SET @omni_probe = ?
: string
| CREATE TABLE t (a INT) CHAR
< ok 0 0
> exec
| PREPARE omni_probe_stmt FROM @omni_probe
< error
| 1064 42000 …line 1 column 27 near "" 
> exec
| SET @omni_probe = ?
: string
| ALTER TABLE t CHAR SET utf8mb4
< ok 0 0
> exec
| PREPARE omni_probe_stmt FROM @omni_probe
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_probe = ?
: string
| ALTER TABLE t CHAR SET utf8mb4
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_probe_stmt FROM @omni_probe
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec
| SET @omni_probe = ?
: string
| LOAD DATA INFILE 'f' INTO TABLE t CHAR SET utf8mb4
< ok 0 0
> exec
| PREPARE omni_probe_stmt FROM @omni_probe
< error
| 1064 42000 …line 1 column 38 near "CHAR SET utf8mb4" 
> exec
| SET @omni_probe = ?
: string
| LOAD DATA INFILE 'f' INTO TABLE t CHARACTER utf8mb4
< ok 0 0
> exec
| PREPARE omni_probe_stmt FROM @omni_probe
< error
| 1064 42000 …line 1 column 51 near "utf8mb4" 
> exec
| SET @omni_probe = ?
: string
| SELECT 1 INTO OUTFILE 'f' CHAR SET utf8mb4
< ok 0 0
> exec
| PREPARE omni_probe_stmt FROM @omni_probe
< error
| 1064 42000 …line 1 column 30 near "CHAR SET utf8mb4" 
> exec
| SET @omni_probe = ?
: string
| SELECT 1 INTO OUTFILE 'f' CHARACTER utf8mb4
< ok 0 0
> exec
| PREPARE omni_probe_stmt FROM @omni_probe
< error
| 1064 42000 …line 1 column 35 near "CHARACTER utf8mb4" 
> exec
| SET @omni_probe = ?
: string
| SELECT 1 INTO OUTFILE 'f' CHARACTER SET utf8mb4
< ok 0 0
> exec
| PREPARE omni_probe_stmt FROM @omni_probe
< error
| 1064 42000 …line 1 column 35 near "CHARACTER SET utf8mb4" 
> exec
| SET @omni_probe = ?
: string
| SELECT 1 INTO OUTFILE 'f' CHARSET utf8mb4
< ok 0 0
> exec
| PREPARE omni_probe_stmt FROM @omni_probe
< error
| 1064 42000 …line 1 column 33 near "CHARSET utf8mb4" 
> exec
| SET @omni_probe = ?
: string
| SHOW CHAR SET
< ok 0 0
> exec
| PREPARE omni_probe_stmt FROM @omni_probe
< ok 0 0
> exec
| SET @omni_probe = ?
: string
| DROP TABLES t
< ok 0 0
> exec
| PREPARE omni_probe_stmt FROM @omni_probe
< ok 0 0
> exec
| SET @omni_probe = ?
: string
| DROP TEMPORARY TABLES IF EXISTS t, u
< ok 0 0
> exec
| PREPARE omni_probe_stmt FROM @omni_probe
< ok 0 0
> exec
| SET @omni_probe = ?
: string
| ANALYZE TABLES t
< ok 0 0
> exec
| PREPARE omni_probe_stmt FROM @omni_probe
< error
| 1064 42000 …line 1 column 14 near "TABLES t" 
> exec
| SET @omni_probe = ?
: string
| ANALYZE LOCAL TABLES t
< ok 0 0
> exec
| PREPARE omni_probe_stmt FROM @omni_probe
< error
| 1064 42000 …line 1 column 20 near "TABLES t" 
> exec
| SET @omni_probe = ?
: string
| CHECK TABLES t
< ok 0 0
> exec
| PREPARE omni_probe_stmt FROM @omni_probe
< error
| 1064 42000 …line 1 column 5 near "CHECK TABLES t" 
> exec
| SET @omni_probe = ?
: string
| OPTIMIZE TABLES t
< ok 0 0
> exec
| PREPARE omni_probe_stmt FROM @omni_probe
< error
| 8200 HY000 OPTIMIZE TABLE is not supported
> exec
| SET @omni_probe = ?
: string
| OPTIMIZE TABLES t
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_probe_stmt FROM @omni_probe
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec
| SET @omni_probe = ?
: string
| REPAIR TABLES t
< ok 0 0
> exec
| PREPARE omni_probe_stmt FROM @omni_probe
< error
| 1064 42000 …line 1 column 6 near "REPAIR TABLES t" 
> exec
| SET @omni_probe = ?
: string
| CREATE LOGFILE GROUP g ADD UNDOFILE 'u' NO_WAIT
< ok 0 0
> exec
| PREPARE omni_probe_stmt FROM @omni_probe
< error
| 1064 42000 …line 1 column 14 near "LOGFILE GROUP g ADD UNDOFILE 'u' NO_WAIT" 
> exec
| SET @omni_probe = ?
: string
| CREATE LOGFILE GROUP g ADD UNDOFILE 'u' WAIT NO_WAIT
< ok 0 0
> exec
| PREPARE omni_probe_stmt FROM @omni_probe
< error
| 1064 42000 …line 1 column 14 near "LOGFILE GROUP g ADD UNDOFILE 'u' WAIT NO_WAIT" 
> exec
| SET @omni_probe = ?
: string
| ALTER LOGFILE GROUP g ADD UNDOFILE 'u' NO_WAIT
< ok 0 0
> exec
| PREPARE omni_probe_stmt FROM @omni_probe
< error
| 1064 42000 …line 1 column 13 near "LOGFILE GROUP g ADD UNDOFILE 'u' NO_WAIT" 
