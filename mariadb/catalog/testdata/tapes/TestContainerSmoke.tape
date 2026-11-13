-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

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
|   `id` int(11) NOT NULL,
|   `name` varchar(100) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
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
| catalog_shared
: string
| t1
< result
< column "COLUMN_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "DATA_TYPE" "VARCHAR"
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
: uint64 10
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
| utf8mb4_uca1400_ai_ci
: uint64 100
: null
: null
