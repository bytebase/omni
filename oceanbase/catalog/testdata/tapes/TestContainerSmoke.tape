-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestContainerSmoke default#1
> exec
| CREATE TABLE t1 (id INT PRIMARY KEY, name VARCHAR(100) NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE t1
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t1
: bytes
| CREATE TABLE `t1` (
|   `id` int(11) NOT NULL,
|   `name` varchar(100) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
|
| 		SELECT COLUMN_NAME, ORDINAL_POSITION, DATA_TYPE, COLUMN_TYPE,
| 		       IS_NULLABLE, COLUMN_DEFAULT, COLUMN_KEY, EXTRA,
| 		       CHARACTER_SET_NAME, COLLATION_NAME,
| 		       CHARACTER_MAXIMUM_LENGTH, NUMERIC_PRECISION, NUMERIC_SCALE
| 		FROM INFORMATION_SCHEMA.COLUMNS
| 		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
| 		ORDER BY ORDINAL_POSITION
: string
| test
: string
| t1
< result
< column "COLUMN_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "DATA_TYPE" "TEXT"
< column "COLUMN_TYPE" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "COLUMN_DEFAULT" "TEXT"
< column "COLUMN_KEY" "VARCHAR"
< column "EXTRA" "VARCHAR"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< column "CHARACTER_MAXIMUM_LENGTH" "UNSIGNED BIGINT"
< column "NUMERIC_PRECISION" "UNSIGNED BIGINT"
< column "NUMERIC_SCALE" "UNSIGNED BIGINT"
< row
: bytes
| id
: uint64 1
: bytes
| int
: bytes
| int(11)
: bytes
| NO
: null
: bytes
| PRI
: bytes
|
: null
: null
: null
: uint64 11
: uint64 0
< row
: bytes
| name
: uint64 2
: bytes
| varchar
: bytes
| varchar(100)
: bytes
| NO
: null
: bytes
|
: bytes
|
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: uint64 100
: null
: null
