-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestDatabaseOptionsOceanBaseOracle default#1
> exec
| DROP DATABASE IF EXISTS obdbopt
< ok 0 0
> exec
| DROP DATABASE IF EXISTS obdbopt_c
< ok 0 0
> exec
| CREATE DATABASE obdbopt
< ok 1 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt CHARACTER SET = 'utf8'
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt CHARACTER SET 'utf8'
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt CHARACTER SET = "utf8"
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt CHARACTER SET = `utf8`
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt CHARACTER SET utf8
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt DEFAULT CHARACTER SET = 'utf8'
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER SCHEMA obdbopt DEFAULT CHARSET = 'UTF8'
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt COLLATE = 'utf8_bin'
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt COLLATE 'utf8_bin'
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt DEFAULT COLLATE = 'utf8_bin'
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt CHARSET = 'nosuchcharset'
< error
| 1115 42000 Unknown character set: 'nosuchcharset'
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt CHARSET = 'utf8mb4' COLLATE = 'utf8_bin'
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt CHARACTER SET = 'utf8' CHARACTER SET = 'utf8mb4'
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt DEFAULT CHARSET = 'utf8mb4' DEFAULT COLLATE = 'utf8mb4_general_ci' CHARACTER SET = 'utf8' COLLATE = 'utf8mb4_bin'
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE CHARACTER SET = 'utf8'
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE COLLATE = 'utf8_bin'
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER SCHEMA COLLATE = 'utf8_bin'
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE CHARSET = 'utf8'
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE CHARSET utf8
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE DEFAULT CHARSET = 'utf8'
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER SCHEMA `` DEFAULT COLLATE = 'UTF8_BiN'
< error
| 1102 42000 Incorrect database name ''
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt CHARSET = binary
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt CHARSET binary
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt CHARACTER SET binary
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt COLLATE = `binary`
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt COLLATE = 'binary'
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt COLLATE = binary
< error
| 1064 42000 …near 'binary' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt COLLATE binary
< error
| 1064 42000 …near 'binary' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt CHARSET=binary COLLATE = binary
< error
| 1064 42000 …near 'binary' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ ONLY
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt CHARSET = utf8 READ ONLY
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ ONLY CHARSET = utf8
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec ^
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ ONLY = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ ONLY = 0
< error
| 1064 42000 …near '= 0' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ ONLY = DEFAULT
< error
| 1064 42000 …near '= DEFAULT' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ ONLY 1
< error
| 1064 42000 …near '1' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE READ ONLY = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ
< error
| 1064 42000 …near 'READ' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt DEFAULT READ ONLY = 1
< error
| 1064 42000 …near 'READ ONLY = 1' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ ONLY, CHARSET = utf8
< error
| 1064 42000 …near ', CHARSET = utf8' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt ENCRYPTION = 'N'
< error
| 1064 42000 …near 'ENCRYPTION = 'N'' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt ENCRYPTION 'N'
< error
| 1064 42000 …near 'ENCRYPTION 'N'' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt DEFAULT ENCRYPTION = 'N'
< error
| 1064 42000 …near 'ENCRYPTION = 'N'' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt ENCRYPTION
< error
| 1064 42000 …near 'ENCRYPTION' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE ENCRYPTION = 'N'
< error
| 1064 42000 …near 'ENCRYPTION = 'N'' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt COLLATE
< error
| 1064 42000 …near 'COLLATE' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt COLLATE =
< error
| 1064 42000 …near '=' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt CHARSET =
< error
| 1064 42000 …near '=' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt CHARACTER SET
< error
| 1064 42000 …near 'SET' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt CHARACTER SET =
< error
| 1064 42000 …near '=' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt CHARACTER SET = 'utf8' COLLATE
< error
| 1064 42000 …near 'COLLATE' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt CHARSET = 1
< error
| 1064 42000 …near '1' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt COLLATE = 1
< error
| 1064 42000 …near '1' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt CHARSET = select
< error
| 1064 42000 …near 'select' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt COLLATE = select
< error
| 1064 42000 …near 'select' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt CHARACTER SET = DEFAULT
< error
| 1064 42000 …near 'DEFAULT' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt COLLATE = DEFAULT
< error
| 1064 42000 …near 'DEFAULT' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt CHARSET = 'utf8', COLLATE = 'utf8_bin'
< error
| 1064 42000 …near ', COLLATE = 'utf8_bin'' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt, CHARSET = 'utf8'
< error
| 1064 42000 …near ', CHARSET = 'utf8'' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE
< error
| 1064 42000 …near 'DATABASE' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt
< error
| 1064 42000 …near 'obdbopt' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER SCHEMA
< error
| 1064 42000 …near 'SCHEMA' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER SCHEMA obdbopt
< error
| 1064 42000 …near 'obdbopt' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt DEFAULT
< error
| 1064 42000 …near 'DEFAULT' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt DEFAULT = 'utf8'
< error
| 1064 42000 …near '= 'utf8'' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt CHARSET = utf8 DEFAULT
< error
| 1064 42000 …near 'DEFAULT' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt CHARSET = utf8 CHARACTER
< error
| 1064 42000 …near 'CHARACTER' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE binary CHARSET = 'utf8'
< error
| 1064 42000 …near 'binary CHARSET = 'utf8'' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt COMMENT = 'x'
< error
| 1064 42000 …near 'COMMENT = 'x'' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt COMMENT 'x'
< error
| 1064 42000 …near 'COMMENT 'x'' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE COMMENT = 'x'
< error
| 1064 42000 …near 'COMMENT = 'x'' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt UPGRADE DATA DIRECTORY NAME
< error
| 1064 42000 …near 'UPGRADE DATA DIRECTORY NAME' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE PLACEMENT POLICY = 'p'
< error
| 1064 42000 …near 'POLICY = 'p'' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| CREATE DATABASE obdbopt_c
< ok 1 0
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| CREATE DATABASE obdbopt_c CHARACTER SET = 'utf8'
< error
| 1007 HY000 Can't create database 'obdbopt_c'; database exists
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| CREATE DATABASE obdbopt_c COLLATE = 'utf8_bin'
< error
| 1007 HY000 Can't create database 'obdbopt_c'; database exists
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| CREATE DATABASE obdbopt_c CHARSET = binary
< error
| 1007 HY000 Can't create database 'obdbopt_c'; database exists
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| CREATE SCHEMA obdbopt_c DEFAULT COLLATE = 'UTF8_BiN'
< error
| 1007 HY000 Can't create database 'obdbopt_c'; database exists
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| CREATE DATABASE obdbopt_c DEFAULT CHARSET = 'utf8mb4' COLLATE = 'utf8mb4_bin'
< error
| 1007 HY000 Can't create database 'obdbopt_c'; database exists
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| CREATE DATABASE obdbopt_c CHARSET = 1
< error
| 1064 42000 …near '1' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| CREATE DATABASE obdbopt_c COLLATE = select
< error
| 1064 42000 …near 'select' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| CREATE DATABASE obdbopt_c CHARSET = 'utf8', COLLATE = 'utf8_bin'
< error
| 1064 42000 …near ', COLLATE = 'utf8_bin'' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| CREATE DATABASE obdbopt_c CHARSET = utf8 CHARACTER
< error
| 1064 42000 …near 'CHARACTER' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| CREATE DATABASE obdbopt_c READ ONLY
< error
| 1007 HY000 Can't create database 'obdbopt_c'; database exists
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| CREATE DATABASE obdbopt_c READ WRITE
< error
| 1007 HY000 Can't create database 'obdbopt_c'; database exists
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| CREATE DATABASE obdbopt_c READ ONLY = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| CREATE DATABASE obdbopt_c READ ONLY 0
< error
| 1064 42000 …near '0' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| CREATE DATABASE obdbopt_c READ ONLY = DEFAULT
< error
| 1064 42000 …near '= DEFAULT' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| CREATE DATABASE obdbopt_c ENCRYPTION = 'N'
< error
| 1064 42000 …near 'ENCRYPTION = 'N'' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| CREATE DATABASE obdbopt_c DEFAULT ENCRYPTION = 'Y'
< error
| 1064 42000 …near 'ENCRYPTION = 'Y'' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| CREATE DATABASE obdbopt_c COMMENT = 'hello'
< error
| 1064 42000 …near 'COMMENT = 'hello'' at line 1
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| USE obdbopt
< ok 0 0
> exec
| ALTER DATABASE obdbopt READ WRITE
< ok 0 0
> exec
| DROP DATABASE IF EXISTS obdbopt
< ok 0 0
> exec
| DROP DATABASE IF EXISTS obdbopt_c
< ok 0 0
