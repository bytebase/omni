-- pin: default starrocks/allin1-ubuntu:3.4.10 sha256:8bc803ea7a5bb490c77117a456b32b792ccae698b1f0e47dfcafdac41542566d 3.4.10

== TestQueryExecuteAndStats default#1
> begin
< ok
> exec
| CREATE TABLE items (id BIGINT NOT NULL, note VARCHAR(64) NULL) ENGINE=OLAP PRIMARY KEY(id) DISTRIBUTED BY HASH(id) BUCKETS 1 PROPERTIES("replication_num" = "1");
| INSERT INTO items VALUES (1, 'a'), (2, 'b'), (3, 'c');
| UPDATE items SET note = 'z' WHERE id <= 2;
< ok 2 0
> commit
< ok
> begin
< ok
> exec
| INSERT INTO items VALUES (10, 'rb');
| INSERT INTO missing_items VALUES (11, 'rb');
< error
| 1064 HY000 Getting analyzing error. Detail message: Table missing_items is not found.
> rollback
< ok
> exec
| INSERT INTO items VALUES (20, 'kept');
| INSERT INTO missing_items VALUES (21, 'lost');
< error
| 1064 HY000 Getting analyzing error. Detail message: Table missing_items is not found.
> query
| SELECT COUNT(*) FROM items;
< result
< column "count(*)" "BIGINT"
< row
: int64 5
> query
| SELECT
| 		42 AS i,
| 		3.5 AS f,
| 		'hello' AS s,
| 		CAST('2026-03-04 05:06:07' AS DATETIME) AS dt,
| 		CAST('' AS VARCHAR) AS empty_s,
| 		NULL AS nul;
< result
< column "i" "TINYINT"
< column "f" "DECIMAL"
< column "s" "VARCHAR"
< column "dt" "DATETIME"
< column "empty_s" "VARCHAR"
< column "nul" "TINYINT"
< row
: int64 42
: bytes
| 3.5
: bytes
| hello
: bytes
| 2026-03-04 05:06:07
: bytes
|
: null
> query
| SELECT id FROM items WHERE id > 0 ORDER BY id LIMIT 2;
< result
< column "id" "BIGINT"
< row
: int64 1
< row
: int64 2
> query
| EXPLAIN SELECT note FROM items WHERE id > 1;
< result
< column "Explain String" "VARCHAR"
< row
: bytes
| PLAN FRAGMENT 0
< row
: bytes
|  OUTPUT EXPRS:2: note
< row
: bytes
|   PARTITION: RANDOM
< row
: bytes
|
< row
: bytes
|   RESULT SINK
< row
: bytes
|
< row
: bytes
|   1:Project
< row
: bytes
|   |  <slot 2> : 2: note
< row
: bytes
|   |  
< row
: bytes
|   0:OlapScanNode
< row
: bytes
|      TABLE: items
< row
: bytes
|      PREAGGREGATION: ON
< row
: bytes
|      PREDICATES: 1: id > 1
< row
: bytes
|      partitions=1/1
< row
: bytes
|      rollup: items
< row
: bytes
|      tabletRatio=1/1
< row
: bytes
|      tabletList=omni_tablet_1
< row
: bytes
|      cardinality=4
< row
: bytes
|      avgRowSize=2.0
> exec
| UPDATE items SET note = 'x' WHERE id = 1;
< ok 1 0
> query
| SELECT 1;
< result
< column "1" "TINYINT"
< row
: int64 1
> query
|
| SELECT * FROM missing_tbl;
< error
| 5502 42602 Getting analyzing error. Detail message: Unknown table 'omni_qe_test.missing_tbl'.
> query
| SELECT id, repeat('x', 80) FROM items ORDER BY id;
< result
< column "id" "BIGINT"
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
< row
: int64 20
: bytes
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
> query
|
| SELECT 7;
< result
< column "7" "TINYINT"
< row
: int64 7
