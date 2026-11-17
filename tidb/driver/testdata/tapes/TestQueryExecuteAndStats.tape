-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestQueryExecuteAndStats default#1
> exec
| USE `omni_qe_test`
< ok 0 0
> begin
< ok
> exec
| CREATE TABLE items (id INT PRIMARY KEY, note VARCHAR(64));
< ok 0 0
> exec
|
| INSERT INTO items VALUES (1, 'a'), (2, 'b'), (3, 'c');
< ok 3 0
> exec
|
| UPDATE items SET note = 'z' WHERE id <= 2;
< ok 2 0
> commit
< ok
> begin
< ok
> exec
| INSERT INTO items VALUES (10, 'rb');
< ok 1 0
> exec
|
| INSERT INTO items VALUES (1, 'dup');
< error
| 1062 23000 Duplicate entry '1' for key 'items.PRIMARY'
> rollback
< ok
> exec
| INSERT INTO items VALUES (10, 'kept');
< ok 1 0
> exec
|
| INSERT INTO items VALUES (1, 'dup');
< error
| 1062 23000 Duplicate entry '1' for key 'items.PRIMARY'
> query
| EXPLAIN DELETE FROM items WHERE id <= 2;
< result
< column "id" "CHAR"
< column "estRows" "CHAR"
< column "task" "CHAR"
< column "access object" "CHAR"
< column "operator info" "CHAR"
< row
: bytes
| Delete_3
: bytes
| N/A
: bytes
| root
: bytes
|
: bytes
| N/A
< row
: bytes
| └─TableReader_7
: bytes
| 3333.33
: bytes
| root
: bytes
|
: bytes
| data:TableRangeScan_6
< row
: bytes
|   └─TableRangeScan_6
: bytes
| 3333.33
: bytes
| cop[tikv]
: bytes
| table:items
: bytes
| range:[-inf,2], keep order:false, stats:pseudo
> begin
< ok
> query
| EXPLAIN UPDATE items SET note = 'x' WHERE id = 1
< result
< column "id" "CHAR"
< column "estRows" "CHAR"
< column "task" "CHAR"
< column "access object" "CHAR"
< column "operator info" "CHAR"
< row
: bytes
| Update_3
: bytes
| N/A
: bytes
| root
: bytes
|
: bytes
| N/A
< row
: bytes
| └─Point_Get_1
: bytes
| 1.00
: bytes
| root
: bytes
| table:items
: bytes
| handle:1, lock
> rollback
< ok
> query
| BATCH ON id LIMIT 2 DRY RUN DELETE FROM items WHERE id > 0;
< result
< column "split statement examples" "CHAR"
< row
: bytes
| DELETE FROM `omni_qe_test`.`items` WHERE (`id` BETWEEN 1 AND 2 AND (`id` > 0))
< row
: bytes
| DELETE FROM `omni_qe_test`.`items` WHERE (`id` BETWEEN 3 AND 10 AND (`id` > 0))
> query
| SELECT
| 		42 AS i,
| 		3.5 AS f,
| 		'hello' AS s,
| 		CAST('2026-03-04 05:06:07.123456' AS DATETIME(6)) AS dt,
| 		NULL AS nul;
< result
< column "i" "BIGINT" "int64"
< column "f" "DECIMAL" "string" precision 4 scale 1
< column "s" "VARCHAR" "string"
< column "dt" "DATETIME" precision 6 scale 6
< column "nul" "NULL" "*interface {}"
< row
: int64 42
: bytes
| 3.5
: bytes
| hello
: bytes
| 2026-03-04 05:06:07.123456
: null
> query
| SELECT `id` FROM `items` ORDER BY `id` LIMIT 2
< result
< column "id" "INT" "int32"
< row
: int64 1
< row
: int64 2
> query
| EXPLAIN FORMAT=tidb_json SELECT note FROM items WHERE id > 1
< result
< column "TiDB_JSON" "CHAR"
< row
: bytes
| [
|     {
|         "id": "TableReader_9",
|         "estRows": "3333.33",
|         "taskType": "root",
|         "operatorInfo": "data:Projection_5",
|         "subOperators": [
|             {
|                 "id": "Projection_5",
|                 "estRows": "3333.33",
|                 "taskType": "cop[tikv]",
|                 "operatorInfo": "omni_qe_test.items.note",
|                 "subOperators": [
|                     {
|                         "id": "TableRangeScan_8",
|                         "estRows": "3333.33",
|                         "taskType": "cop[tikv]",
|                         "accessObject": "table:items",
|                         "operatorInfo": "range:(1,+inf], keep order:false, stats:pseudo"
|                     }
|                 ]
|             }
|         ]
|     }
| ]
|
> query
| SELECT COUNT(*) FROM items WHERE id = 99;
< result
< column "COUNT(*)" "BIGINT" "int64"
< row
: int64 0
> query
| SELECT 1;
< result
< column "1" "BIGINT" "int64"
< row
: int64 1
> query
|
| SELECT * FROM missing_tbl;
< error
| 1146 42S02 Table 'omni_qe_test.missing_tbl' doesn't exist
> query
| SELECT id, repeat('x', 80) FROM items ORDER BY id;
< result
< column "id" "INT" "int32"
< column "repeat('x', 80)" "TEXT" "string"
< row
: int64 1
: bytes
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 2
: bytes
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 3
: bytes
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 10
: bytes
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
> query
|
| SELECT 7;
< result
< column "7" "BIGINT" "int64"
< row
: int64 7
