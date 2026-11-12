-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_JSONValueAlias default#1
> exec
| DROP DATABASE IF EXISTS zz_json_value_alias
< ok 0 0
> exec
| CREATE DATABASE zz_json_value_alias
< ok 1 0
> exec
| USE zz_json_value_alias
< ok 0 0
> exec
| CREATE TABLE t (a JSON)
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value_alias.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` json DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE OR REPLACE VIEW v AS SELECT JSON_VALUE(a, '$') FROM t
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'zz_json_value_alias' AND TABLE_NAME = 'v'
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| JSON_VALUE(a, '$')
> exec
| CREATE OR REPLACE VIEW v AS SELECT JSON_VALUE(a, '$' RETURNING SIGNED) FROM t
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'zz_json_value_alias' AND TABLE_NAME = 'v'
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| JSON_VALUE(a, '$' RETURNING SIGNED)
> exec
| CREATE OR REPLACE VIEW v AS SELECT JSON_VALUE(a, '$' RETURNING CHAR(3)) FROM t
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'zz_json_value_alias' AND TABLE_NAME = 'v'
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| JSON_VALUE(a, '$' RETURNING CHAR(3))
> exec
| CREATE OR REPLACE VIEW v AS SELECT JSON_VALUE(a, '$' NULL ON EMPTY) FROM t
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'zz_json_value_alias' AND TABLE_NAME = 'v'
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| JSON_VALUE(a, '$' NULL ON EMPTY)
> exec
| CREATE OR REPLACE VIEW v AS SELECT JSON_VALUE(a, '$' DEFAULT 'x' ON EMPTY) FROM t
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'zz_json_value_alias' AND TABLE_NAME = 'v'
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| JSON_VALUE(a, '$' DEFAULT 'x' ON EMPTY)
> exec
| CREATE OR REPLACE VIEW v AS SELECT JSON_VALUE(a, '$' DEFAULT -1 ON ERROR) FROM t
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'zz_json_value_alias' AND TABLE_NAME = 'v'
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| JSON_VALUE(a, '$' DEFAULT -1 ON ERROR)
> exec
| CREATE OR REPLACE VIEW v AS SELECT JSON_VALUE(a, '$' RETURNING CHAR(3) TRUNCATE ASCII) FROM t
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'zz_json_value_alias' AND TABLE_NAME = 'v'
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| JSON_VALUE(a, '$' RETURNING CHAR(3) TRUNCATE ASCII)
