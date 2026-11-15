-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8
-- pin: mariadb107 mariadb:10.7.8 sha256:9a48ac9f196f3d4fd6fea2cab59a49df9e7ca459bf14b2f7b85a0e38a5454571 mariadb-10.7.8

== TestQueryExecuteAndStats/default default#1
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
| 1062 23000 Duplicate entry '1' for key 'PRIMARY'
> rollback
< ok
> exec
| INSERT INTO items VALUES (10, 'kept');
< ok 1 0
> exec
|
| INSERT INTO items VALUES (1, 'dup');
< error
| 1062 23000 Duplicate entry '1' for key 'PRIMARY'
> exec
| CREATE PROCEDURE sp_items_count(OUT n INT)
| BEGIN
|   SELECT COUNT(*) INTO n FROM items;
| END
< ok 0 0
> exec
| SET SESSION explain_json_format_version = 1
< error
| 1193 HY000 Unknown system variable 'explain_json_format_version'
> query
| EXPLAIN FORMAT=JSON SELECT 1 FROM items WHERE id <= 2
< result
< column "EXPLAIN" "VARCHAR"
< row
: bytes
| {
|   "query_block": {
|     "select_id": 1,
|     "cost": 0.00282078,
|     "nested_loop": [
|       {
|         "table": {
|           "table_name": "items",
|           "access_type": "range",
|           "possible_keys": ["PRIMARY"],
|           "key": "PRIMARY",
|           "key_length": "4",
|           "used_key_parts": ["id"],
|           "loops": 1,
|           "rows": 2,
|           "cost": 0.00282078,
|           "filtered": 100,
|           "attached_condition": "items.`id` <= 2"
|         }
|       }
|     ]
|   }
| }
> begin
< ok
> query
| EXPLAIN UPDATE items SET note = 'x' WHERE id = 1
< result
< column "id" "UNSIGNED BIGINT"
< column "select_type" "VARCHAR"
< column "table" "VARCHAR"
< column "type" "VARCHAR"
< column "possible_keys" "VARCHAR"
< column "key" "VARCHAR"
< column "key_len" "VARCHAR"
< column "ref" "VARCHAR"
< column "rows" "VARCHAR"
< column "Extra" "VARCHAR"
< row
: uint64 1
: bytes
| SIMPLE
: bytes
| items
: bytes
| range
: bytes
| PRIMARY
: bytes
| PRIMARY
: bytes
| 4
: null
: bytes
| 1
: bytes
| Using where
> rollback
< ok
> query
| SELECT
| 				42 AS i,
| 				3.5 AS f,
| 				'hello' AS s,
| 				CAST('2026-03-04 05:06:07.123456' AS DATETIME(6)) AS dt,
| 				CAST(0xdeadbeef AS BINARY(4)) AS bin,
| 				NULL AS nul;
< result
< column "i" "INT"
< column "f" "DECIMAL"
< column "s" "VARCHAR"
< column "dt" "DATETIME"
< column "bin" "VARBINARY"
< column "nul" "NULL"
< row
: int64 42
: bytes
| 3.5
: bytes
| hello
: bytes
| 2026-03-04 05:06:07.123456
: bytes
|q "ޭ\xbe\xef"
: null
> query
| SELECT id FROM items ORDER BY id LIMIT 2;
< result
< column "id" "INT"
< row
: int64 1
< row
: int64 2
> query
| EXPLAIN FORMAT=JSON SELECT note FROM items WHERE id > 1
< result
< column "EXPLAIN" "VARCHAR"
< row
: bytes
| {
|   "query_block": {
|     "select_id": 1,
|     "cost": 0.00301265,
|     "nested_loop": [
|       {
|         "table": {
|           "table_name": "items",
|           "access_type": "range",
|           "possible_keys": ["PRIMARY"],
|           "key": "PRIMARY",
|           "key_length": "4",
|           "used_key_parts": ["id"],
|           "loops": 1,
|           "rows": 3,
|           "cost": 0.00301265,
|           "filtered": 100,
|           "attached_condition": "items.`id` > 1"
|         }
|       }
|     ]
|   }
| }
> exec
| DELETE FROM items WHERE id = 10 RETURNING note;
< ok 0 0
> query
| SELECT COUNT(*) FROM items WHERE id = 99;
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 0
> query
| SELECT 1;
< result
< column "1" "INT"
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
< column "id" "INT"
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
|
| SELECT 7;
< result
< column "7" "INT"
< row
: int64 7

== TestQueryExecuteAndStats/mariadb107 mariadb107#1
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
| 1062 23000 Duplicate entry '1' for key 'PRIMARY'
> rollback
< ok
> exec
| INSERT INTO items VALUES (10, 'kept');
< ok 1 0
> exec
|
| INSERT INTO items VALUES (1, 'dup');
< error
| 1062 23000 Duplicate entry '1' for key 'PRIMARY'
> exec
| CREATE PROCEDURE sp_items_count(OUT n INT)
| BEGIN
|   SELECT COUNT(*) INTO n FROM items;
| END
< ok 0 0
> exec
| SET SESSION explain_json_format_version = 1
< error
| 1193 HY000 Unknown system variable 'explain_json_format_version'
> query
| EXPLAIN FORMAT=JSON SELECT 1 FROM items WHERE id <= 2
< result
< column "EXPLAIN" "VARCHAR"
< row
: bytes
| {
|   "query_block": {
|     "select_id": 1,
|     "table": {
|       "table_name": "items",
|       "access_type": "range",
|       "possible_keys": ["PRIMARY"],
|       "key": "PRIMARY",
|       "key_length": "4",
|       "used_key_parts": ["id"],
|       "rows": 2,
|       "filtered": 100,
|       "attached_condition": "items.`id` <= 2",
|       "using_index": true
|     }
|   }
| }
> begin
< ok
> query
| EXPLAIN UPDATE items SET note = 'x' WHERE id = 1
< result
< column "id" "UNSIGNED BIGINT"
< column "select_type" "VARCHAR"
< column "table" "VARCHAR"
< column "type" "VARCHAR"
< column "possible_keys" "VARCHAR"
< column "key" "VARCHAR"
< column "key_len" "VARCHAR"
< column "ref" "VARCHAR"
< column "rows" "VARCHAR"
< column "Extra" "VARCHAR"
< row
: uint64 1
: bytes
| SIMPLE
: bytes
| items
: bytes
| range
: bytes
| PRIMARY
: bytes
| PRIMARY
: bytes
| 4
: null
: bytes
| 1
: bytes
| Using where
> rollback
< ok
> query
| SELECT
| 				42 AS i,
| 				3.5 AS f,
| 				'hello' AS s,
| 				CAST('2026-03-04 05:06:07.123456' AS DATETIME(6)) AS dt,
| 				CAST(0xdeadbeef AS BINARY(4)) AS bin,
| 				NULL AS nul;
< result
< column "i" "INT"
< column "f" "DECIMAL"
< column "s" "VARCHAR"
< column "dt" "DATETIME"
< column "bin" "VARBINARY"
< column "nul" "NULL"
< row
: int64 42
: bytes
| 3.5
: bytes
| hello
: bytes
| 2026-03-04 05:06:07.123456
: bytes
|q "ޭ\xbe\xef"
: null
> query
| SELECT id FROM items ORDER BY id LIMIT 2;
< result
< column "id" "INT"
< row
: int64 1
< row
: int64 2
> query
| EXPLAIN FORMAT=JSON SELECT note FROM items WHERE id > 1
< result
< column "EXPLAIN" "VARCHAR"
< row
: bytes
| {
|   "query_block": {
|     "select_id": 1,
|     "table": {
|       "table_name": "items",
|       "access_type": "range",
|       "possible_keys": ["PRIMARY"],
|       "key": "PRIMARY",
|       "key_length": "4",
|       "used_key_parts": ["id"],
|       "rows": 3,
|       "filtered": 100,
|       "attached_condition": "items.`id` > 1"
|     }
|   }
| }
> exec
| DELETE FROM items WHERE id = 10 RETURNING note;
< ok 0 0
> query
| SELECT COUNT(*) FROM items WHERE id = 99;
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 0
> query
| SELECT 1;
< result
< column "1" "INT"
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
< column "id" "INT"
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
|
| SELECT 7;
< result
< column "7" "INT"
< row
: int64 7
