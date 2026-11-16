-- pin: default dyrnq/doris:3.1.4 sha256:88f99819c0877567b77d78cc8367f0e2dd4b6c1602fe17c6f4b1e332d8e315dd 3.1.4-rc02

== TestQueryExecuteAndStats default#1
> begin
< ok
> exec
| CREATE TABLE items (id BIGINT NOT NULL, note VARCHAR(64) NULL) ENGINE=OLAP UNIQUE KEY(id) DISTRIBUTED BY HASH(id) BUCKETS 1 PROPERTIES("replication_num" = "1");
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
| INSERT INTO missing_items VALUES (11, 'rb');
< error
| 1105 HY000 errCode = 2, detailMessage = Table [missing_items] does not exist in database [omni_qe_test].
> rollback
< ok
> exec
| INSERT INTO items VALUES (20, 'kept');
< ok 1 0
> exec
|
| INSERT INTO missing_items VALUES (21, 'lost');
< error
| 1105 HY000 errCode = 2, detailMessage = Table [missing_items] does not exist in database [omni_qe_test].
> query
| SELECT COUNT(*) FROM items;
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 5
> query
| SELECT
| 		42 AS i,
| 		3.5 AS f,
| 		'hello' AS s,
| 		CAST('2026-03-04 05:06:07.123456' AS DATETIME(6)) AS dt,
| 		CAST('' AS VARCHAR) AS empty_s,
| 		NULL AS nul;
< result
< column "i" "TINYINT"
< column "f" "DECIMAL"
< column "s" "CHAR"
< column "dt" "DATETIME"
< column "empty_s" "CHAR"
< column "nul" "CHAR"
< row
: int64 42
: bytes
| 3.5
: bytes
| hello
: bytes
| 2026-03-04 05:06:07.123456
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
< column "Explain String(Nereids Planner)" "CHAR"
< row
: bytes
| PLAN FRAGMENT 0
< row
: bytes
|   OUTPUT EXPRS:
< row
: bytes
|     note[#4]
< row
: bytes
|   PARTITION: UNPARTITIONED
< row
: bytes
|
< row
: bytes
|   HAS_COLO_PLAN_NODE: false
< row
: bytes
|
< row
: bytes
|   VRESULT SINK
< row
: bytes
|      MYSQL_PROTOCAL
< row
: bytes
|
< row
: bytes
|   1:VEXCHANGE
< row
: bytes
|      offset: 0
< row
: bytes
|      distribute expr lists: 
< row
: bytes
|
< row
: bytes
| PLAN FRAGMENT 1
< row
: bytes
|
< row
: bytes
|   PARTITION: HASH_PARTITIONED: id[#0]
< row
: bytes
|
< row
: bytes
|   HAS_COLO_PLAN_NODE: false
< row
: bytes
|
< row
: bytes
|   STREAM DATA SINK
< row
: bytes
|     EXCHANGE ID: 01
< row
: bytes
|     UNPARTITIONED
< row
: bytes
|
< row
: bytes
|   0:VOlapScanNode(141)
< row
: bytes
|      TABLE: omni_qe_test.items(items), PREAGGREGATION: ON
< row
: bytes
|      PREDICATES: ((id[#0] > 1) AND (__DORIS_DELETE_SIGN__[#2] = 0))
< row
: bytes
|      partitions=1/1 (items)
< row
: bytes
|      tablets=1/1, tabletList=omni_tablet_1
< row
: bytes
|      cardinality=1, avgRowSize=0.0, numNodes=1
< row
: bytes
|      pushAggOp=NONE
< row
: bytes
|      final projections: note[#1]
< row
: bytes
|      final project output tuple id: 1
< row
: bytes
|
< row
: bytes
|
< row
: bytes
|
< row
: bytes
| ========== STATISTICS ==========
< row
: bytes
| planed with unknown column statistics
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
| 1105 HY000 errCode = 2, detailMessage = Table [missing_tbl] does not exist in database [omni_qe_test].
> query
| SELECT id, repeat('x', 80) FROM items ORDER BY id;
< result
< column "id" "BIGINT"
< column "repeat('x', 80)" "CHAR"
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
