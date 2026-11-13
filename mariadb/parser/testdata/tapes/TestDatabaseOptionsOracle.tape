-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestDatabaseOptionsOracle default#1
> exec
| PREPARE omni_parses_stmt FROM ?
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
< error
| 1253 42000 COLLATION 'utf8mb3_bin' is not valid for CHARACTER SET 'utf8mb4'
> exec ^
: string
| ALTER DATABASE d CHARACTER SET = 'utf8' CHARACTER SET = 'utf8mb4'
< error
| 1302 HY000 Conflicting declarations: 'CHARACTER SET utf8mb3' and 'CHARACTER SET utf8mb4'
> exec ^
: string
| ALTER DATABASE d CHARACTER SET = 'utf8' CHARACTER SET = 'utf8mb4'
| ; SELECT FROM
< error
| 1302 HY000 Conflicting declarations: 'CHARACTER SET utf8mb3' and 'CHARACTER SET utf8mb4'
> exec ^
: string
| ALTER DATABASE d DEFAULT CHARSET = 'utf8mb4' DEFAULT COLLATE = 'utf8mb4_general_ci' CHARACTER SET = 'utf8' COLLATE = 'utf8mb4_bin'
< error
| 1302 HY000 Conflicting declarations: 'CHARACTER SET utf8mb4' and 'CHARACTER SET utf8mb3'
> exec ^
: string
| ALTER DATABASE d DEFAULT CHARSET = 'utf8mb4' DEFAULT COLLATE = 'utf8mb4_general_ci' CHARACTER SET = 'utf8' COLLATE = 'utf8mb4_bin'
| ; SELECT FROM
< error
| 1302 HY000 Conflicting declarations: 'CHARACTER SET utf8mb4' and 'CHARACTER SET utf8mb3'
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
< error
| 1253 42000 COLLATION 'utf8mb3_bin' is not valid for CHARACTER SET 'utf8mb4'
> exec ^
: string
| ALTER SCHEMA `` DEFAULT COLLATE = 'UTF8_BiN'
< ok 0 0
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
< error
| 1064 42000 …near 'ENCRYPTION = 'N'' at line 1
> exec ^
: string
| ALTER DATABASE d ENCRYPTION 'N'
< error
| 1064 42000 …near 'ENCRYPTION 'N'' at line 1
> exec ^
: string
| ALTER DATABASE d DEFAULT ENCRYPTION = 'N'
< error
| 1064 42000 …near 'ENCRYPTION = 'N'' at line 1
> exec ^
: string
| ALTER DATABASE d ENCRYPTION
< error
| 1064 42000 …near 'ENCRYPTION' at line 1
> exec ^
: string
| ALTER DATABASE ENCRYPTION = 'N'
< error
| 1064 42000 …near '= 'N'' at line 1
> exec ^
: string
| CREATE DATABASE d ENCRYPTION = 'N'
< error
| 1064 42000 …near 'ENCRYPTION = 'N'' at line 1
> exec ^
: string
| CREATE DATABASE d DEFAULT ENCRYPTION = 'Y'
< error
| 1064 42000 …near 'ENCRYPTION = 'Y'' at line 1
> exec ^
: string
| ALTER DATABASE d READ ONLY = 1
< error
| 1064 42000 …near 'READ ONLY = 1' at line 1
> exec ^
: string
| ALTER DATABASE d READ ONLY = 0
< error
| 1064 42000 …near 'READ ONLY = 0' at line 1
> exec ^
: string
| ALTER DATABASE d READ ONLY = DEFAULT
< error
| 1064 42000 …near 'READ ONLY = DEFAULT' at line 1
> exec ^
: string
| ALTER DATABASE d READ ONLY 1
< error
| 1064 42000 …near 'READ ONLY 1' at line 1
> exec ^
: string
| ALTER DATABASE d READ ONLY
< error
| 1064 42000 …near 'READ ONLY' at line 1
> exec ^
: string
| ALTER DATABASE READ ONLY = 1
< error
| 1064 42000 …near 'READ ONLY = 1' at line 1
> exec ^
: string
| ALTER DATABASE d READ
< error
| 1064 42000 …near 'READ' at line 1
> exec ^
: string
| ALTER DATABASE d READ WRITE
< error
| 1064 42000 …near 'READ WRITE' at line 1
> exec ^
: string
| ALTER DATABASE d DEFAULT READ ONLY = 1
< error
| 1064 42000 …near 'READ ONLY = 1' at line 1
> exec ^
: string
| ALTER DATABASE d CHARSET = utf8 READ ONLY
< error
| 1064 42000 …near 'READ ONLY' at line 1
> exec ^
: string
| ALTER DATABASE d COLLATE
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| ALTER DATABASE d COLLATE =
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| ALTER DATABASE d CHARSET =
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| ALTER DATABASE d CHARACTER SET
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| ALTER DATABASE d CHARACTER SET =
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| ALTER DATABASE d CHARACTER SET = 'utf8' COLLATE
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| ALTER DATABASE d CHARSET = 1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| ALTER DATABASE d COLLATE = 1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| ALTER DATABASE d CHARSET = select
< error
| 1064 42000 …near 'select' at line 1
> exec ^
: string
| ALTER DATABASE d COLLATE = select
< error
| 1064 42000 …near 'select' at line 1
> exec ^
: string
| CREATE DATABASE d CHARSET = 1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| CREATE DATABASE d COLLATE = select
< error
| 1064 42000 …near 'select' at line 1
> exec ^
: string
| ALTER DATABASE d CHARSET = 'utf8', COLLATE = 'utf8_bin'
< error
| 1064 42000 …near ' COLLATE = 'utf8_bin'' at line 1
> exec ^
: string
| ALTER DATABASE d, CHARSET = 'utf8'
< error
| 1064 42000 …near ' CHARSET = 'utf8'' at line 1
> exec ^
: string
| CREATE DATABASE d CHARSET = 'utf8', COLLATE = 'utf8_bin'
< error
| 1064 42000 …near ' COLLATE = 'utf8_bin'' at line 1
> exec ^
: string
| ALTER DATABASE
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| ALTER DATABASE d
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| ALTER SCHEMA
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| ALTER SCHEMA d
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| ALTER DATABASE d DEFAULT
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| ALTER DATABASE d DEFAULT = 'utf8'
< error
| 1064 42000 …near '= 'utf8'' at line 1
> exec ^
: string
| ALTER DATABASE d CHARSET = utf8 DEFAULT
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| ALTER DATABASE d CHARSET = utf8 CHARACTER
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| CREATE DATABASE d CHARSET = utf8 CHARACTER
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| ALTER DATABASE binary CHARSET = 'utf8'
< error
| 1064 42000 …near 'binary CHARSET = 'utf8'' at line 1
> exec ^
: string
| CREATE DATABASE d READ ONLY = 1
< error
| 1064 42000 …near 'READ ONLY = 1' at line 1
> exec ^
: string
| CREATE DATABASE d READ ONLY 0
< error
| 1064 42000 …near 'READ ONLY 0' at line 1
> exec ^
: string
| CREATE DATABASE d READ ONLY = DEFAULT
< error
| 1064 42000 …near 'READ ONLY = DEFAULT' at line 1
> exec ^
: string
| ALTER DATABASE CHARSET = 'utf8'
< error
| 1064 42000 …near '= 'utf8'' at line 1
> exec ^
: string
| ALTER DATABASE CHARSET utf8
< error
| 1064 42000 …near 'utf8' at line 1
> exec ^
: string
| ALTER DATABASE d PLACEMENT POLICY SET DEFAULT
< error
| 1064 42000 …near 'PLACEMENT POLICY SET DEFAULT' at line 1
> exec ^
: string
| ALTER DATABASE d PLACEMENT POLICY = 'p'
< error
| 1064 42000 …near 'PLACEMENT POLICY = 'p'' at line 1
> exec ^
: string
| ALTER DATABASE PLACEMENT POLICY = 'p'
< error
| 1064 42000 …near 'POLICY = 'p'' at line 1
> exec ^
: string
| ALTER DATABASE d PLACEMENT
< error
| 1064 42000 …near 'PLACEMENT' at line 1
> exec ^
: string
| ALTER DATABASE d SET TIFLASH REPLICA 1
< error
| 1064 42000 …near 'SET TIFLASH REPLICA 1' at line 1
> exec ^
: string
| ALTER DATABASE SET TIFLASH REPLICA 1
< error
| 1064 42000 …near 'SET TIFLASH REPLICA 1' at line 1
