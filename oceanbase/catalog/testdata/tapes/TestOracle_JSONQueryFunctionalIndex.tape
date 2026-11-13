-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_JSONQueryFunctionalIndex default#1
> exec
| DROP DATABASE IF EXISTS zz_json_query_index
< ok 0 0
> exec
| CREATE DATABASE zz_json_query_index
< ok 1 0
> exec
| USE zz_json_query_index
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a JSON, KEY k ((JSON_QUERY(a, '$.x' ))))
< error
| 5443 42000 invalid data type for return value
> exec
| DROP DATABASE IF EXISTS zz_json_query_index
< ok 0 0
> exec
| CREATE DATABASE zz_json_query_index
< ok 1 0
> exec
| USE zz_json_query_index
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a JSON, KEY k ((JSON_QUERY(a, '$.x' RETURNING JSON))))
< error
| 3753 42000 Cannot create a functional index on a function that returns a JSON or GEOMETRY value.
> exec
| DROP DATABASE IF EXISTS zz_json_query_index
< ok 0 0
> exec
| CREATE DATABASE zz_json_query_index
< ok 1 0
> exec
| USE zz_json_query_index
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a JSON, KEY k ((JSON_QUERY(a, '$.x' RETURNING UNSIGNED))))
< error
| 5443 42000 invalid data type for return value
> exec
| DROP DATABASE IF EXISTS zz_json_query_index
< ok 0 0
> exec
| CREATE DATABASE zz_json_query_index
< ok 1 0
> exec
| USE zz_json_query_index
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a JSON, KEY k ((JSON_QUERY(a, '$.x' RETURNING JSON MULTIVALUE))))
< error
| 3753 42000 Cannot create a functional index on a function that returns a JSON or GEOMETRY value.
> exec
| DROP DATABASE IF EXISTS zz_json_query_index
< ok 0 0
> exec
| CREATE DATABASE zz_json_query_index
< ok 1 0
> exec
| USE zz_json_query_index
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a JSON, KEY k ((JSON_QUERY(a, '$.x' RETURNING UNSIGNED MULTIVALUE))))
< error
| 5443 42000 invalid data type for return value
> exec
| DROP DATABASE IF EXISTS zz_json_query_index
< ok 0 0
> exec
| CREATE DATABASE zz_json_query_index
< ok 1 0
> exec
| USE zz_json_query_index
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a JSON, KEY k ((JSON_QUERY(a, '$.x' ASIS))))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query_index.f
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| f
: bytes
| CREATE TABLE `f` (
|   `id` int(11) NOT NULL,
|   `a` json DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   INDEX `k` ((json_query(`a`,'$.x' RETURNING char(512) charset utf8mb4 asis null on mismatch))) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS zz_json_query_index
< ok 6 0
> exec
| CREATE DATABASE zz_json_query_index
< ok 1 0
> exec
| USE zz_json_query_index
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a JSON, KEY k ((JSON_QUERY(a, '$.x' RETURNING UNSIGNED ASIS))))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query_index.f
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| f
: bytes
| CREATE TABLE `f` (
|   `id` int(11) NOT NULL,
|   `a` json DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   INDEX `k` ((json_query(`a`,'$.x' RETURNING unsigned asis null on mismatch))) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS zz_json_query_index
< ok 6 0
> exec
| CREATE DATABASE zz_json_query_index
< ok 1 0
> exec
| USE zz_json_query_index
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a JSON, KEY k ((JSON_QUERY(a, '$.x' RETURNING CHAR(40) ASIS))))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query_index.f
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| f
: bytes
| CREATE TABLE `f` (
|   `id` int(11) NOT NULL,
|   `a` json DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   INDEX `k` ((json_query(`a`,'$.x' RETURNING char(40) charset utf8mb4 asis null on mismatch))) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS zz_json_query_index
< ok 6 0
> exec
| CREATE DATABASE zz_json_query_index
< ok 1 0
> exec
| USE zz_json_query_index
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a JSON, KEY k ((JSON_QUERY(a, '$.x' RETURNING JSON ASIS))))
< error
| 3753 42000 Cannot create a functional index on a function that returns a JSON or GEOMETRY value.
> exec
| DROP DATABASE IF EXISTS zz_json_query_index
< ok 0 0
> exec
| CREATE DATABASE zz_json_query_index
< ok 1 0
> exec
| USE zz_json_query_index
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a JSON, KEY k ((JSON_QUERY(a, '$.x' ASIS MULTIVALUE))))
< error
| 3753 42000 Cannot create a functional index on a function that returns a JSON or GEOMETRY value.
> exec
| DROP DATABASE IF EXISTS zz_json_query_index
< ok 0 0
> exec
| CREATE DATABASE zz_json_query_index
< ok 1 0
> exec
| USE zz_json_query_index
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a JSON, KEY k ((JSON_QUERY(a, '$.x' RETURNING UNSIGNED ASIS MULTIVALUE))))
< error
| 3753 42000 Cannot create a functional index on a function that returns a JSON or GEOMETRY value.
> exec
| DROP DATABASE IF EXISTS zz_json_query_index
< ok 0 0
> exec
| CREATE DATABASE zz_json_query_index
< ok 1 0
> exec
| USE zz_json_query_index
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a JSON, KEY k ((JSON_QUERY(a, '$.x' RETURNING UNSIGNED WITHOUT ARRAY WRAPPER ASIS ERROR ON ERROR NULL ON EMPTY NULL ON MISMATCH))))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query_index.f
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| f
: bytes
| CREATE TABLE `f` (
|   `id` int(11) NOT NULL,
|   `a` json DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   INDEX `k` ((json_query(`a`,'$.x' RETURNING unsigned WITHOUT ARRAY WRAPPER asis error on error null on empty null on mismatch))) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS zz_json_query_index
< ok 6 0
> exec
| CREATE DATABASE zz_json_query_index
< ok 1 0
> exec
| USE zz_json_query_index
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a JSON, KEY k ((JSON_QUERY(a, '$.x' RETURNING DATETIME(3) WITHOUT ARRAY WRAPPER ASIS ERROR ON ERROR NULL ON EMPTY NULL ON MISMATCH))))
< ok 0 0
> query
| SHOW CREATE TABLE zz_json_query_index.f
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| f
: bytes
| CREATE TABLE `f` (
|   `id` int(11) NOT NULL,
|   `a` json DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   INDEX `k` ((json_query(`a`,'$.x' RETURNING datetime(3) WITHOUT ARRAY WRAPPER asis error on error null on empty null on mismatch))) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS zz_json_query_index
< ok 6 0
> exec
| CREATE DATABASE zz_json_query_index
< ok 1 0
> exec
| USE zz_json_query_index
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a JSON, KEY k ((JSON_QUERY(a, '$.x' RETURNING UNSIGNED WITHOUT ARRAY WRAPPER ASIS ERROR ON ERROR NULL ON EMPTY NULL ON MISMATCH MULTIVALUE))))
< error
| 3753 42000 Cannot create a functional index on a function that returns a JSON or GEOMETRY value.
