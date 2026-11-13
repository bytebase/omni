-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestDatabaseOptionsOracle default#1
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| ALTER DATABASE d CHARACTER SET = 'utf8'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER DATABASE d CHARACTER SET = 'utf8'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER DATABASE d CHARACTER SET 'utf8'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER DATABASE d CHARACTER SET 'utf8'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER DATABASE d CHARACTER SET = "utf8"
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER DATABASE d CHARACTER SET = "utf8"
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER DATABASE d CHARACTER SET = `utf8`
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER DATABASE d CHARACTER SET = `utf8`
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER DATABASE d CHARACTER SET utf8
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER DATABASE d CHARACTER SET utf8
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER DATABASE d DEFAULT CHARACTER SET = 'utf8'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER DATABASE d DEFAULT CHARACTER SET = 'utf8'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER SCHEMA d DEFAULT CHARSET = 'UTF8'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER SCHEMA d DEFAULT CHARSET = 'UTF8'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER DATABASE d CHARSET = binary
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER DATABASE d CHARSET = binary
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER DATABASE d CHARSET binary
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER DATABASE d CHARSET binary
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER DATABASE d CHARACTER SET binary
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER DATABASE d CHARACTER SET binary
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER DATABASE d COLLATE = 'utf8_bin'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER DATABASE d COLLATE = 'utf8_bin'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER DATABASE d COLLATE 'utf8_bin'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER DATABASE d COLLATE 'utf8_bin'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER DATABASE d COLLATE = `utf8_bin`
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER DATABASE d COLLATE = `utf8_bin`
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER DATABASE d DEFAULT COLLATE = 'utf8_bin'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER DATABASE d DEFAULT COLLATE = 'utf8_bin'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER DATABASE COLLATE = 'utf8_bin'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER SCHEMA COLLATE = 'utf8_bin'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER SCHEMA COLLATE = 'utf8_bin'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER DATABASE DEFAULT CHARSET = 'utf8mb4' COLLATE = 'utf8_bin'
< error
| 1253 42000 COLLATION 'utf8mb3_bin' is not valid for CHARACTER SET 'utf8mb4'
> exec ^
: string
| ALTER SCHEMA `` DEFAULT COLLATE = 'UTF8_BiN'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER SCHEMA `` DEFAULT COLLATE = 'UTF8_BiN'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER DATABASE d COLLATE = binary
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER DATABASE d COLLATE binary
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER DATABASE d COLLATE binary
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER DATABASE d CHARSET=binary COLLATE = binary
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER DATABASE d CHARSET=binary COLLATE = binary
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER DATABASE d ENCRYPTION = 'N'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER DATABASE d ENCRYPTION = 'N'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER DATABASE d ENCRYPTION 'N'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER DATABASE d ENCRYPTION 'N'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER DATABASE d DEFAULT ENCRYPTION = 'N'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER DATABASE d DEFAULT ENCRYPTION = 'N'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near '' at line 1
> exec ^
: string
| ALTER DATABASE d ENCRYPTION = N
< error
| 1064 42000 …near 'N' at line 1
> exec ^
: string
| ALTER DATABASE ENCRYPTION = 'N'
< error
| 1064 42000 …near '= 'N'' at line 1
> exec ^
: string
| ALTER DATABASE d READ ONLY = 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER DATABASE d READ ONLY = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER DATABASE d READ ONLY = 0
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER DATABASE d READ ONLY = 0
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER DATABASE d READ ONLY = DEFAULT
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER DATABASE d READ ONLY = DEFAULT
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER DATABASE d READ ONLY 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER DATABASE d READ ONLY 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER DATABASE READ ONLY = 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER DATABASE READ ONLY = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER DATABASE d CHARSET = utf8 READ ONLY = 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER DATABASE d CHARSET = utf8 READ ONLY = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER DATABASE d READ ONLY = 1 CHARSET = utf8
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER DATABASE d READ ONLY = 1 CHARSET = utf8
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER DATABASE d READ ONLY
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| ALTER DATABASE d READ ONLY = 2
< error
| 1064 42000 …near '2' at line 1
> exec ^
: string
| ALTER DATABASE d READ ONLY = TRUE
< error
| 1064 42000 …near 'TRUE' at line 1
> exec ^
: string
| ALTER DATABASE d READ ONLY = 'x'
< error
| 1064 42000 …near ''x'' at line 1
> exec ^
: string
| ALTER DATABASE d READ ONLY = abc
< error
| 1064 42000 …near 'abc' at line 1
> exec ^
: string
| ALTER DATABASE d READ
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| ALTER DATABASE d READ WRITE
< error
| 1064 42000 …near 'WRITE' at line 1
> exec ^
: string
| ALTER DATABASE d DEFAULT READ ONLY = 1
< error
| 1064 42000 …near 'READ ONLY = 1' at line 1
> exec ^
: string
| ALTER DATABASE d READ ONLY, CHARSET = utf8
< error
| 1064 42000 …near ', CHARSET = utf8' at line 1
> exec ^
: string
| ALTER DATABASE d CHARSET = utf8 READ ONLY
< error
| 1064 42000 …near '' at line 1
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
| 1064 42000 …near ', COLLATE = 'utf8_bin'' at line 1
> exec ^
: string
| ALTER DATABASE d, CHARSET = 'utf8'
< error
| 1064 42000 …near ', CHARSET = 'utf8'' at line 1
> exec ^
: string
| CREATE DATABASE d CHARSET = 'utf8', COLLATE = 'utf8_bin'
< error
| 1064 42000 …near ', COLLATE = 'utf8_bin'' at line 1
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
