-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestContainer_Section_1_2_StringTypes default#1
> exec
| DROP TABLE IF EXISTS t_char10
< ok 0 0
> exec
| CREATE TABLE t_char10 (a CHAR(10))
< ok 0 0
> query
| SHOW CREATE TABLE t_char10
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_char10
: bytes
| CREATE TABLE `t_char10` (
|   `a` char(10) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t_char1
< ok 0 0
> exec
| CREATE TABLE t_char1 (a CHAR)
< ok 0 0
> query
| SHOW CREATE TABLE t_char1
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_char1
: bytes
| CREATE TABLE `t_char1` (
|   `a` char(1) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t_varchar255
< ok 0 0
> exec
| CREATE TABLE t_varchar255 (a VARCHAR(255))
< ok 0 0
> query
| SHOW CREATE TABLE t_varchar255
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_varchar255
: bytes
| CREATE TABLE `t_varchar255` (
|   `a` varchar(255) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t_varchar16383
< ok 0 0
> exec
| CREATE TABLE t_varchar16383 (a VARCHAR(16383))
< ok 0 0
> query
| SHOW CREATE TABLE t_varchar16383
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_varchar16383
: bytes
| CREATE TABLE `t_varchar16383` (
|   `a` varchar(16383) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t_tinytext
< ok 0 0
> exec
| CREATE TABLE t_tinytext (a TINYTEXT)
< ok 0 0
> query
| SHOW CREATE TABLE t_tinytext
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_tinytext
: bytes
| CREATE TABLE `t_tinytext` (
|   `a` tinytext DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t_text
< ok 0 0
> exec
| CREATE TABLE t_text (a TEXT)
< ok 0 0
> query
| SHOW CREATE TABLE t_text
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_text
: bytes
| CREATE TABLE `t_text` (
|   `a` text DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t_mediumtext
< ok 0 0
> exec
| CREATE TABLE t_mediumtext (a MEDIUMTEXT)
< ok 0 0
> query
| SHOW CREATE TABLE t_mediumtext
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_mediumtext
: bytes
| CREATE TABLE `t_mediumtext` (
|   `a` mediumtext DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t_longtext
< ok 0 0
> exec
| CREATE TABLE t_longtext (a LONGTEXT)
< ok 0 0
> query
| SHOW CREATE TABLE t_longtext
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_longtext
: bytes
| CREATE TABLE `t_longtext` (
|   `a` longtext DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t_text1000
< ok 0 0
> exec
| CREATE TABLE t_text1000 (a TEXT(1000))
< ok 0 0
> query
| SHOW CREATE TABLE t_text1000
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_text1000
: bytes
| CREATE TABLE `t_text1000` (
|   `a` text DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t_enum
< ok 0 0
> exec
| CREATE TABLE t_enum (a ENUM('a','b','c'))
< ok 0 0
> query
| SHOW CREATE TABLE t_enum
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_enum
: bytes
| CREATE TABLE `t_enum` (
|   `a` enum('a','b','c') DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t_enum_sc
< ok 0 0
> exec
| CREATE TABLE t_enum_sc (a ENUM('it''s','hello,world','a"b'))
< ok 0 0
> query
| SHOW CREATE TABLE t_enum_sc
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_enum_sc
: bytes
| CREATE TABLE `t_enum_sc` (
|   `a` enum('it''s','hello,world','a"b') DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t_set
< ok 0 0
> exec
| CREATE TABLE t_set (a SET('x','y','z'))
< ok 0 0
> query
| SHOW CREATE TABLE t_set
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_set
: bytes
| CREATE TABLE `t_set` (
|   `a` set('x','y','z') DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t_char_cs
< ok 0 0
> exec
| CREATE TABLE t_char_cs (a CHAR(10) CHARACTER SET latin1)
< ok 0 0
> query
| SHOW CREATE TABLE t_char_cs
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_char_cs
: bytes
| CREATE TABLE `t_char_cs` (
|   `a` char(10) CHARACTER SET latin1 COLLATE latin1_bin DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t_varchar_cc
< ok 0 0
> exec
| CREATE TABLE t_varchar_cc (a VARCHAR(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_bin)
< ok 0 0
> query
| SHOW CREATE TABLE t_varchar_cc
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_varchar_cc
: bytes
| CREATE TABLE `t_varchar_cc` (
|   `a` varchar(100) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t_nchar
< ok 0 0
> exec
| CREATE TABLE t_nchar (a NATIONAL CHAR(10))
< ok 0 0
> query
| SHOW CREATE TABLE t_nchar
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_nchar
: bytes
| CREATE TABLE `t_nchar` (
|   `a` char(10) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t_nchar2
< ok 0 0
> exec
| CREATE TABLE t_nchar2 (a NCHAR(10))
< ok 0 0
> query
| SHOW CREATE TABLE t_nchar2
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_nchar2
: bytes
| CREATE TABLE `t_nchar2` (
|   `a` char(10) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t_nvarchar
< ok 0 0
> exec
| CREATE TABLE t_nvarchar (a NVARCHAR(100))
< ok 0 0
> query
| SHOW CREATE TABLE t_nvarchar
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_nvarchar
: bytes
| CREATE TABLE `t_nvarchar` (
|   `a` varchar(100) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
