-- pin: lctn1 mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46
-- pin: mysql84 mysql:8.4.13 sha256:a81d043930e51745522bd5c430961cf8286f0a23f0a9a2a6e0108f7ea71b0902 mysql-8.4.13

== TestCountAffectedRowsAndExplain/lctn1 lctn1#1
> begin
< ok
> exec
| /*app=bytebase*/ CREATE TABLE sitems (id INT PRIMARY KEY, note VARCHAR(64));
< ok 0 0
> exec
| /*app=bytebase*/ 
| INSERT INTO sitems VALUES (1, 'a'), (2, 'b'), (3, 'c'), (4, 'd');
< ok 4 0
> commit
< ok
> exec
| SET SESSION explain_json_format_version = 1
< error
| 1193 HY000 Unknown system variable 'explain_json_format_version'
> query
| EXPLAIN FORMAT=JSON SELECT 1 FROM sitems WHERE id <= 2
< result
< column "EXPLAIN" "VARCHAR" "string"
< row
: bytes
| {
|   "query_block": {
|     "select_id": 1,
|     "cost_info": {
|       "query_cost": "0.65"
|     },
|     "table": {
|       "table_name": "sitems",
|       "access_type": "index",
|       "possible_keys": [
|         "PRIMARY"
|       ],
|       "key": "PRIMARY",
|       "used_key_parts": [
|         "id"
|       ],
|       "key_length": "4",
|       "rows_examined_per_scan": 4,
|       "rows_produced_per_join": 2,
|       "filtered": "50.00",
|       "using_index": true,
|       "cost_info": {
|         "read_cost": "0.45",
|         "eval_cost": "0.20",
|         "prefix_cost": "0.65",
|         "data_read_per_join": "528"
|       },
|       "used_columns": [
|         "id"
|       ],
|       "attached_condition": "(`omni_qe_stats`.`sitems`.`id` <= 2)"
|     }
|   }
| }
> exec
| SET SESSION explain_json_format_version = 1
< error
| 1193 HY000 Unknown system variable 'explain_json_format_version'
> query
| EXPLAIN FORMAT=JSON SELECT 1 FROM sitems WHERE id = 2
< result
< column "EXPLAIN" "VARCHAR" "string"
< row
: bytes
| {
|   "query_block": {
|     "select_id": 1,
|     "cost_info": {
|       "query_cost": "1.00"
|     },
|     "table": {
|       "table_name": "sitems",
|       "access_type": "const",
|       "possible_keys": [
|         "PRIMARY"
|       ],
|       "key": "PRIMARY",
|       "used_key_parts": [
|         "id"
|       ],
|       "key_length": "4",
|       "ref": [
|         "const"
|       ],
|       "rows_examined_per_scan": 1,
|       "rows_produced_per_join": 1,
|       "filtered": "100.00",
|       "using_index": true,
|       "cost_info": {
|         "read_cost": "0.00",
|         "eval_cost": "0.10",
|         "prefix_cost": "0.00",
|         "data_read_per_join": "264"
|       },
|       "used_columns": [
|         "id"
|       ]
|     }
|   }
| }
> exec
| SET SESSION explain_json_format_version = 1
< error
| 1193 HY000 Unknown system variable 'explain_json_format_version'
> query
| EXPLAIN FORMAT=JSON SELECT * FROM sitems
< result
< column "EXPLAIN" "VARCHAR" "string"
< row
: bytes
| {
|   "query_block": {
|     "select_id": 1,
|     "cost_info": {
|       "query_cost": "0.65"
|     },
|     "table": {
|       "table_name": "sitems",
|       "access_type": "ALL",
|       "rows_examined_per_scan": 4,
|       "rows_produced_per_join": 4,
|       "filtered": "100.00",
|       "cost_info": {
|         "read_cost": "0.25",
|         "eval_cost": "0.40",
|         "prefix_cost": "0.65",
|         "data_read_per_join": "1K"
|       },
|       "used_columns": [
|         "id",
|         "note"
|       ]
|     }
|   }
| }
> begin
< ok
> query
| EXPLAIN UPDATE sitems SET note = 'x' WHERE id = 1
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
| UPDATE
: bytes
| sitems
: null
: bytes
| range
: bytes
| PRIMARY
: bytes
| PRIMARY
: bytes
| 4
: bytes
| const
: uint64 1
: float64 100
: bytes
| Using where
> rollback
< ok

== TestCountAffectedRowsAndExplain/mysql84 mysql84#1
> begin
< ok
> exec
| /*app=bytebase*/ CREATE TABLE sitems (id INT PRIMARY KEY, note VARCHAR(64));
< ok 0 0
> exec
| /*app=bytebase*/ 
| INSERT INTO sitems VALUES (1, 'a'), (2, 'b'), (3, 'c'), (4, 'd');
< ok 4 0
> commit
< ok
> exec
| SET SESSION explain_json_format_version = 1
< ok 0 0
> query
| EXPLAIN FORMAT=JSON SELECT 1 FROM sitems WHERE id <= 2
< result
< column "EXPLAIN" "VARCHAR" "string"
< row
: bytes
| {
|   "query_block": {
|     "select_id": 1,
|     "cost_info": {
|       "query_cost": "0.65"
|     },
|     "table": {
|       "table_name": "sitems",
|       "access_type": "index",
|       "possible_keys": [
|         "PRIMARY"
|       ],
|       "key": "PRIMARY",
|       "used_key_parts": [
|         "id"
|       ],
|       "key_length": "4",
|       "rows_examined_per_scan": 4,
|       "rows_produced_per_join": 2,
|       "filtered": "50.00",
|       "using_index": true,
|       "cost_info": {
|         "read_cost": "0.45",
|         "eval_cost": "0.20",
|         "prefix_cost": "0.65",
|         "data_read_per_join": "528"
|       },
|       "used_columns": [
|         "id"
|       ],
|       "attached_condition": "(`omni_qe_stats`.`sitems`.`id` <= 2)"
|     }
|   }
| }
> exec
| SET SESSION explain_json_format_version = DEFAULT
< ok 0 0
> exec
| SET SESSION explain_json_format_version = 1
< ok 0 0
> query
| EXPLAIN FORMAT=JSON SELECT 1 FROM sitems WHERE id = 2
< result
< column "EXPLAIN" "VARCHAR" "string"
< row
: bytes
| {
|   "query_block": {
|     "select_id": 1,
|     "cost_info": {
|       "query_cost": "1.00"
|     },
|     "table": {
|       "table_name": "sitems",
|       "access_type": "const",
|       "possible_keys": [
|         "PRIMARY"
|       ],
|       "key": "PRIMARY",
|       "used_key_parts": [
|         "id"
|       ],
|       "key_length": "4",
|       "ref": [
|         "const"
|       ],
|       "rows_examined_per_scan": 1,
|       "rows_produced_per_join": 1,
|       "filtered": "100.00",
|       "using_index": true,
|       "cost_info": {
|         "read_cost": "0.00",
|         "eval_cost": "0.10",
|         "prefix_cost": "0.00",
|         "data_read_per_join": "264"
|       },
|       "used_columns": [
|         "id"
|       ]
|     }
|   }
| }
> exec
| SET SESSION explain_json_format_version = DEFAULT
< ok 0 0
> exec
| SET SESSION explain_json_format_version = 1
< ok 0 0
> query
| EXPLAIN FORMAT=JSON SELECT * FROM sitems
< result
< column "EXPLAIN" "VARCHAR" "string"
< row
: bytes
| {
|   "query_block": {
|     "select_id": 1,
|     "cost_info": {
|       "query_cost": "0.65"
|     },
|     "table": {
|       "table_name": "sitems",
|       "access_type": "ALL",
|       "rows_examined_per_scan": 4,
|       "rows_produced_per_join": 4,
|       "filtered": "100.00",
|       "cost_info": {
|         "read_cost": "0.25",
|         "eval_cost": "0.40",
|         "prefix_cost": "0.65",
|         "data_read_per_join": "1K"
|       },
|       "used_columns": [
|         "id",
|         "note"
|       ]
|     }
|   }
| }
> exec
| SET SESSION explain_json_format_version = DEFAULT
< ok 0 0
> begin
< ok
> query
| EXPLAIN UPDATE sitems SET note = 'x' WHERE id = 1
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
| UPDATE
: bytes
| sitems
: null
: bytes
| range
: bytes
| PRIMARY
: bytes
| PRIMARY
: bytes
| 4
: bytes
| const
: uint64 1
: float64 100
: bytes
| Using where
> rollback
< ok
