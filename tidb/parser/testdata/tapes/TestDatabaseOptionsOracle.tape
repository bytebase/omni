-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestDatabaseOptionsOracle default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| ALTER DATABASE d CHARACTER SET = 'utf8'
< ok 0 0
> exec ^
: string
| ALTER DATABASE d CHARACTER SET 'utf8'
< ok 0 0
> exec ^
: string
| ALTER DATABASE d CHARACTER SET = "utf8"
< ok 0 0
> exec ^
: string
| ALTER DATABASE d CHARACTER SET = `utf8`
< ok 0 0
> exec ^
: string
| ALTER DATABASE d CHARACTER SET utf8
< ok 0 0
> exec ^
: string
| ALTER DATABASE d DEFAULT CHARACTER SET = 'utf8'
< ok 0 0
> exec ^
: string
| ALTER SCHEMA d DEFAULT CHARSET = 'UTF8'
< ok 0 0
> exec ^
: string
| ALTER DATABASE d CHARSET = binary
< ok 0 0
> exec ^
: string
| ALTER DATABASE d CHARSET binary
< ok 0 0
> exec ^
: string
| ALTER DATABASE d CHARACTER SET binary
< ok 0 0
> exec ^
: string
| ALTER DATABASE d COLLATE = 'utf8_bin'
< ok 0 0
> exec ^
: string
| ALTER DATABASE d COLLATE 'utf8_bin'
< ok 0 0
> exec ^
: string
| ALTER DATABASE d COLLATE = `utf8_bin`
< ok 0 0
> exec ^
: string
| ALTER DATABASE d DEFAULT COLLATE = 'utf8_bin'
< ok 0 0
> exec ^
: string
| ALTER DATABASE d CHARSET = 'nosuchcharset'
< error
| 1115 42000 Unknown character set: 'nosuchcharset'
> exec ^
: string
| ALTER DATABASE d CHARSET = 'utf8mb4' COLLATE = 'utf8_bin'
< ok 0 0
> exec ^
: string
| ALTER DATABASE d CHARACTER SET = 'utf8' CHARACTER SET = 'utf8mb4'
< ok 0 0
> exec ^
: string
| ALTER DATABASE d DEFAULT CHARSET = 'utf8mb4' DEFAULT COLLATE = 'utf8mb4_general_ci' CHARACTER SET = 'utf8' COLLATE = 'utf8mb4_bin'
< ok 0 0
> exec ^
: string
| ALTER DATABASE COLLATE = 'utf8_bin'
< ok 0 0
> exec ^
: string
| ALTER SCHEMA COLLATE = 'utf8_bin'
< ok 0 0
> exec ^
: string
| ALTER DATABASE DEFAULT CHARSET = 'utf8mb4' COLLATE = 'utf8_bin'
< ok 0 0
> exec ^
: string
| ALTER SCHEMA `` DEFAULT COLLATE = 'UTF8_BiN'
< error
| 1102 42000 Incorrect database name ''
> exec ^
: string
| ALTER SCHEMA `` DEFAULT COLLATE = 'UTF8_BiN'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE DATABASE d
< ok 0 0
> exec ^
: string
| CREATE DATABASE d CHARACTER SET = 'utf8'
< ok 0 0
> exec ^
: string
| CREATE DATABASE d COLLATE = 'utf8_bin'
< ok 0 0
> exec ^
: string
| CREATE DATABASE d CHARSET = binary
< ok 0 0
> exec ^
: string
| CREATE SCHEMA d DEFAULT COLLATE = 'UTF8_BiN'
< ok 0 0
> exec ^
: string
| CREATE DATABASE d DEFAULT CHARSET = 'utf8mb4' COLLATE = 'utf8mb4_bin'
< ok 0 0
> exec ^
: string
| ALTER DATABASE d COLLATE = binary
< ok 0 0
> exec ^
: string
| ALTER DATABASE d COLLATE binary
< ok 0 0
> exec ^
: string
| ALTER DATABASE d CHARSET=binary COLLATE = binary
< ok 0 0
> exec ^
: string
| ALTER DATABASE d ENCRYPTION = 'N'
< ok 0 0
> exec ^
: string
| ALTER DATABASE d ENCRYPTION 'N'
< ok 0 0
> exec ^
: string
| ALTER DATABASE d DEFAULT ENCRYPTION = 'N'
< ok 0 0
> exec ^
: string
| CREATE DATABASE d ENCRYPTION = 'N'
< ok 0 0
> exec ^
: string
| CREATE DATABASE d DEFAULT ENCRYPTION = 'Y'
< ok 0 0
> exec ^
: string
| ALTER DATABASE d ENCRYPTION
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| ALTER DATABASE d ENCRYPTION = N
< error
| 1064 42000 …line 1 column 31 near "N" 
> exec ^
: string
| ALTER DATABASE ENCRYPTION = 'N'
< error
| 1064 42000 …line 1 column 27 near "= 'N'" 
> exec ^
: string
| ALTER DATABASE d READ ONLY = 1
< error
| 1064 42000 …line 1 column 21 near "READ ONLY = 1" 
> exec ^
: string
| ALTER DATABASE d READ ONLY = 0
< error
| 1064 42000 …line 1 column 21 near "READ ONLY = 0" 
> exec ^
: string
| ALTER DATABASE d READ ONLY = DEFAULT
< error
| 1064 42000 …line 1 column 21 near "READ ONLY = DEFAULT" 
> exec ^
: string
| ALTER DATABASE d READ ONLY 1
< error
| 1064 42000 …line 1 column 21 near "READ ONLY 1" 
> exec ^
: string
| ALTER DATABASE d READ ONLY
< error
| 1064 42000 …line 1 column 21 near "READ ONLY" 
> exec ^
: string
| ALTER DATABASE READ ONLY = 1
< error
| 1064 42000 …line 1 column 19 near "READ ONLY = 1" 
> exec ^
: string
| ALTER DATABASE d READ
< error
| 1064 42000 …line 1 column 21 near "READ" 
> exec ^
: string
| ALTER DATABASE d READ WRITE
< error
| 1064 42000 …line 1 column 21 near "READ WRITE" 
> exec ^
: string
| ALTER DATABASE d DEFAULT READ ONLY = 1
< error
| 1064 42000 …line 1 column 29 near "READ ONLY = 1" 
> exec ^
: string
| ALTER DATABASE d CHARSET = utf8 READ ONLY
< error
| 1064 42000 …line 1 column 36 near "READ ONLY" 
> exec ^
: string
| ALTER DATABASE d COLLATE
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| ALTER DATABASE d COLLATE =
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| ALTER DATABASE d CHARSET =
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| ALTER DATABASE d CHARACTER SET
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| ALTER DATABASE d CHARACTER SET =
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| ALTER DATABASE d CHARACTER SET = 'utf8' COLLATE
< error
| 1064 42000 …line 1 column 47 near "" 
> exec ^
: string
| ALTER DATABASE d CHARSET = 1
< error
| 1064 42000 …line 1 column 28 near "1" 
> exec ^
: string
| ALTER DATABASE d COLLATE = 1
< error
| 1064 42000 …line 1 column 28 near "1" 
> exec ^
: string
| ALTER DATABASE d CHARSET = select
< error
| 1064 42000 …line 1 column 33 near "select" 
> exec ^
: string
| ALTER DATABASE d COLLATE = select
< error
| 1064 42000 …line 1 column 33 near "select" 
> exec ^
: string
| CREATE DATABASE d CHARSET = 1
< error
| 1064 42000 …line 1 column 29 near "1" 
> exec ^
: string
| CREATE DATABASE d COLLATE = select
< error
| 1064 42000 …line 1 column 34 near "select" 
> exec ^
: string
| ALTER DATABASE d CHARSET = 'utf8', COLLATE = 'utf8_bin'
< error
| 1064 42000 …line 1 column 34 near ", COLLATE = 'utf8_bin'" 
> exec ^
: string
| ALTER DATABASE d, CHARSET = 'utf8'
< error
| 1064 42000 …line 1 column 17 near ", CHARSET = 'utf8'" 
> exec ^
: string
| CREATE DATABASE d CHARSET = 'utf8', COLLATE = 'utf8_bin'
< error
| 1064 42000 …line 1 column 35 near ", COLLATE = 'utf8_bin'" 
> exec ^
: string
| ALTER DATABASE COLLATE binary zzz
< error
| 1064 42000 …line 1 column 33 near "zzz" 
> exec ^
: string
| ALTER DATABASE
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| ALTER DATABASE d
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| ALTER SCHEMA
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| ALTER SCHEMA d
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| ALTER DATABASE d DEFAULT
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| ALTER DATABASE d DEFAULT = 'utf8'
< error
| 1064 42000 …line 1 column 26 near "= 'utf8'" 
> exec ^
: string
| ALTER DATABASE d CHARSET = utf8 DEFAULT
< error
| 1064 42000 …line 1 column 39 near "" 
> exec ^
: string
| ALTER DATABASE d CHARSET = utf8 CHARACTER
< error
| 1064 42000 …line 1 column 41 near "" 
> exec ^
: string
| CREATE DATABASE d CHARSET = utf8 CHARACTER
< error
| 1064 42000 …line 1 column 42 near "" 
> exec ^
: string
| ALTER DATABASE binary CHARSET = 'utf8'
< error
| 1064 42000 …line 1 column 21 near "binary CHARSET = 'utf8'" 
> exec ^
: string
| CREATE DATABASE d READ ONLY = 1
< error
| 1064 42000 …line 1 column 22 near "READ ONLY = 1" 
> exec ^
: string
| CREATE DATABASE d READ ONLY 0
< error
| 1064 42000 …line 1 column 22 near "READ ONLY 0" 
> exec ^
: string
| CREATE DATABASE d READ ONLY = DEFAULT
< error
| 1064 42000 …line 1 column 22 near "READ ONLY = DEFAULT" 
> exec ^
: string
| ALTER DATABASE CHARSET = 'utf8'
< error
| 1064 42000 …line 1 column 24 near "= 'utf8'" 
> exec ^
: string
| ALTER DATABASE CHARSET utf8
< error
| 1064 42000 …line 1 column 27 near "utf8" 
> exec ^
: string
| ALTER DATABASE d PLACEMENT POLICY SET DEFAULT
< ok 0 0
> exec ^
: string
| ALTER DATABASE d CHARSET = utf8 PLACEMENT POLICY SET DEFAULT
< ok 0 0
> exec ^
: string
| ALTER DATABASE d DEFAULT PLACEMENT POLICY = 'p'
< ok 0 0
> exec ^
: string
| ALTER DATABASE PLACEMENT POLICY = 'p'
< ok 0 0
> exec ^
: string
| ALTER DATABASE d SET TIFLASH REPLICA 1
< ok 0 0
> exec ^
: string
| ALTER DATABASE SET TIFLASH REPLICA 1
< ok 0 0
> exec ^
: string
| ALTER DATABASE d PLACEMENT POLICY SET
< error
| 1064 42000 …line 1 column 37 near "" 
> exec ^
: string
| ALTER DATABASE d PLACEMENT POLICY SET x
< error
| 1064 42000 …line 1 column 39 near "x" 
> exec ^
: string
| ALTER DATABASE d PLACEMENT
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| ALTER DATABASE d DEFAULT SET TIFLASH REPLICA 1
< error
| 1064 42000 …line 1 column 28 near "SET TIFLASH REPLICA 1" 
