-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_JSONValueStoredForms default#1
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 0 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' )))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' ) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING SIGNED)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING signed)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING SIGNED) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING signed) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING SIGNED INTEGER)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING signed)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING SIGNED INTEGER) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING signed) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING UNSIGNED)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING unsigned)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING UNSIGNED) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING unsigned) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING UNSIGNED INTEGER)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING unsigned)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING UNSIGNED INTEGER) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING unsigned) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING CHAR)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char charset utf8mb4)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING CHAR) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char charset utf8mb4) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING CHAR(3))))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(3) charset utf8mb4)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING CHAR(3)) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(3) charset utf8mb4) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING CHAR(3) CHARSET ascii)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(3) charset ascii)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING CHAR(3) CHARSET ascii) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(3) charset ascii) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING CHAR CHARSET ascii)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char charset ascii)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING CHAR CHARSET ascii) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char charset ascii) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING CHAR(3) BINARY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(3) charset utf8mb4)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING CHAR(3) BINARY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(3) charset utf8mb4) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING BINARY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char charset binary)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING BINARY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char charset binary) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING BINARY(3))))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(3) charset binary)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING BINARY(3)) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(3) charset binary) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING NCHAR(3))))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(3) charset utf8mb4)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING NCHAR(3)) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(3) charset utf8mb4) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING CHAR(3) CHARSET utf8mb3)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(3) charset utf8mb4)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING CHAR(3) CHARSET utf8mb3) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(3) charset utf8mb4) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING CHAR(3) CHARSET utf8)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(3) charset utf8mb4)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING CHAR(3) CHARSET utf8) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(3) charset utf8mb4) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING DECIMAL)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING number(10,0))) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING DECIMAL) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING number(10,0)) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING DECIMAL(5))))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING number(5,0))) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING DECIMAL(5)) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING number(5,0)) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING DECIMAL(5,2))))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING number(5,2))) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING DECIMAL(5,2)) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING number(5,2)) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING NUMBER)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING number(10,0))) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING NUMBER) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING number(10,0)) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING NUMBER(5,2))))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING number(5,2))) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING NUMBER(5,2)) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING number(5,2)) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING NUMERIC(5,2))))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING number(5,2))) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING NUMERIC(5,2)) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING number(5,2)) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING FLOAT)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING float)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING FLOAT) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING float) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING FLOAT(5))))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING float)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING FLOAT(5)) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING float) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING FLOAT(30))))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING float)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING FLOAT(30)) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING float) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING DOUBLE)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING double)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING DOUBLE) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING double) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING DATE)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING date)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING DATE) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING date) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING TIME)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING time(0))) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING TIME) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING time(0)) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING TIME(2))))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING time(2))) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING TIME(2)) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING time(2)) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING DATETIME)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING datetime(0))) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING DATETIME) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING datetime(0)) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING DATETIME(3))))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING datetime(3))) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING DATETIME(3)) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING datetime(3)) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING YEAR)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING year(0))) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING YEAR) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING year(0)) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING JSON)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING json)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING JSON) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING json) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING CHAR(3) ASCII)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(3) charset utf8mb4 ASCII)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING CHAR(3) ASCII) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(3) charset utf8mb4 ASCII) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING CHAR(3) CHARSET ascii ASCII)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(3) charset ascii ASCII)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING CHAR(3) CHARSET ascii ASCII) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(3) charset ascii ASCII) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING CHAR(3) TRUNCATE)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(3) charset utf8mb4 TRUNCATE)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING CHAR(3) TRUNCATE) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(3) charset utf8mb4 TRUNCATE) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING CHAR(3) TRUNCATE ASCII)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(3) charset utf8mb4 TRUNCATE ASCII)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING CHAR(3) TRUNCATE ASCII) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(3) charset utf8mb4 TRUNCATE ASCII) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' NULL ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 null on empty)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' NULL ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 null on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' ERROR ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 error on empty)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' ERROR ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 error on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' NULL ON ERROR)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 null on error)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' NULL ON ERROR) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 null on error) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' ERROR ON ERROR)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 error on error)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' ERROR ON ERROR) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 error on error) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' NULL ON EMPTY ERROR ON ERROR)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 null on empty error on error)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' NULL ON EMPTY ERROR ON ERROR) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 null on empty error on error) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT 'x' ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 default 'x' on empty)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT 'x' ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 default 'x' on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT '001' ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 default '001' on empty)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT '001' ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 default '001' on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT '' ON ERROR)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 default '' on error)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT '' ON ERROR) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 default '' on error) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT 'it''s' ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 default 'it\'s' on empty)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT 'it''s' ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 default 'it\'s' on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT 'a\\b' ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 default 'a\\b' on empty)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT 'a\\b' ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 default 'a\\b' on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT _utf8mb4'x' ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 default 'x' on empty)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT _utf8mb4'x' ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 default 'x' on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT N'x' ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 default 'x' on empty)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT N'x' ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 default 'x' on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT 7 ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 default 7 on empty)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT 7 ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 default 7 on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT -7 ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 default -7 on empty)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT -7 ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 default -7 on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT +7 ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 default 7 on empty)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT +7 ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 default 7 on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT 1.50 ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 default 1.50 on empty)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT 1.50 ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 default 1.50 on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT -1.5 ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 default -1.5 on empty)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT -1.5 ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 default -1.5 on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT .5 ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 default 0.5 on empty)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT .5 ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 default 0.5 on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT x'61' ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 default X'61' on empty)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT x'61' ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 default X'61' on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT X'6A6b' ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 default X'6A6B' on empty)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT X'6A6b' ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 default X'6A6B' on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT 0xABC ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 default X'0ABC' on empty)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT 0xABC ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 default X'0ABC' on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT 0x6a ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 default X'6A' on empty)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT 0x6a ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 default X'6A' on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT b'1' ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 default X'01' on empty)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT b'1' ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 default X'01' on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT b'101010101' ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 default X'0155' on empty)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT b'101010101' ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 default X'0155' on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT 0b11 ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 default X'03' on empty)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT 0b11 ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 default X'03' on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT DATE'2000-01-01' ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 default DATE '2000-01-01' on empty)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT DATE'2000-01-01' ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 default DATE '2000-01-01' on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT TIME'10:00:00' ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 default TIME '10:00:00' on empty)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT TIME'10:00:00' ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 default TIME '10:00:00' on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT TIMESTAMP'2000-01-01 00:00:00' ON ERROR)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 default TIMESTAMP '2000-01-01 00:00:00' on error)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT TIMESTAMP'2000-01-01 00:00:00' ON ERROR) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 default TIMESTAMP '2000-01-01 00:00:00' on error) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING CHAR(0))))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(0) charset utf8mb4)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING CHAR(0)) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(0) charset utf8mb4) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING BINARY(0))))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(0) charset binary)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING BINARY(0)) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(0) charset binary) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT 18446744073709551615 ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 default 18446744073709551615 on empty)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT 18446744073709551615 ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 default 18446744073709551615 on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT -9223372036854775808 ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 default -9223372036854775808 on empty)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT -9223372036854775808 ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 default -9223372036854775808 on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT 99999999999999999999 ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 default 99999999999999999999 on empty)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT 99999999999999999999 ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 default 99999999999999999999 on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT 007 ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 default 7 on empty)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT 007 ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 default 7 on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT b'' ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 default X'' on empty)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT b'' ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 default X'' on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT b'0000000001000001' ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 default X'0041' on empty)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT b'0000000001000001' ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 default X'0041' on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT x'' ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 default X'' on empty)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT x'' ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 default X'' on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT 0x0041 ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 default X'0041' on empty)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT 0x0041 ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 default X'0041' on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT _latin1 X'61' ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 default 'a' on empty)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT _latin1 X'61' ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 default 'a' on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT _latin1 b'01100001' ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 default 'a' on empty)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT _latin1 b'01100001' ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 default 'a' on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' DEFAULT _utf8mb4 0x61 ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING char(512) charset utf8mb4 default 'a' on empty)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' DEFAULT _utf8mb4 0x61 ON EMPTY) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING char(512) charset utf8mb4 default 'a' on empty) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_value
< ok 4 0
> exec
| CREATE DATABASE zz_json_value
< ok 1 0
> exec
| USE zz_json_value
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' RETURNING SIGNED DEFAULT 7 ON EMPTY DEFAULT -1 ON ERROR)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_value(`a`,'$.b' RETURNING signed default 7 on empty default -1 on error)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_VALUE('{}', '$.b' RETURNING SIGNED DEFAULT 7 ON EMPTY DEFAULT -1 ON ERROR) AS x
< ok 0 0
> query
| SHOW CREATE VIEW v
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` AS select json_value('{}','$.b' RETURNING signed default 7 on empty default -1 on error) AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
