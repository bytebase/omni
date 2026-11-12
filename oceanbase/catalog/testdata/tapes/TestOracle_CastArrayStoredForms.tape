-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_CastArrayStoredForms default#1
> exec
| DROP DATABASE IF EXISTS zz_cast_array
< ok 0 0
> exec
| CREATE DATABASE zz_cast_array
< ok 1 0
> exec
| USE zz_cast_array
< ok 0 0
> exec
| CREATE TABLE zz (a JSON)
< ok 0 0
> query
| SHOW CREATE TABLE zz_cast_array.zz
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| zz
: bytes
| CREATE TABLE `zz` (
|   `a` json DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE OR REPLACE VIEW v AS SELECT CAST(a AS ARRAY(ARRAY(INT))) x FROM zz
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
| CREATE VIEW `v` AS select cast(`zz_cast_array`.`zz`.`a` as ARRAY(ARRAY(INT))) AS `x` from `zz_cast_array`.`zz`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT CAST(a AS ARRAY(ARRAY<INT>)) x FROM zz
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
| CREATE VIEW `v` AS select cast(`zz_cast_array`.`zz`.`a` as ARRAY(ARRAY(INT))) AS `x` from `zz_cast_array`.`zz`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT CAST(a AS ARRAY(BIGINT UNSIGNED)) x FROM zz
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
| CREATE VIEW `v` AS select cast(`zz_cast_array`.`zz`.`a` as ARRAY(BIGINT UNSIGNED)) AS `x` from `zz_cast_array`.`zz`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT CAST(a AS ARRAY(BIGINT(5))) x FROM zz
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
| CREATE VIEW `v` AS select cast(`zz_cast_array`.`zz`.`a` as ARRAY(BIGINT)) AS `x` from `zz_cast_array`.`zz`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT CAST(a AS ARRAY(BIGINT)) x FROM zz
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
| CREATE VIEW `v` AS select cast(`zz_cast_array`.`zz`.`a` as ARRAY(BIGINT)) AS `x` from `zz_cast_array`.`zz`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT CAST(a AS ARRAY(BOOL)) x FROM zz
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
| CREATE VIEW `v` AS select cast(`zz_cast_array`.`zz`.`a` as ARRAY(TINYINT)) AS `x` from `zz_cast_array`.`zz`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT CAST(a AS ARRAY(DOUBLE(5,2))) x FROM zz
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
| CREATE VIEW `v` AS select cast(`zz_cast_array`.`zz`.`a` as ARRAY(DOUBLE)) AS `x` from `zz_cast_array`.`zz`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT CAST(a AS ARRAY(FLOAT(30))) x FROM zz
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
| CREATE VIEW `v` AS select cast(`zz_cast_array`.`zz`.`a` as ARRAY(FLOAT)) AS `x` from `zz_cast_array`.`zz`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT CAST(a AS ARRAY(FLOAT(5,2))) x FROM zz
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
| CREATE VIEW `v` AS select cast(`zz_cast_array`.`zz`.`a` as ARRAY(FLOAT)) AS `x` from `zz_cast_array`.`zz`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT CAST(a AS ARRAY(FLOAT)) x FROM zz
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
| CREATE VIEW `v` AS select cast(`zz_cast_array`.`zz`.`a` as ARRAY(FLOAT)) AS `x` from `zz_cast_array`.`zz`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT CAST(a AS ARRAY(INT ZEROFILL)) x FROM zz
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
| CREATE VIEW `v` AS select cast(`zz_cast_array`.`zz`.`a` as ARRAY(INT UNSIGNED)) AS `x` from `zz_cast_array`.`zz`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT CAST(a AS ARRAY(INT)) x FROM zz
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
| CREATE VIEW `v` AS select cast(`zz_cast_array`.`zz`.`a` as ARRAY(INT)) AS `x` from `zz_cast_array`.`zz`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT CAST(a AS ARRAY(INTEGER)) x FROM zz
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
| CREATE VIEW `v` AS select cast(`zz_cast_array`.`zz`.`a` as ARRAY(INT)) AS `x` from `zz_cast_array`.`zz`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT CAST(a AS ARRAY(NVARCHAR(3))) x FROM zz
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
| CREATE VIEW `v` AS select cast(`zz_cast_array`.`zz`.`a` as ARRAY(VARCHAR(3))) AS `x` from `zz_cast_array`.`zz`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT CAST(a AS ARRAY(REAL)) x FROM zz
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
| CREATE VIEW `v` AS select cast(`zz_cast_array`.`zz`.`a` as ARRAY(DOUBLE)) AS `x` from `zz_cast_array`.`zz`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT CAST(a AS ARRAY(SMALLINT)) x FROM zz
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
| CREATE VIEW `v` AS select cast(`zz_cast_array`.`zz`.`a` as ARRAY(SMALLINT)) AS `x` from `zz_cast_array`.`zz`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT CAST(a AS ARRAY(TINYINT(1))) x FROM zz
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
| CREATE VIEW `v` AS select cast(`zz_cast_array`.`zz`.`a` as ARRAY(TINYINT)) AS `x` from `zz_cast_array`.`zz`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT CAST(a AS ARRAY(VARCHAR(0))) x FROM zz
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
| CREATE VIEW `v` AS select cast(`zz_cast_array`.`zz`.`a` as ARRAY(VARCHAR(0))) AS `x` from `zz_cast_array`.`zz`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT CAST(a AS ARRAY(VARCHAR(3) CHARSET ascii)) x FROM zz
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
| CREATE VIEW `v` AS select cast(`zz_cast_array`.`zz`.`a` as ARRAY(VARCHAR(3))) AS `x` from `zz_cast_array`.`zz`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT CAST(a AS ARRAY(VARCHAR(3))) x FROM zz
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
| CREATE VIEW `v` AS select cast(`zz_cast_array`.`zz`.`a` as ARRAY(VARCHAR(3))) AS `x` from `zz_cast_array`.`zz`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT CAST(a AS ARRAY(VECTOR(3))) x FROM zz
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
| CREATE VIEW `v` AS select cast(`zz_cast_array`.`zz`.`a` as ARRAY(VECTOR(3))) AS `x` from `zz_cast_array`.`zz`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT CAST(a AS ARRAY<ARRAY(INT)>) x FROM zz
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
| CREATE VIEW `v` AS select cast(`zz_cast_array`.`zz`.`a` as ARRAY(ARRAY(INT))) AS `x` from `zz_cast_array`.`zz`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT CAST(a AS ARRAY<ARRAY<ARRAY<INT>>>) x FROM zz
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
| CREATE VIEW `v` AS select cast(`zz_cast_array`.`zz`.`a` as ARRAY(ARRAY(ARRAY(INT)))) AS `x` from `zz_cast_array`.`zz`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT CAST(a AS ARRAY<ARRAY<INT>>) x FROM zz
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
| CREATE VIEW `v` AS select cast(`zz_cast_array`.`zz`.`a` as ARRAY(ARRAY(INT))) AS `x` from `zz_cast_array`.`zz`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT CAST(a AS ARRAY<INT>) x FROM zz
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
| CREATE VIEW `v` AS select cast(`zz_cast_array`.`zz`.`a` as ARRAY(INT)) AS `x` from `zz_cast_array`.`zz`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE OR REPLACE VIEW v AS SELECT CAST(a AS ARRAY<VARCHAR(3)>) x FROM zz
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
| CREATE VIEW `v` AS select cast(`zz_cast_array`.`zz`.`a` as ARRAY(VARCHAR(3))) AS `x` from `zz_cast_array`.`zz`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
