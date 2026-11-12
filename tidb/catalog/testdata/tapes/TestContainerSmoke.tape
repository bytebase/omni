-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestContainerSmoke default#1
> exec
| CREATE TABLE t1 (id INT PRIMARY KEY, name VARCHAR(100) NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE t1
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t1
: bytes
| CREATE TABLE `t1` (
|   `id` int NOT NULL,
|   `name` varchar(100) NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
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
< column "ORDINAL_POSITION" "BIGINT"
< column "DATA_TYPE" "VARCHAR"
< column "COLUMN_TYPE" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "COLUMN_DEFAULT" "TEXT"
< column "COLUMN_KEY" "VARCHAR"
< column "EXTRA" "VARCHAR"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< column "CHARACTER_MAXIMUM_LENGTH" "BIGINT"
< column "NUMERIC_PRECISION" "BIGINT"
< column "NUMERIC_SCALE" "BIGINT"
< row
: bytes
| id
: int64 1
: bytes
| int
: bytes
| int
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
: int64 11
: int64 0
< row
: bytes
| name
: int64 2
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
| utf8mb4_bin
: int64 100
: null
: null
