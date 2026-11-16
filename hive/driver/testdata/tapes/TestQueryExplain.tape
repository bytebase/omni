-- pin: default apache/hive:4.0.1 sha256:5194161ef50b80875f937dff04936df047cbe894e9727ab8f0606349345b1bd3 rel/release-4.0.1

== TestQueryExplain default#1
> exec
| CREATE EXTERNAL TABLE ext_users (id INT COMMENT 'the id', name STRING) ROW FORMAT DELIMITED FIELDS TERMINATED BY ',' STORED AS TEXTFILE LOCATION '/tmp/omni/ext_users'
< ok - -
> exec
| CREATE TABLE orders (id BIGINT, amount DECIMAL(10,2), note STRING COMMENT 'a note') STORED AS PARQUET
< ok - -
> exec
| ALTER TABLE orders SET TBLPROPERTIES('EXTERNAL'='FALSE')
< ok - -
> exec
| INSERT INTO orders VALUES (1, 1.5, 'n1')
< ok - -
> exec
| CREATE TABLE events (id INT, body STRING) PARTITIONED BY (dt STRING) STORED AS PARQUET
< ok - -
> exec
| ALTER TABLE events SET TBLPROPERTIES('EXTERNAL'='FALSE')
< ok - -
> exec
| ALTER TABLE events ADD PARTITION (dt='2024-01-01')
< ok - -
> exec
| INSERT INTO events PARTITION (dt='2024-01-01') VALUES (1, 'x'), (2, 'y')
< ok - -
> exec
| CREATE VIEW v_users AS SELECT id, name FROM ext_users
< ok - -
> query
| EXPLAIN SELECT id FROM orders
< result
< column "Explain" "STRING"
< row
: string
| STAGE DEPENDENCIES:
< row
: string
|   Stage-0 is a root stage
< row
: string
|
< row
: string
| STAGE PLANS:
< row
: string
|   Stage: Stage-0
< row
: string
|     Fetch Operator
< row
: string
|       limit: -1
< row
: string
|       Processor Tree:
< row
: string
|         TableScan
< row
: string
|           alias: orders
< row
: string
|           Statistics: Num rows: 1 Data size: 8 Basic stats: COMPLETE Column stats: COMPLETE
< row
: string
|           Select Operator
< row
: string
|             expressions: id (type: bigint)
< row
: string
|             outputColumnNames: _col0
< row
: string
|             Statistics: Num rows: 1 Data size: 8 Basic stats: COMPLETE Column stats: COMPLETE
< row
: string
|             ListSink
< row
: string
|
> query
| EXPLAIN FORMATTED SELECT id FROM orders
< result
< column "Explain" "STRING"
< row
: string
| {"CBOPlan":"{\n  \"rels\": [\n    {\n      \"id\": \"0\",\n      \"relOp\": \"org.apache.hadoop.hive.ql.optimizer.calcite.reloperators.HiveTableScan\",\n      \"table\": [\n        \"omni_query_explain\",\n        \"orders\"\n      ],\n      \"table:alias\": \"orders\",\n      \"inputs\": [],\n      \"rowCount\": 1.0,\n      \"avgRowSize\": 8.0,\n      \"rowType\": [\n        {\n          \"type\": \"BIGINT\",\n          \"nullable\": true,\n          \"name\": \"id\"\n        },\n        {\n          \"type\": \"DECIMAL\",\n          \"nullable\": true,\n          \"precision\": 10,\n          \"scale\": 2,\n          \"name\": \"amount\"\n        },\n        {\n          \"type\": \"VARCHAR\",\n          \"nullable\": true,\n          \"precision\": 1000000001,\n          \"name\": \"note\"\n        },\n        {\n          \"type\": \"BIGINT\",\n          \"nullable\": true,\n          \"name\": \"BLOCK__OFFSET__INSIDE__FILE\"\n        },\n        {\n          \"type\": \"VARCHAR\",\n          \"nullable\": true,\n          \"precision\": 1000000001,\n          \"name\": \"INPUT__FILE__NAME\"\n        },\n        {\n          \"fields\": [\n            {\n              \"type\": \"BIGINT\",\n              \"nullable\": true,\n              \"name\": \"writeid\"\n            },\n            {\n              \"type\": \"INTEGER\",\n              \"nullable\": true,\n              \"name\": \"bucketid\"\n            },\n            {\n              \"type\": \"BIGINT\",\n              \"nullable\": true,\n              \"name\": \"rowid\"\n            }\n          ],\n          \"name\": \"ROW__ID\"\n        },\n        {\n          \"type\": \"BOOLEAN\",\n          \"nullable\": true,\n          \"name\": \"ROW__IS__DELETED\"\n        }\n      ],\n      \"colStats\": [\n        {\n          \"name\": \"id\",\n          \"ndv\": 1,\n          \"minValue\": 1,\n          \"maxValue\": 1\n        },\n        {\n          \"name\": \"amount\",\n          \"ndv\": 1,\n          \"minValue\": 1.5,\n          \"maxValue\": 1.5\n        },\n        {\n          \"name\": \"note\",\n          \"ndv\": 1\n        }\n      ]\n    },\n    {\n      \"id\": \"1\",\n      \"relOp\": \"org.apache.hadoop.hive.ql.optimizer.calcite.reloperators.HiveProject\",\n      \"fields\": [\n        \"id\"\n      ],\n      \"exprs\": [\n        {\n          \"input\": 0,\n          \"name\": \"$0\"\n        }\n      ],\n      \"rowCount\": 1.0\n    }\n  ]\n}","optimizedSQL":"SELECT `id`\nFROM `omni_query_explain`.`orders`","cboInfo":"Plan optimized by CBO.","STAGE DEPENDENCIES":{"Stage-0":{"ROOT STAGE":"TRUE"}},"STAGE PLANS":{"Stage-0":{"Fetch Operator":{"limit:":"-1","Processor Tree:":{"TableScan":{"alias:":"orders","columns:":["id"],"database:":"omni_query_explain","Statistics:":"Num rows: 1 Data size: 8 Basic stats: COMPLETE Column stats: COMPLETE","table:":"orders","isTempTable:":"false","OperatorId:":"TS_0","children":{"Select Operator":{"expressions:":"id (type: bigint)","columnExprMap:":{"_col0":"id"},"outputColumnNames:":["_col0"],"Statistics:":"Num rows: 1 Data size: 8 Basic stats: COMPLETE Column stats: COMPLETE","OperatorId:":"SEL_1","children":{"ListSink":{"OperatorId:":"LIST_SINK_3"}}}}}}}}}}
> query
| EXPLAIN SELECT id FROM orders
< result
< column "Explain" "STRING"
< row
: string
| STAGE DEPENDENCIES:
< row
: string
|   Stage-0 is a root stage
< row
: string
|
< row
: string
| STAGE PLANS:
< row
: string
|   Stage: Stage-0
< row
: string
|     Fetch Operator
< row
: string
|       limit: -1
< row
: string
|       Processor Tree:
< row
: string
|         TableScan
< row
: string
|           alias: orders
< row
: string
|           Statistics: Num rows: 1 Data size: 8 Basic stats: COMPLETE Column stats: COMPLETE
< row
: string
|           Select Operator
< row
: string
|             expressions: id (type: bigint)
< row
: string
|             outputColumnNames: _col0
< row
: string
|             Statistics: Num rows: 1 Data size: 8 Basic stats: COMPLETE Column stats: COMPLETE
< row
: string
|             ListSink
< row
: string
|
