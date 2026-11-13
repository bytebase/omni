-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestContainer_Section_2_5_AlterTableTableLevel default#1
> exec
| DROP TABLE IF EXISTS t_rename_src
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_rename_dst
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_engine
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_conv_cs
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_def_cs
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_comment
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_autoinc
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_rowfmt
< ok 0 0
> exec
| CREATE TABLE t_rename_src (id INT PRIMARY KEY, name VARCHAR(100))
< ok 0 0
> exec
| ALTER TABLE t_rename_src RENAME TO t_rename_dst
< ok 0 0
> query
| SHOW CREATE TABLE t_rename_dst
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_rename_dst
: bytes
| CREATE TABLE `t_rename_dst` (
|   `id` int NOT NULL,
|   `name` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_rename_src
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_rename_dst
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_engine
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_conv_cs
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_def_cs
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_comment
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_autoinc
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_rowfmt
< ok 0 0
> exec
| CREATE TABLE t_engine (id INT PRIMARY KEY, val INT)
< ok 0 0
> exec
| ALTER TABLE t_engine ENGINE=MyISAM
< ok 0 0
> query
| SHOW CREATE TABLE t_engine
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_engine
: bytes
| CREATE TABLE `t_engine` (
|   `id` int NOT NULL,
|   `val` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_rename_src
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_rename_dst
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_engine
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_conv_cs
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_def_cs
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_comment
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_autoinc
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_rowfmt
< ok 0 0
> exec
| CREATE TABLE t_conv_cs (id INT PRIMARY KEY, name VARCHAR(100)) DEFAULT CHARSET=latin1
< ok 0 0
> exec
| ALTER TABLE t_conv_cs CONVERT TO CHARACTER SET utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE t_conv_cs
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_conv_cs
: bytes
| CREATE TABLE `t_conv_cs` (
|   `id` int NOT NULL,
|   `name` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_rename_src
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_rename_dst
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_engine
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_conv_cs
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_def_cs
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_comment
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_autoinc
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_rowfmt
< ok 0 0
> exec
| CREATE TABLE t_def_cs (id INT PRIMARY KEY, name VARCHAR(100))
< ok 0 0
> exec
| ALTER TABLE t_def_cs DEFAULT CHARACTER SET latin1
< ok 0 0
> query
| SHOW CREATE TABLE t_def_cs
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_def_cs
: bytes
| CREATE TABLE `t_def_cs` (
|   `id` int NOT NULL,
|   `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=latin1
> exec
| DROP TABLE IF EXISTS t_rename_src
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_rename_dst
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_engine
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_conv_cs
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_def_cs
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_comment
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_autoinc
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_rowfmt
< ok 0 0
> exec
| CREATE TABLE t_comment (id INT PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE t_comment COMMENT='new comment'
< ok 0 0
> query
| SHOW CREATE TABLE t_comment
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_comment
: bytes
| CREATE TABLE `t_comment` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='new comment'
> exec
| DROP TABLE IF EXISTS t_rename_src
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_rename_dst
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_engine
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_conv_cs
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_def_cs
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_comment
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_autoinc
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_rowfmt
< ok 0 0
> exec
| CREATE TABLE t_autoinc (id INT PRIMARY KEY AUTO_INCREMENT, val INT)
< ok 0 0
> exec
| ALTER TABLE t_autoinc AUTO_INCREMENT=1000
< ok 0 0
> query
| SHOW CREATE TABLE t_autoinc
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_autoinc
: bytes
| CREATE TABLE `t_autoinc` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `val` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_rename_src
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_rename_dst
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_engine
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_conv_cs
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_def_cs
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_comment
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_autoinc
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_rowfmt
< ok 0 0
> exec
| CREATE TABLE t_rowfmt (id INT PRIMARY KEY, val INT)
< ok 0 0
> exec
| ALTER TABLE t_rowfmt ROW_FORMAT=COMPRESSED
< ok 0 0
> query
| SHOW CREATE TABLE t_rowfmt
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_rowfmt
: bytes
| CREATE TABLE `t_rowfmt` (
|   `id` int NOT NULL,
|   `val` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=COMPRESSED
