-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestContainer_Section_1_20_CharsetCollationInheritance default#1
> exec
| DROP DATABASE IF EXISTS db_latin1
< ok 0 0
> exec
| CREATE DATABASE db_latin1 CHARACTER SET latin1
< ok 0 0
> exec
| USE db_latin1
< ok 0 0
> exec
| CREATE TABLE t_inherit_db (a VARCHAR(100))
< ok 0 0
> query
| SHOW CREATE TABLE t_inherit_db
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_inherit_db
: bytes
| CREATE TABLE `t_inherit_db` (
|   `a` varchar(100) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_bin
> exec
| DROP DATABASE IF EXISTS db_utf8mb4_tbl
< ok 0 0
> exec
| CREATE DATABASE db_utf8mb4_tbl
< ok 0 0
> exec
| USE db_utf8mb4_tbl
< ok 0 0
> exec
| CREATE TABLE t_col_inherit (a VARCHAR(50)) DEFAULT CHARSET=latin1
< ok 0 0
> query
| SHOW CREATE TABLE t_col_inherit
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_col_inherit
: bytes
| CREATE TABLE `t_col_inherit` (
|   `a` varchar(50) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_bin
> exec
| DROP DATABASE IF EXISTS db_override_cs
< ok 0 0
> exec
| CREATE DATABASE db_override_cs
< ok 0 0
> exec
| USE db_override_cs
< ok 0 0
> exec
| CREATE TABLE t_col_override_cs (a VARCHAR(50) CHARACTER SET latin1) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE t_col_override_cs
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_col_override_cs
: bytes
| CREATE TABLE `t_col_override_cs` (
|   `a` varchar(50) CHARACTER SET latin1 COLLATE latin1_bin DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS db_override_coll
< ok 0 0
> exec
| CREATE DATABASE db_override_coll
< ok 0 0
> exec
| USE db_override_coll
< ok 0 0
> exec
| CREATE TABLE t_col_override_coll (a VARCHAR(50) COLLATE utf8mb4_bin) DEFAULT CHARSET=utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE t_col_override_coll
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_col_override_coll
: bytes
| CREATE TABLE `t_col_override_coll` (
|   `a` varchar(50) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS db_display
< ok 0 0
> exec
| CREATE DATABASE db_display CHARACTER SET utf8mb4
< ok 0 0
> exec
| USE db_display
< ok 0 0
> exec
| CREATE TABLE t_display_rules (a VARCHAR(50), b VARCHAR(50) CHARACTER SET latin1, c VARCHAR(50) COLLATE utf8mb4_bin, d VARCHAR(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_bin)
< ok 0 0
> query
| SHOW CREATE TABLE t_display_rules
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_display_rules
: bytes
| CREATE TABLE `t_display_rules` (
|   `a` varchar(50) DEFAULT NULL,
|   `b` varchar(50) CHARACTER SET latin1 COLLATE latin1_bin DEFAULT NULL,
|   `c` varchar(50) DEFAULT NULL,
|   `d` varchar(50) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS db_binary
< ok 0 0
> exec
| CREATE DATABASE db_binary
< ok 0 0
> exec
| USE db_binary
< ok 0 0
> exec
| CREATE TABLE t_binary_cs (a VARCHAR(50) CHARACTER SET binary)
< ok 0 0
> query
| SHOW CREATE TABLE t_binary_cs
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_binary_cs
: bytes
| CREATE TABLE `t_binary_cs` (
|   `a` varbinary(50) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
