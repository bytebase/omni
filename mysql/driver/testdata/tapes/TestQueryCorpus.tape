-- pin: lctn1 mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46
-- pin: mysql84 mysql:8.4.13 sha256:a81d043930e51745522bd5c430961cf8286f0a23f0a9a2a6e0108f7ea71b0902 mysql-8.4.13

== TestQueryCorpus/lctn1 lctn1#1
> exec
| CREATE TABLE qitems (
|   id BIGINT PRIMARY KEY,
|   note VARCHAR(64),
|   amount DECIMAL(10,2),
|   created DATETIME(6),
|   flag TINYINT(1),
|   payload JSON,
|   raw VARBINARY(8),
|   bigu BIGINT UNSIGNED
| );
< ok 0 0
> exec
|
| INSERT INTO qitems VALUES
|   (1, 'a', 1.25, '2026-03-04 05:06:07.123456', 1, '{"a": 1}', 0x00ff, 18446744073709551615),
|   (2, 'b', 2.50, '2026-03-04 05:06:08', 0, '{"b": 2}', 0x10, 10),
|   (3, 'c', NULL, NULL, NULL, NULL, NULL, NULL);
< ok 3 0
> query
| /*app=bytebase*/ SELECT * FROM qitems ORDER BY id;
< result
< column "id" "BIGINT" "int64"
< column "note" "VARCHAR"
< column "amount" "DECIMAL" precision 10 scale 2
< column "created" "DATETIME" precision 6 scale 6
< column "flag" "TINYINT"
< column "payload" "JSON"
< column "raw" "VARBINARY" "[]uint8"
< column "bigu" "UNSIGNED BIGINT"
< row
: int64 1
: bytes
| a
: bytes
| 1.25
: bytes
| 2026-03-04 05:06:07.123456
: int64 1
: bytes
| {"a": 1}
: bytes
|q "\x00\xff"
: uint64 18446744073709551615
< row
: int64 2
: bytes
| b
: bytes
| 2.50
: bytes
| 2026-03-04 05:06:08.000000
: int64 0
: bytes
| {"b": 2}
: bytes
|q "\x10"
: uint64 10
< row
: int64 3
: bytes
| c
: null
: null
: null
: null
: null
: null
> query
| /*app=bytebase*/ SELECT id FROM qitems ORDER BY id LIMIT 2;
< result
< column "id" "BIGINT" "int64"
< row
: int64 1
< row
: int64 2
> query ^
< result
< column "id" "BIGINT" "int64"
< row
: int64 1
< row
: int64 2
> query
| /*app=bytebase*/ SHOW TABLES;
< result
< column "Tables_in_omni_qe_query" "VARCHAR"
< row
: bytes
| qitems
> query
| /*app=bytebase*/ EXPLAIN FORMAT=JSON SELECT note FROM qitems WHERE id > 1;
< result
< column "EXPLAIN" "VARCHAR" "string"
< row
: bytes
| {
|   "query_block": {
|     "select_id": 1,
|     "cost_info": {
|       "query_cost": "0.91"
|     },
|     "table": {
|       "table_name": "qitems",
|       "access_type": "range",
|       "possible_keys": [
|         "PRIMARY"
|       ],
|       "key": "PRIMARY",
|       "used_key_parts": [
|         "id"
|       ],
|       "key_length": "8",
|       "rows_examined_per_scan": 2,
|       "rows_produced_per_join": 2,
|       "filtered": "100.00",
|       "cost_info": {
|         "read_cost": "0.71",
|         "eval_cost": "0.20",
|         "prefix_cost": "0.91",
|         "data_read_per_join": "624"
|       },
|       "used_columns": [
|         "id",
|         "note"
|       ],
|       "attached_condition": "(`omni_qe_query`.`qitems`.`id` > 1)"
|     }
|   }
| }
> query
| /*app=bytebase*/ EXPLAIN SELECT note FROM qitems WHERE id > 1;
< result
< column "id" "UNSIGNED BIGINT" "uint64"
< column "select_type" "VARCHAR" "string"
< column "table" "VARCHAR"
< column "partitions" "MEDIUMTEXT"
< column "type" "VARCHAR"
< column "possible_keys" "VARCHAR"
< column "key" "VARCHAR"
< column "key_len" "VARCHAR"
< column "ref" "VARCHAR"
< column "rows" "UNSIGNED BIGINT"
< column "filtered" "DOUBLE" precision 9223372036854775807 scale 2
< column "Extra" "VARCHAR" "string"
< row
: uint64 1
: bytes
| SIMPLE
: bytes
| qitems
: null
: bytes
| range
: bytes
| PRIMARY
: bytes
| PRIMARY
: bytes
| 8
: null
: uint64 2
: float64 100
: bytes
| Using where
> query
| /*app=bytebase*/ SELECT 1;
< result
< column "1" "BIGINT" "int64"
< row
: int64 1
> query
| /*app=bytebase*/ 
| SELECT * FROM missing_tbl;
< error
| 1146 42S02 Table 'omni_qe_query.missing_tbl' doesn't exist
> query
| /*app=bytebase*/ SELECT COUNT(*) FROM qitems WHERE id = 99;
< result
< column "COUNT(*)" "BIGINT" "int64"
< row
: int64 0
> query
| /*app=bytebase*/ SELECT id, repeat('x', 80) FROM qitems ORDER BY id;
< result
< column "id" "BIGINT" "int64"
< column "repeat('x', 80)" "VARCHAR"
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
> query
| /*app=bytebase*/ 
| SELECT 7;
< result
< column "7" "BIGINT" "int64"
< row
: int64 7
> exec
| UPDATE qitems SET note = 'z' WHERE id <= 2;
< ok 2 0
> query
| /*app=bytebase*/ SELECT id, note FROM qitems ORDER BY id;
< result
< column "id" "BIGINT" "int64"
< column "note" "VARCHAR"
< row
: int64 1
: bytes
| z
< row
: int64 2
: bytes
| z
< row
: int64 3
: bytes
| c
> query ^
< result
< column "id" "BIGINT" "int64"
< column "note" "VARCHAR"
< row
: int64 1
: bytes
| z
< row
: int64 2
: bytes
| z
< row
: int64 3
: bytes
| c

== TestQueryCorpus/mysql84 mysql84#1
> exec
| CREATE TABLE qitems (
|   id BIGINT PRIMARY KEY,
|   note VARCHAR(64),
|   amount DECIMAL(10,2),
|   created DATETIME(6),
|   flag TINYINT(1),
|   payload JSON,
|   raw VARBINARY(8),
|   bigu BIGINT UNSIGNED
| );
< ok 0 0
> exec
|
| INSERT INTO qitems VALUES
|   (1, 'a', 1.25, '2026-03-04 05:06:07.123456', 1, '{"a": 1}', 0x00ff, 18446744073709551615),
|   (2, 'b', 2.50, '2026-03-04 05:06:08', 0, '{"b": 2}', 0x10, 10),
|   (3, 'c', NULL, NULL, NULL, NULL, NULL, NULL);
< ok 3 0
> query
| /*app=bytebase*/ SELECT * FROM qitems ORDER BY id;
< result
< column "id" "BIGINT" "int64"
< column "note" "VARCHAR"
< column "amount" "DECIMAL" precision 10 scale 2
< column "created" "DATETIME" precision 6 scale 6
< column "flag" "TINYINT"
< column "payload" "JSON"
< column "raw" "VARBINARY" "[]uint8"
< column "bigu" "UNSIGNED BIGINT"
< row
: int64 1
: bytes
| a
: bytes
| 1.25
: bytes
| 2026-03-04 05:06:07.123456
: int64 1
: bytes
| {"a": 1}
: bytes
|q "\x00\xff"
: uint64 18446744073709551615
< row
: int64 2
: bytes
| b
: bytes
| 2.50
: bytes
| 2026-03-04 05:06:08.000000
: int64 0
: bytes
| {"b": 2}
: bytes
|q "\x10"
: uint64 10
< row
: int64 3
: bytes
| c
: null
: null
: null
: null
: null
: null
> query
| /*app=bytebase*/ SELECT id FROM qitems ORDER BY id LIMIT 2;
< result
< column "id" "BIGINT" "int64"
< row
: int64 1
< row
: int64 2
> query ^
< result
< column "id" "BIGINT" "int64"
< row
: int64 1
< row
: int64 2
> query
| /*app=bytebase*/ SHOW TABLES;
< result
< column "Tables_in_omni_qe_query" "VARCHAR" "string"
< row
: bytes
| qitems
> query
| /*app=bytebase*/ EXPLAIN FORMAT=JSON SELECT note FROM qitems WHERE id > 1;
< result
< column "EXPLAIN" "VARCHAR" "string"
< row
: bytes
| {
|   "query_block": {
|     "select_id": 1,
|     "cost_info": {
|       "query_cost": "0.91"
|     },
|     "table": {
|       "table_name": "qitems",
|       "access_type": "range",
|       "possible_keys": [
|         "PRIMARY"
|       ],
|       "key": "PRIMARY",
|       "used_key_parts": [
|         "id"
|       ],
|       "key_length": "8",
|       "rows_examined_per_scan": 2,
|       "rows_produced_per_join": 2,
|       "filtered": "100.00",
|       "cost_info": {
|         "read_cost": "0.71",
|         "eval_cost": "0.20",
|         "prefix_cost": "0.91",
|         "data_read_per_join": "624"
|       },
|       "used_columns": [
|         "id",
|         "note"
|       ],
|       "attached_condition": "(`omni_qe_query`.`qitems`.`id` > 1)"
|     }
|   }
| }
> query
| /*app=bytebase*/ EXPLAIN SELECT note FROM qitems WHERE id > 1;
< result
< column "id" "UNSIGNED BIGINT" "uint64"
< column "select_type" "VARCHAR" "string"
< column "table" "VARCHAR"
< column "partitions" "MEDIUMTEXT"
< column "type" "VARCHAR"
< column "possible_keys" "VARCHAR"
< column "key" "VARCHAR"
< column "key_len" "VARCHAR"
< column "ref" "VARCHAR"
< column "rows" "UNSIGNED BIGINT"
< column "filtered" "DOUBLE" precision 9223372036854775807 scale 2
< column "Extra" "VARCHAR" "string"
< row
: uint64 1
: bytes
| SIMPLE
: bytes
| qitems
: null
: bytes
| range
: bytes
| PRIMARY
: bytes
| PRIMARY
: bytes
| 8
: null
: uint64 2
: float64 100
: bytes
| Using where
> query
| /*app=bytebase*/ SELECT 1;
< result
< column "1" "BIGINT" "int64"
< row
: int64 1
> query
| /*app=bytebase*/ 
| SELECT * FROM missing_tbl;
< error
| 1146 42S02 Table 'omni_qe_query.missing_tbl' doesn't exist
> query
| /*app=bytebase*/ SELECT COUNT(*) FROM qitems WHERE id = 99;
< result
< column "COUNT(*)" "BIGINT" "int64"
< row
: int64 0
> query
| /*app=bytebase*/ SELECT id, repeat('x', 80) FROM qitems ORDER BY id;
< result
< column "id" "BIGINT" "int64"
< column "repeat('x', 80)" "VARCHAR"
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
> query
| /*app=bytebase*/ 
| SELECT 7;
< result
< column "7" "BIGINT" "int64"
< row
: int64 7
> exec
| UPDATE qitems SET note = 'z' WHERE id <= 2;
< ok 2 0
> query
| /*app=bytebase*/ SELECT id, note FROM qitems ORDER BY id;
< result
< column "id" "BIGINT" "int64"
< column "note" "VARCHAR"
< row
: int64 1
: bytes
| z
< row
: int64 2
: bytes
| z
< row
: int64 3
: bytes
| c
> query ^
< result
< column "id" "BIGINT" "int64"
< column "note" "VARCHAR"
< row
: int64 1
: bytes
| z
< row
: int64 2
: bytes
| z
< row
: int64 3
: bytes
| c
