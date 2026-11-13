-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestJSONValueStoredForms default#1
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' )))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512))) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' ) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512)) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING SIGNED)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning signed)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING SIGNED) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning signed) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING SIGNED INTEGER)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning signed)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING SIGNED INTEGER) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning signed) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING UNSIGNED)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning unsigned)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING UNSIGNED) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning unsigned) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING UNSIGNED INTEGER)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning unsigned)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING UNSIGNED INTEGER) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning unsigned) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING CHAR)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(1073741823))) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING CHAR) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(1073741823)) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING CHAR(3))))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(3))) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING CHAR(3)) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(3)) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING CHAR(3) CHARSET latin1)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(3) character set latin1)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING CHAR(3) CHARSET latin1) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(3) character set latin1) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING CHAR CHARSET latin1)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(4294967295) character set latin1)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING CHAR CHARSET latin1) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(4294967295) character set latin1) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING CHAR CHARSET utf8mb3)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(1431655764) character set utf8mb3)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING CHAR CHARSET utf8mb3) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(1431655764) character set utf8mb3) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING CHAR CHARSET utf8)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(1431655764) character set utf8mb3)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING CHAR CHARSET utf8) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(1431655764) character set utf8mb3) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING CHAR CHARSET gbk)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(2147483647) character set gbk)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING CHAR CHARSET gbk) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(2147483647) character set gbk) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING CHAR CHARSET utf32)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(1073741823) character set utf32)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING CHAR CHARSET utf32) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(1073741823) character set utf32) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING CHAR CHARSET binary)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning binary(4294967295))) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING CHAR CHARSET binary) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning binary(4294967295)) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING CHAR(3) CHARSET binary)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning binary(3))) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING CHAR(3) CHARSET binary) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning binary(3)) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING CHAR BINARY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(1073741823) character set utf8mb4)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING CHAR BINARY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(1073741823) character set utf8mb4) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING NCHAR)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(1431655764) character set utf8mb3)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING NCHAR) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(1431655764) character set utf8mb3) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING CHAR(3) BINARY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(3) character set utf8mb4)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING CHAR(3) BINARY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(3) character set utf8mb4) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING BINARY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning binary(4294967295))) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING BINARY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning binary(4294967295)) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING BINARY(3))))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning binary(3))) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING BINARY(3)) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning binary(3)) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING NCHAR(3))))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(3) character set utf8mb3)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING NCHAR(3)) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(3) character set utf8mb3) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING DECIMAL)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning decimal(10, 0))) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING DECIMAL) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning decimal(10, 0)) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING DECIMAL(5))))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning decimal(5, 0))) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING DECIMAL(5)) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning decimal(5, 0)) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING DECIMAL(5,2))))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning decimal(5, 2))) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING DECIMAL(5,2)) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning decimal(5, 2)) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING FLOAT)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning float)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING FLOAT) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning float) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING FLOAT(5))))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning float)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING FLOAT(5)) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning float) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING FLOAT(30))))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning float)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING FLOAT(30)) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning float) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING DOUBLE)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning double)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING DOUBLE) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning double) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING REAL)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning double)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING REAL) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning double) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING DATE)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning date)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING DATE) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning date) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING TIME)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning time)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING TIME) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning time) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING TIME(2))))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning time(2))) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING TIME(2)) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning time(2)) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING DATETIME)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning datetime)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING DATETIME) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning datetime) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING DATETIME(3))))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning datetime(3))) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING DATETIME(3)) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning datetime(3)) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING YEAR)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning year)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING YEAR) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning year) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING JSON)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning json)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING JSON) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning json) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' NULL ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) null on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' NULL ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) null on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' ERROR ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) error on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' ERROR ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) error on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' NULL ON ERROR)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) null on error)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' NULL ON ERROR) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) null on error) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' ERROR ON ERROR)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) error on error)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' ERROR ON ERROR) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) error on error) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' NULL ON EMPTY ERROR ON ERROR)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) null on empty error on error)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' NULL ON EMPTY ERROR ON ERROR) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) null on empty error on error) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT 'x' ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default _utf8mb4'x' on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT 'x' ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default 'x' on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT '001' ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default _utf8mb4'001' on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT '001' ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default '001' on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT '' ON ERROR)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default _utf8mb4'' on error)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT '' ON ERROR) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default '' on error) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT 'it''s' ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default _utf8mb4'it\'s' on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT 'it''s' ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default 'it\'s' on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT 'a\\b' ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default _utf8mb4'a\\b' on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT 'a\\b' ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default 'a\\b' on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT _latin1'x' ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default _latin1'x' on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT _latin1'x' ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default _latin1'x' on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT N'x' ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default _utf8mb3'x' on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT N'x' ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default 'x' on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT 7 ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default 7 on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT 7 ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default 7 on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT -7 ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default -7 on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT -7 ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default -7 on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT +7 ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default 7 on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT +7 ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default 7 on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT 1.50 ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default 1.50 on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT 1.50 ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default 1.50 on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT -1.5 ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default -1.5 on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT -1.5 ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default -1.5 on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT .5 ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default 0.5 on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT .5 ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default 0.5 on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT 1E2 ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default 1E2 on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT 1E2 ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default 1E2 on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT TRUE ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default true on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT TRUE ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default true on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT FALSE ON ERROR)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default false on error)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT FALSE ON ERROR) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default false on error) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT x'61' ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default 0x61 on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT x'61' ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default 0x61 on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT X'6A6b' ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default 0x6a6b on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT X'6A6b' ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default 0x6a6b on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT 0xABC ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default 0x0abc on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT 0xABC ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default 0x0abc on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT 0x6a ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default 0x6a on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT 0x6a ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default 0x6a on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT b'1' ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default 0x01 on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT b'1' ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default 0x01 on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT b'101010101' ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default 0x0155 on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT b'101010101' ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default 0x0155 on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT 0b11 ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default 0x03 on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT 0b11 ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default 0x03 on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT DATE'2000-01-01' ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default DATE'2000-01-01' on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT DATE'2000-01-01' ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default DATE'2000-01-01' on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT TIME'10:00:00' ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default TIME'10:00:00' on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT TIME'10:00:00' ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default TIME'10:00:00' on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT TIMESTAMP'2000-01-01 00:00:00' ON ERROR)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default TIMESTAMP'2000-01-01 00:00:00' on error)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT TIMESTAMP'2000-01-01 00:00:00' ON ERROR) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default TIMESTAMP'2000-01-01 00:00:00' on error) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING CHAR(0))))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(0))) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING CHAR(0)) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(0)) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING BINARY(0))))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning binary(0))) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING BINARY(0)) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning binary(0)) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING DECIMAL(0))))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning decimal(10, 0))) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING DECIMAL(0)) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning decimal(10, 0)) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT 18446744073709551615 ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default 18446744073709551615 on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT 18446744073709551615 ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default 18446744073709551615 on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT -9223372036854775808 ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default -9223372036854775808 on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT -9223372036854775808 ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default -9223372036854775808 on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT 99999999999999999999 ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default 99999999999999999999 on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT 99999999999999999999 ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default 99999999999999999999 on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT 007 ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default 7 on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT 007 ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default 7 on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT b'' ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default X'' on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT b'' ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default X'' on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT b'0000000001000001' ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default 0x0041 on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT b'0000000001000001' ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default 0x0041 on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT x'' ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default X'' on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT x'' ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default X'' on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT 0x0041 ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default 0x0041 on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT 0x0041 ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default 0x0041 on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT _latin1 X'61' ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default _latin1'a' on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT _latin1 X'61' ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default _latin1'a' on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT _latin1 b'01100001' ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default _latin1'a' on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT _latin1 b'01100001' ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default _latin1'a' on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT _utf8mb4 0x61 ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning char(512) default _utf8mb4'a' on empty)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT _utf8mb4 0x61 ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning char(512) default _utf8mb4'a' on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING SIGNED DEFAULT 7 ON EMPTY DEFAULT -1 ON ERROR)))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning signed default 7 on empty default -1 on error)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING SIGNED DEFAULT 7 ON EMPTY DEFAULT -1 ON ERROR) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select json_value('{}', '$.b' returning signed default 7 on empty default -1 on error) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
