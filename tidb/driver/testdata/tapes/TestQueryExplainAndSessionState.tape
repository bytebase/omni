-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestQueryExplainAndSessionState default#1
> exec
| USE `omni_qe_explain_session`
< ok 0 0
> exec
| CREATE TABLE plan_target (id INT PRIMARY KEY, note VARCHAR(64))
< ok 0 0
> exec
| SET @omni_session_probe = 42
< ok 0 0
> query
| SELECT @omni_session_probe
< result
< column "@omni_session_probe" "BIGINT" "int64"
< row
: int64 42
> query
| EXPLAIN FORMAT=tidb_json SELECT note FROM plan_target WHERE id > 100
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
|                 "operatorInfo": "omni_qe_explain_session.plan_target.note",
|                 "subOperators": [
|                     {
|                         "id": "TableRangeScan_8",
|                         "estRows": "3333.33",
|                         "taskType": "cop[tikv]",
|                         "accessObject": "table:plan_target",
|                         "operatorInfo": "range:(100,+inf], keep order:false, stats:pseudo"
|                     }
|                 ]
|             }
|         ]
|     }
| ]
|
> query
| EXPLAIN SELECT note FROM plan_target WHERE id > 100
< result
< column "id" "CHAR"
< column "estRows" "CHAR"
< column "task" "CHAR"
< column "access object" "CHAR"
< column "operator info" "CHAR"
< row
: bytes
| TableReader_9
: bytes
| 3333.33
: bytes
| root
: bytes
|
: bytes
| data:Projection_5
< row
: bytes
| └─Projection_5
: bytes
| 3333.33
: bytes
| cop[tikv]
: bytes
|
: bytes
| omni_qe_explain_session.plan_target.note
< row
: bytes
|   └─TableRangeScan_8
: bytes
| 3333.33
: bytes
| cop[tikv]
: bytes
| table:plan_target
: bytes
| range:(100,+inf], keep order:false, stats:pseudo
> query
| EXPLAIN FORMAT=tidb_json SELECT note FROM plan_target WHERE id > 100
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
|                 "operatorInfo": "omni_qe_explain_session.plan_target.note",
|                 "subOperators": [
|                     {
|                         "id": "TableRangeScan_8",
|                         "estRows": "3333.33",
|                         "taskType": "cop[tikv]",
|                         "accessObject": "table:plan_target",
|                         "operatorInfo": "range:(100,+inf], keep order:false, stats:pseudo"
|                     }
|                 ]
|             }
|         ]
|     }
| ]
|
> query ^
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
|                 "operatorInfo": "omni_qe_explain_session.plan_target.note",
|                 "subOperators": [
|                     {
|                         "id": "TableRangeScan_8",
|                         "estRows": "3333.33",
|                         "taskType": "cop[tikv]",
|                         "accessObject": "table:plan_target",
|                         "operatorInfo": "range:(100,+inf], keep order:false, stats:pseudo"
|                     }
|                 ]
|             }
|         ]
|     }
| ]
|
> query
| EXPLAIN SELECT note FROM plan_target WHERE id > 100
< result
< column "id" "CHAR"
< column "estRows" "CHAR"
< column "task" "CHAR"
< column "access object" "CHAR"
< column "operator info" "CHAR"
< row
: bytes
| TableReader_9
: bytes
| 3333.33
: bytes
| root
: bytes
|
: bytes
| data:Projection_5
< row
: bytes
| └─Projection_5
: bytes
| 3333.33
: bytes
| cop[tikv]
: bytes
|
: bytes
| omni_qe_explain_session.plan_target.note
< row
: bytes
|   └─TableRangeScan_8
: bytes
| 3333.33
: bytes
| cop[tikv]
: bytes
| table:plan_target
: bytes
| range:(100,+inf], keep order:false, stats:pseudo
