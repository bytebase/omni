-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_JSONQueryStoredForms default#1
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 0 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING JSON)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING json null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING JSON) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING json null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING JSON TRUNCATE)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING json TRUNCATE null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING JSON TRUNCATE) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING json TRUNCATE null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING JSON ALLOW SCALARS)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING json ALLOW SCALARS null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING JSON ALLOW SCALARS) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING json ALLOW SCALARS null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING JSON DISALLOW SCALARS)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING json DISALLOW SCALARS null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING JSON DISALLOW SCALARS) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING json DISALLOW SCALARS null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING JSON PRETTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING json PRETTY null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING JSON PRETTY) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING json PRETTY null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING JSON ASCII)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING json ASCII null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING JSON ASCII) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING json ASCII null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING JSON ASIS)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING json asis null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING JSON ASIS) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING json asis null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING JSON MULTIVALUE)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING json null on mismatch multivalue)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING JSON MULTIVALUE) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING json null on mismatch multivalue) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING JSON WITHOUT WRAPPER)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING json WITHOUT WRAPPER null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING JSON WITHOUT WRAPPER) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING json WITHOUT WRAPPER null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING JSON WITHOUT ARRAY WRAPPER)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING json WITHOUT ARRAY WRAPPER null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING JSON WITHOUT ARRAY WRAPPER) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING json WITHOUT ARRAY WRAPPER null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING JSON WITH WRAPPER)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING json WITH WRAPPER null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING JSON WITH WRAPPER) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING json WITH WRAPPER null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING JSON WITH ARRAY WRAPPER)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING json WITH ARRAY WRAPPER null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING JSON WITH ARRAY WRAPPER) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING json WITH ARRAY WRAPPER null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING JSON WITH UNCONDITIONAL WRAPPER)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING json WITH UNCONDITIONAL WRAPPER null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING JSON WITH UNCONDITIONAL WRAPPER) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING json WITH UNCONDITIONAL WRAPPER null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING JSON WITH CONDITIONAL WRAPPER)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING json WITH CONDITIONAL WRAPPER null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING JSON WITH CONDITIONAL WRAPPER) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING json WITH CONDITIONAL WRAPPER null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING JSON WITH UNCONDITIONAL ARRAY WRAPPER)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING json WITH UNCONDITIONAL ARRAY WRAPPER null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING JSON WITH UNCONDITIONAL ARRAY WRAPPER) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING json WITH UNCONDITIONAL ARRAY WRAPPER null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING JSON WITH CONDITIONAL ARRAY WRAPPER)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING json WITH CONDITIONAL ARRAY WRAPPER null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING JSON WITH CONDITIONAL ARRAY WRAPPER) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING json WITH CONDITIONAL ARRAY WRAPPER null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING JSON ERROR ON ERROR)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING json error on error null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING JSON ERROR ON ERROR) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING json error on error null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING JSON NULL ON ERROR)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING json null on error null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING JSON NULL ON ERROR) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING json null on error null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING JSON EMPTY ON ERROR)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING json EMPTY on error null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING JSON EMPTY ON ERROR) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING json EMPTY on error null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING JSON EMPTY ARRAY ON ERROR)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING json EMPTY ARRAY on error null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING JSON EMPTY ARRAY ON ERROR) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING json EMPTY ARRAY on error null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING JSON EMPTY OBJECT ON ERROR)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING json EMPTY OBJECT on error null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING JSON EMPTY OBJECT ON ERROR) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING json EMPTY OBJECT on error null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING JSON ERROR ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING json error on empty null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING JSON ERROR ON EMPTY) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING json error on empty null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING JSON NULL ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING json null on empty null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING JSON NULL ON EMPTY) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING json null on empty null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING JSON EMPTY ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING json EMPTY on empty null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING JSON EMPTY ON EMPTY) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING json EMPTY on empty null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING JSON EMPTY ARRAY ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING json EMPTY ARRAY on empty null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING JSON EMPTY ARRAY ON EMPTY) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING json EMPTY ARRAY on empty null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING JSON EMPTY OBJECT ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING json EMPTY OBJECT on empty null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING JSON EMPTY OBJECT ON EMPTY) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING json EMPTY OBJECT on empty null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING JSON ERROR ON MISMATCH)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING json error on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING JSON ERROR ON MISMATCH) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING json error on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING JSON NULL ON MISMATCH)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING json null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING JSON NULL ON MISMATCH) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING json null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING JSON ERROR ON ERROR NULL ON EMPTY)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING json error on error null on empty null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING JSON ERROR ON ERROR NULL ON EMPTY) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING json error on error null on empty null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING JSON NULL ON EMPTY ERROR ON ERROR)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING json error on error null on empty null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING JSON NULL ON EMPTY ERROR ON ERROR) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING json error on error null on empty null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING JSON ERROR ON ERROR ERROR ON MISMATCH)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING json error on error error on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING JSON ERROR ON ERROR ERROR ON MISMATCH) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING json error on error error on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING JSON NULL ON EMPTY ERROR ON MISMATCH)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING json null on empty error on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING JSON NULL ON EMPTY ERROR ON MISMATCH) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING json null on empty error on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING JSON NULL ON EMPTY ERROR ON ERROR ERROR ON MISMATCH)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING json error on error null on empty error on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING JSON NULL ON EMPTY ERROR ON ERROR ERROR ON MISMATCH) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING json error on error null on empty error on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING JSON TRUNCATE ALLOW SCALARS PRETTY ASCII WITH CONDITIONAL ARRAY WRAPPER ASIS EMPTY OBJECT ON ERROR EMPTY ARRAY ON EMPTY ERROR ON MISMATCH MULTIVALUE)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING json TRUNCATE ALLOW SCALARS PRETTY ASCII WITH CONDITIONAL ARRAY WRAPPER asis EMPTY OBJECT on error EMPTY ARRAY on empty error on mismatch multivalue)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING JSON TRUNCATE ALLOW SCALARS PRETTY ASCII WITH CONDITIONAL ARRAY WRAPPER ASIS EMPTY OBJECT ON ERROR EMPTY ARRAY ON EMPTY ERROR ON MISMATCH MULTIVALUE) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING json TRUNCATE ALLOW SCALARS PRETTY ASCII WITH CONDITIONAL ARRAY WRAPPER asis EMPTY OBJECT on error EMPTY ARRAY on empty error on mismatch multivalue) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' ASIS)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING char(512) charset utf8mb4 asis null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' ASIS) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING char(512) charset utf8mb4 asis null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' ASIS ERROR ON MISMATCH)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING char(512) charset utf8mb4 asis error on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' ASIS ERROR ON MISMATCH) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING char(512) charset utf8mb4 asis error on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' ASIS MULTIVALUE)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING char(512) charset utf8mb4 asis null on mismatch multivalue)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' ASIS MULTIVALUE) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING char(512) charset utf8mb4 asis null on mismatch multivalue) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING CHAR(10) ASIS)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING char(10) charset utf8mb4 asis null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING CHAR(10) ASIS) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING char(10) charset utf8mb4 asis null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING CHAR(10) CHARSET ascii ASIS)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING char(10) charset ascii asis null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING CHAR(10) CHARSET ascii ASIS) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING char(10) charset ascii asis null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING CHAR(10) BINARY ASIS)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING char(10) charset utf8mb4 asis null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING CHAR(10) BINARY ASIS) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING char(10) charset utf8mb4 asis null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING CHAR(0) ASIS)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING char(0) charset utf8mb4 asis null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING CHAR(0) ASIS) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING char(0) charset utf8mb4 asis null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING NCHAR(3) ASIS)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING char(3) charset utf8mb4 asis null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING NCHAR(3) ASIS) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING char(3) charset utf8mb4 asis null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING BINARY(4) ASIS)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING char(4) charset binary asis null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING BINARY(4) ASIS) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING char(4) charset binary asis null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING SIGNED ASIS)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING signed asis null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING SIGNED ASIS) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING signed asis null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING SIGNED INTEGER ASIS)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING signed asis null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING SIGNED INTEGER ASIS) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING signed asis null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING UNSIGNED INTEGER ASIS)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING unsigned asis null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING UNSIGNED INTEGER ASIS) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING unsigned asis null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING DECIMAL ASIS)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING number(10,0) asis null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING DECIMAL ASIS) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING number(10,0) asis null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING DECIMAL(5) ASIS)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING number(5,0) asis null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING DECIMAL(5) ASIS) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING number(5,0) asis null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING DECIMAL(5,2) ASIS)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING number(5,2) asis null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING DECIMAL(5,2) ASIS) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING number(5,2) asis null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING NUMBER(5,2) ASIS)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING number(5,2) asis null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING NUMBER(5,2) ASIS) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING number(5,2) asis null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING FLOAT ASIS)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING float asis null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING FLOAT ASIS) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING float asis null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING FLOAT(5) ASIS)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING float asis null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING FLOAT(5) ASIS) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING float asis null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING DOUBLE ASIS)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING double asis null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING DOUBLE ASIS) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING double asis null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING DATE ASIS)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING date asis null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING DATE ASIS) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING date asis null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING DATETIME ASIS)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING datetime(0) asis null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING DATETIME ASIS) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING datetime(0) asis null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING DATETIME(3) ASIS)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING datetime(3) asis null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING DATETIME(3) ASIS) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING datetime(3) asis null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING TIME ASIS)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING time(0) asis null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING TIME ASIS) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING time(0) asis null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING TIME(2) ASIS)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING time(2) asis null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING TIME(2) ASIS) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING time(2) asis null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP DATABASE IF EXISTS zz_json_query
< ok 4 0
> exec
| CREATE DATABASE zz_json_query
< ok 1 0
> exec
| USE zz_json_query
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' RETURNING UNSIGNED WITHOUT ARRAY WRAPPER ASIS ERROR ON ERROR NULL ON EMPTY NULL ON MISMATCH)))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query.g
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` json DEFAULT NULL,
|   `b` varchar(600) GENERATED ALWAYS AS (json_query(`a`,'$.b' RETURNING unsigned WITHOUT ARRAY WRAPPER asis error on error null on empty null on mismatch)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE VIEW v AS SELECT JSON_QUERY(a, '$.b' RETURNING UNSIGNED WITHOUT ARRAY WRAPPER ASIS ERROR ON ERROR NULL ON EMPTY NULL ON MISMATCH) AS x FROM g
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
| CREATE VIEW `v` AS select json_query(`zz_json_query`.`g`.`a`,'$.b' RETURNING unsigned WITHOUT ARRAY WRAPPER asis error on error null on empty null on mismatch) AS `x` from `zz_json_query`.`g`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
