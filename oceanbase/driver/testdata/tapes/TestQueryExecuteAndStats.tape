-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestQueryExecuteAndStats default#1
> begin
< ok
> exec
| CREATE TABLE items (id INT AUTO_INCREMENT PRIMARY KEY, note VARCHAR(64));
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
> query
| EXPLAIN FORMAT=JSON DELETE FROM items WHERE id <= 2;
< result
< column "Query Plan" "VARCHAR"
< row
: bytes
| {
< row
: bytes
|   "ID":0,
< row
: bytes
|   "OPERATOR":"DELETE",
< row
: bytes
|   "NAME":"",
< row
: bytes
|   "EST.ROWS":2,
< row
: bytes
|   "EST.TIME(us)":19,
< row
: bytes
|   "output":"",
< row
: bytes
|   "CHILD_1": {
< row
: bytes
|     "ID":1,
< row
: bytes
|     "OPERATOR":"TABLE RANGE SCAN",
< row
: bytes
|     "NAME":"items",
< row
: bytes
|     "EST.ROWS":2,
< row
: bytes
|     "EST.TIME(us)":3,
< row
: bytes
|     "output":"output([items.id], [items.note])"
< row
: bytes
|   }
< row
: bytes
| }
> begin
< ok
> query
| EXPLAIN UPDATE items SET note = 'x' WHERE id = 1
< result
< column "Query Plan" "VARCHAR"
< row
: bytes
| ============================================
< row
: bytes
| |ID|OPERATOR   |NAME |EST.ROWS|EST.TIME(us)|
< row
: bytes
| --------------------------------------------
< row
: bytes
| |0 |UPDATE     |     |1       |37          |
< row
: bytes
| |1 |└─TABLE GET|items|1       |5           |
< row
: bytes
| ============================================
< row
: bytes
| Outputs & filters:
< row
: bytes
| -------------------------------------
< row
: bytes
|   0 - output(nil), filter(nil)
< row
: bytes
|       table_columns([{items: ({items: (items.id, items.note)})}]),
< row
: bytes
|       update([items.note=column_conv(VARCHAR,utf8mb4_general_ci,length:64,NULL,'x')])
< row
: bytes
|   1 - output([items.id], [items.note]), filter(nil), rowset=16
< row
: bytes
|       access([items.id], [items.note]), partitions(p0)
< row
: bytes
|       is_index_back=false, is_global_index=false, 
< row
: bytes
|       range_key([items.id]), range[1 ; 1], 
< row
: bytes
|       range_cond([items.id = 1])
> rollback
< ok
> begin
< ok
> query
| EXPLAIN FORMAT=JSON DELETE FROM items WHERE id = 1
< result
< column "Query Plan" "VARCHAR"
< row
: bytes
| {
< row
: bytes
|   "ID":0,
< row
: bytes
|   "OPERATOR":"DELETE",
< row
: bytes
|   "NAME":"",
< row
: bytes
|   "EST.ROWS":1,
< row
: bytes
|   "EST.TIME(us)":17,
< row
: bytes
|   "output":"",
< row
: bytes
|   "CHILD_1": {
< row
: bytes
|     "ID":1,
< row
: bytes
|     "OPERATOR":"TABLE GET",
< row
: bytes
|     "NAME":"items",
< row
: bytes
|     "EST.ROWS":1,
< row
: bytes
|     "EST.TIME(us)":5,
< row
: bytes
|     "output":"output([items.id], [items.note])"
< row
: bytes
|   }
< row
: bytes
| }
> rollback
< ok
> query
| SELECT EXTRA FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND COLUMN_NAME = ?
: string
| omni_qe_test
: string
| items
: string
| id
< result
< column "EXTRA" "VARCHAR"
< row
: bytes
| auto_increment
> query
| SELECT
| 		42 AS i,
| 		3.5 AS f,
| 		'hello' AS s,
| 		CAST('2026-03-04 05:06:07.123456' AS DATETIME(6)) AS dt,
| 		CAST(0xdeadbeef AS BINARY(4)) AS bin,
| 		NULL AS nul;
< result
< column "i" "BIGINT"
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
< column "Query Plan" "VARCHAR"
< row
: bytes
| {
< row
: bytes
|   "ID":0,
< row
: bytes
|   "OPERATOR":"TABLE RANGE SCAN",
< row
: bytes
|   "NAME":"items",
< row
: bytes
|   "EST.ROWS":3,
< row
: bytes
|   "EST.TIME(us)":3,
< row
: bytes
|   "output":"output([items.note])"
< row
: bytes
| }
> query
| SELECT COUNT(*) FROM items WHERE id = 99;
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 0
> query
| SELECT 1;
< result
< column "1" "BIGINT"
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
< row
: int64 10
: bytes
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
> query
|
| SELECT 7;
< result
< column "7" "BIGINT"
< row
: int64 7
