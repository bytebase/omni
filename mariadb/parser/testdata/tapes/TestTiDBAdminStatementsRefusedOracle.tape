-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestTiDBAdminStatementsRefusedOracle default#1
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| RECOMMEND INDEX RUN
< error
| 1064 42000 …near 'RECOMMEND INDEX RUN' at line 1
> exec ^
: string
| RECOMMEND INDEX RUN WITH a = 1
< error
| 1064 42000 …near 'RECOMMEND INDEX RUN WITH a = 1' at line 1
> exec ^
: string
| RECOMMEND INDEX RUN FOR 'x'
< error
| 1064 42000 …near 'RECOMMEND INDEX RUN FOR 'x'' at line 1
> exec ^
: string
| RECOMMEND INDEX SHOW
< error
| 1064 42000 …near 'RECOMMEND INDEX SHOW' at line 1
> exec ^
: string
| RECOMMEND INDEX APPLY 1
< error
| 1064 42000 …near 'RECOMMEND INDEX APPLY 1' at line 1
> exec ^
: string
| RECOMMEND INDEX IGNORE 1
< error
| 1064 42000 …near 'RECOMMEND INDEX IGNORE 1' at line 1
> exec ^
: string
| RECOMMEND INDEX SET a = 1
< error
| 1064 42000 …near 'RECOMMEND INDEX SET a = 1' at line 1
> exec ^
: string
| SET CONFIG tikv log.level = 'info'
< error
| 1064 42000 …near 'tikv log.level = 'info'' at line 1
> exec ^
: string
| SET CONFIG '127.0.0.1:3306' a = 1
< error
| 1064 42000 …near ''127.0.0.1:3306' a = 1' at line 1
> exec ^
: string
| SET CONFIG tikv a-b.c = ON
< error
| 1064 42000 …near 'tikv a-b.c = ON' at line 1
> exec ^
: string
| PAUSE BACKUP LOGS
< error
| 1064 42000 …near 'PAUSE BACKUP LOGS' at line 1
> exec ^
: string
| RESUME BACKUP LOGS
< error
| 1064 42000 …near 'RESUME BACKUP LOGS' at line 1
> exec ^
: string
| STOP BACKUP LOGS
< error
| 1064 42000 …near 'BACKUP LOGS' at line 1
> exec ^
: string
| PURGE BACKUP LOGS FROM 'noop://'
< error
| 1064 42000 …near 'BACKUP LOGS FROM 'noop://'' at line 1
> exec ^
: string
| SHOW BACKUP LOGS STATUS
< error
| 1064 42000 …near 'BACKUP LOGS STATUS' at line 1
> exec ^
: string
| SHOW BACKUP LOGS METADATA FROM 'noop://'
< error
| 1064 42000 …near 'BACKUP LOGS METADATA FROM 'noop://'' at line 1
> exec ^
: string
| SHOW BACKUP METADATA FROM 'noop://'
< error
| 1064 42000 …near 'BACKUP METADATA FROM 'noop://'' at line 1
> exec ^
: string
| SHOW BR JOB 1
< error
| 1064 42000 …near 'BR JOB 1' at line 1
> exec ^
: string
| SHOW BR JOB QUERY 1
< error
| 1064 42000 …near 'BR JOB QUERY 1' at line 1
> exec ^
: string
| CANCEL BR JOB 1
< error
| 1064 42000 …near 'CANCEL BR JOB 1' at line 1
> exec ^
: string
| CANCEL DISTRIBUTION JOB 1
< error
| 1064 42000 …near 'CANCEL DISTRIBUTION JOB 1' at line 1
> exec ^
: string
| FLUSH TIDB PLUGINS a
< error
| 1064 42000 …near 'TIDB PLUGINS a' at line 1
> exec ^
: string
| FLUSH TIDB PLUGINS a, b
< error
| 1064 42000 …near 'TIDB PLUGINS a, b' at line 1
> exec ^
: string
| FLUSH LOCAL TIDB PLUGINS a
< error
| 1064 42000 …near 'TIDB PLUGINS a' at line 1
> exec ^
: string
| ALTER RANGE global PLACEMENT POLICY p
< error
| 1064 42000 …near 'RANGE global PLACEMENT POLICY p' at line 1
> exec ^
: string
| ALTER RANGE global PLACEMENT POLICY DEFAULT
< error
| 1064 42000 …near 'RANGE global PLACEMENT POLICY DEFAULT' at line 1
> exec ^
: string
| ALTER RANGE global PLACEMENT POLICY SET DEFAULT
< error
| 1064 42000 …near 'RANGE global PLACEMENT POLICY SET DEFAULT' at line 1
> exec
| SET @omni_kill = ?
: string
| KILL TIDB 1
< ok 0 0
> exec
| PREPARE omni_kill_stmt FROM @omni_kill
< error
| 1054 42S22 Unknown column 'TIDB' in 'KILL'
> exec
| SET @omni_kill = ?
: string
| KILL TIDB 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_kill_stmt FROM @omni_kill
< error
| 1054 42S22 Unknown column 'TIDB' in 'KILL'
> exec
| SET @omni_kill = ?
: string
| KILL TIDB CONNECTION 1
< ok 0 0
> exec
| PREPARE omni_kill_stmt FROM @omni_kill
< error
| 1054 42S22 Unknown column 'TIDB' in 'KILL'
> exec
| SET @omni_kill = ?
: string
| KILL TIDB CONNECTION 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_kill_stmt FROM @omni_kill
< error
| 1054 42S22 Unknown column 'TIDB' in 'KILL'
> exec
| SET @omni_kill = ?
: string
| KILL TIDB QUERY 1
< ok 0 0
> exec
| PREPARE omni_kill_stmt FROM @omni_kill
< error
| 1054 42S22 Unknown column 'TIDB' in 'KILL'
> exec
| SET @omni_kill = ?
: string
| KILL TIDB QUERY 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_kill_stmt FROM @omni_kill
< error
| 1054 42S22 Unknown column 'TIDB' in 'KILL'
> exec
| SET @omni_kill = ?
: string
| KILL TIDB CONNECTION_ID()
< ok 0 0
> exec
| PREPARE omni_kill_stmt FROM @omni_kill
< error
| 1054 42S22 Unknown column 'TIDB' in 'KILL'
> exec
| SET @omni_kill = ?
: string
| KILL TIDB CONNECTION_ID()
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_kill_stmt FROM @omni_kill
< error
| 1054 42S22 Unknown column 'TIDB' in 'KILL'
