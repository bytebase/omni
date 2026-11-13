-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestContainer_Section_1_19_TableOptions default#1
> exec
| DROP TABLE IF EXISTS t_eng_innodb
< ok 0 0
> exec
| CREATE TABLE t_eng_innodb (a INT) ENGINE=InnoDB
< ok 0 0
> query
| SHOW CREATE TABLE t_eng_innodb
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_eng_innodb
: bytes
| CREATE TABLE `t_eng_innodb` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_eng_myisam
< ok 0 0
> exec
| CREATE TABLE t_eng_myisam (a INT) ENGINE=MyISAM
< ok 0 0
> query
| SHOW CREATE TABLE t_eng_myisam
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_eng_myisam
: bytes
| CREATE TABLE `t_eng_myisam` (
|   `a` int DEFAULT NULL
| ) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_eng_memory
< ok 0 0
> exec
| CREATE TABLE t_eng_memory (a INT) ENGINE=MEMORY
< ok 0 0
> query
| SHOW CREATE TABLE t_eng_memory
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_eng_memory
: bytes
| CREATE TABLE `t_eng_memory` (
|   `a` int DEFAULT NULL
| ) ENGINE=MEMORY DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_cs_utf8mb4
< ok 0 0
> exec
| CREATE TABLE t_cs_utf8mb4 (a INT) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
< ok 0 0
> query
| SHOW CREATE TABLE t_cs_utf8mb4
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_cs_utf8mb4
: bytes
| CREATE TABLE `t_cs_utf8mb4` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_cs_latin1
< ok 0 0
> exec
| CREATE TABLE t_cs_latin1 (a INT) DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci
< ok 0 0
> query
| SHOW CREATE TABLE t_cs_latin1
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_cs_latin1
: bytes
| CREATE TABLE `t_cs_latin1` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=latin1
> exec
| DROP TABLE IF EXISTS t_cs_unicode
< ok 0 0
> exec
| CREATE TABLE t_cs_unicode (a INT) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
< ok 0 0
> query
| SHOW CREATE TABLE t_cs_unicode
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_cs_unicode
: bytes
| CREATE TABLE `t_cs_unicode` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
> exec
| DROP TABLE IF EXISTS t_comment
< ok 0 0
> exec
| CREATE TABLE t_comment (a INT) COMMENT='table description'
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
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='table description'
> exec
| DROP TABLE IF EXISTS t_comment_sc
< ok 0 0
> exec
| CREATE TABLE t_comment_sc (a INT) COMMENT='it\'s a \\test'
< ok 0 0
> query
| SHOW CREATE TABLE t_comment_sc
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_comment_sc
: bytes
| CREATE TABLE `t_comment_sc` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='it''s a \\test'
> exec
| DROP TABLE IF EXISTS t_rf_dyn
< ok 0 0
> exec
| CREATE TABLE t_rf_dyn (a INT) ROW_FORMAT=DYNAMIC
< ok 0 0
> query
| SHOW CREATE TABLE t_rf_dyn
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_rf_dyn
: bytes
| CREATE TABLE `t_rf_dyn` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC
> exec
| DROP TABLE IF EXISTS t_rf_comp
< ok 0 0
> exec
| CREATE TABLE t_rf_comp (a INT) ROW_FORMAT=COMPRESSED
< ok 0 0
> query
| SHOW CREATE TABLE t_rf_comp
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_rf_comp
: bytes
| CREATE TABLE `t_rf_comp` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=COMPRESSED
> exec
| DROP TABLE IF EXISTS t_ai1000
< ok 0 0
> exec
| CREATE TABLE t_ai1000 (id INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (id)) AUTO_INCREMENT=1000
< ok 0 0
> query
| SHOW CREATE TABLE t_ai1000
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_ai1000
: bytes
| CREATE TABLE `t_ai1000` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_kbs8
< ok 0 0
> exec
| CREATE TABLE t_kbs8 (a INT) KEY_BLOCK_SIZE=8
< ok 0 0
> query
| SHOW CREATE TABLE t_kbs8
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_kbs8
: bytes
| CREATE TABLE `t_kbs8` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci KEY_BLOCK_SIZE=8
> exec
| DROP TABLE IF EXISTS t_multi_opts
< ok 0 0
> exec
| CREATE TABLE t_multi_opts (id INT NOT NULL AUTO_INCREMENT, name VARCHAR(100), PRIMARY KEY (id)) ENGINE=InnoDB AUTO_INCREMENT=500 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='multi opts'
< ok 0 0
> query
| SHOW CREATE TABLE t_multi_opts
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_multi_opts
: bytes
| CREATE TABLE `t_multi_opts` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB AUTO_INCREMENT=500 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='multi opts'
