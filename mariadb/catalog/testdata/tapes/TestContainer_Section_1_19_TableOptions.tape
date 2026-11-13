-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

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
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
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
|   `a` int(11) DEFAULT NULL
| ) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
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
|   `a` int(11) DEFAULT NULL
| ) ENGINE=MEMORY DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
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
|   `a` int(11) DEFAULT NULL
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
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci
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
|   `a` int(11) DEFAULT NULL
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
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci COMMENT='table description'
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
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci COMMENT='it''s a \\test'
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
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci ROW_FORMAT=DYNAMIC
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
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci ROW_FORMAT=COMPRESSED
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
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
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
|   `a` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci KEY_BLOCK_SIZE=8
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
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `name` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB AUTO_INCREMENT=500 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='multi opts'
