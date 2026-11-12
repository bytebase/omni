-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_JSONValueFunctionalIndex default#1
> exec
| DROP DATABASE IF EXISTS zz_json_value_index
< ok 0 0
> exec
| CREATE DATABASE zz_json_value_index
< ok 1 0
> exec
| USE zz_json_value_index
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' ))))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value_index.f
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| f
: bytes
| CREATE TABLE `f` (
|   `a` json DEFAULT NULL,
|   KEY `k` ((json_value(`a`,'$.x' RETURNING char(512) charset utf8mb4))) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS zz_json_value_index
< ok 4 0
> exec
| CREATE DATABASE zz_json_value_index
< ok 1 0
> exec
| USE zz_json_value_index
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING SIGNED))))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value_index.f
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| f
: bytes
| CREATE TABLE `f` (
|   `a` json DEFAULT NULL,
|   KEY `k` ((json_value(`a`,'$.x' RETURNING signed))) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS zz_json_value_index
< ok 4 0
> exec
| CREATE DATABASE zz_json_value_index
< ok 1 0
> exec
| USE zz_json_value_index
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING CHAR(10)))))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value_index.f
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| f
: bytes
| CREATE TABLE `f` (
|   `a` json DEFAULT NULL,
|   KEY `k` ((json_value(`a`,'$.x' RETURNING char(10) charset utf8mb4))) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS zz_json_value_index
< ok 4 0
> exec
| CREATE DATABASE zz_json_value_index
< ok 1 0
> exec
| USE zz_json_value_index
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING DECIMAL(5,2)))))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value_index.f
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| f
: bytes
| CREATE TABLE `f` (
|   `a` json DEFAULT NULL,
|   KEY `k` ((json_value(`a`,'$.x' RETURNING number(5,2)))) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS zz_json_value_index
< ok 4 0
> exec
| CREATE DATABASE zz_json_value_index
< ok 1 0
> exec
| USE zz_json_value_index
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING DATETIME(3)))))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value_index.f
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| f
: bytes
| CREATE TABLE `f` (
|   `a` json DEFAULT NULL,
|   KEY `k` ((json_value(`a`,'$.x' RETURNING datetime(3)))) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS zz_json_value_index
< ok 4 0
> exec
| CREATE DATABASE zz_json_value_index
< ok 1 0
> exec
| USE zz_json_value_index
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING DOUBLE))))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value_index.f
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| f
: bytes
| CREATE TABLE `f` (
|   `a` json DEFAULT NULL,
|   KEY `k` ((json_value(`a`,'$.x' RETURNING double))) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS zz_json_value_index
< ok 4 0
> exec
| CREATE DATABASE zz_json_value_index
< ok 1 0
> exec
| USE zz_json_value_index
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING JSON))))
< error
| 3753 42000 Cannot create a functional index on a function that returns a JSON or GEOMETRY value.
> exec
| DROP DATABASE IF EXISTS zz_json_value_index
< ok 0 0
> exec
| CREATE DATABASE zz_json_value_index
< ok 1 0
> exec
| USE zz_json_value_index
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING CHAR))))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value_index.f
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| f
: bytes
| CREATE TABLE `f` (
|   `a` json DEFAULT NULL,
|   KEY `k` ((json_value(`a`,'$.x' RETURNING char charset utf8mb4))) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS zz_json_value_index
< ok 4 0
> exec
| CREATE DATABASE zz_json_value_index
< ok 1 0
> exec
| USE zz_json_value_index
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING CHAR(512)))))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value_index.f
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| f
: bytes
| CREATE TABLE `f` (
|   `a` json DEFAULT NULL,
|   KEY `k` ((json_value(`a`,'$.x' RETURNING char(512) charset utf8mb4))) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS zz_json_value_index
< ok 4 0
> exec
| CREATE DATABASE zz_json_value_index
< ok 1 0
> exec
| USE zz_json_value_index
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING CHAR(513)))))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value_index.f
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| f
: bytes
| CREATE TABLE `f` (
|   `a` json DEFAULT NULL,
|   KEY `k` ((json_value(`a`,'$.x' RETURNING char(513) charset utf8mb4))) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS zz_json_value_index
< ok 4 0
> exec
| CREATE DATABASE zz_json_value_index
< ok 1 0
> exec
| USE zz_json_value_index
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING BINARY(3)))))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value_index.f
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| f
: bytes
| CREATE TABLE `f` (
|   `a` json DEFAULT NULL,
|   KEY `k` ((json_value(`a`,'$.x' RETURNING char(3) charset binary))) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS zz_json_value_index
< ok 4 0
> exec
| CREATE DATABASE zz_json_value_index
< ok 1 0
> exec
| USE zz_json_value_index
< ok 0 0
> exec
| CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' RETURNING YEAR))))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_value_index.f
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| f
: bytes
| CREATE TABLE `f` (
|   `a` json DEFAULT NULL,
|   KEY `k` ((json_value(`a`,'$.x' RETURNING year(0)))) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
