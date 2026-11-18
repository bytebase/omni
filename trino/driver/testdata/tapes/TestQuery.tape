-- pin: default trinodb/trino:482 sha256:90b35b7c603eaa1f889bf03981a62b75f998ee6c0f851d9f4e341b49a57022b6 482

== TestQuery/SHOW_takes_no_limit_rewrite default#1
> query
| SHOW SCHEMAS
< result
< column "Schema" "VARCHAR"
< row
: string
| default
< row
: string
| information_schema
< row
: string
| omni_query
< row
: string
| omni_query_count
< row
: string
| omni_query_family
< row
: string
| omni_query_scoped

== TestQuery/a_comment_ending_at_a_carriage_return default#1
> query
|q "-- note\rSELECT 1"
< result
< column "_col0" "INTEGER"
< row
: int64 1

== TestQuery/a_count_only_where_the_statement_states_one default#1
> exec
| CREATE SCHEMA memory.omni_query_count
< ok 0 -
> exec
| CREATE TABLE memory.omni_query_count.plain (id integer)
< ok 0 -
> exec
| CREATE TABLE memory.omni_query_count.ctas AS SELECT x FROM UNNEST(SEQUENCE(1, 3)) AS t(x)
< ok 3 -
> exec
| SET SESSION query_max_run_time = '1m'
< ok 0 -

== TestQuery/a_leading_comment_does_not_change_the_statement_kind default#1
> query
| -- a comment
| /* another */ SELECT 1
< result
< column "_col0" "INTEGER"
< row
: int64 1

== TestQuery/a_user_EXPLAIN_plans_the_same_statement default#1
> query
| EXPLAIN SELECT 1
< result
< column "Query Plan" "VARCHAR"
< row
: string
| Trino version: 482
| Fragment 0 [SINGLE]
|     Output layout: [expr]
|     Output partitioning: SINGLE []
|     Output[columnNames = [_col0]]
|     │   Layout: [expr:integer]
|     │   Estimates: {rows: 1 (5B), cpu: 0, memory: 0B, network: 0B}
|     │   _col0 := expr
|     └─ Values[]
|            Layout: [expr:integer]
|            Estimates: {rows: 1 (5B), cpu: 0, memory: 0B, network: 0B}
|            (integer '1')
|
|

== TestQuery/an_exact_DECIMAL_keeps_its_value default#1
> query
| SELECT CAST('9007199254740993' AS DECIMAL(16,0)) AS v
< result
< column "v" "DECIMAL" precision 16 scale 0
< row
: string
| 9007199254740993

== TestQuery/change_and_query_statements default#1
> exec
| CREATE SCHEMA memory.omni_query
< ok 0 -
> exec
| CREATE TABLE memory.omni_query.items (id integer, note varchar)
< ok 0 -
> exec
| INSERT INTO memory.omni_query.items VALUES (1, 'a')
< ok 1 -
> exec
| INSERT INTO memory.omni_query.items VALUES (2, 'b')
< ok 1 -
> query
| SELECT id, note FROM memory.omni_query.items ORDER BY id
< result
< column "id" "INTEGER"
< column "note" "VARCHAR"
< row
: int64 1
: string
| a
< row
: int64 2
: string
| b

== TestQuery/explain default#1
> query
| EXPLAIN SELECT 1
< result
< column "Query Plan" "VARCHAR"
< row
: string
| Trino version: 482
| Fragment 0 [SINGLE]
|     Output layout: [expr]
|     Output partitioning: SINGLE []
|     Output[columnNames = [_col0]]
|     │   Layout: [expr:integer]
|     │   Estimates: {rows: 1 (5B), cpu: 0, memory: 0B, network: 0B}
|     │   _col0 := expr
|     └─ Values[]
|            Layout: [expr:integer]
|            Estimates: {rows: 1 (5B), cpu: 0, memory: 0B, network: 0B}
|            (integer '1')
|
|
> query
| EXPLAIN (FORMAT JSON) SELECT 1
< result
< column "Query Plan" "VARCHAR"
< row
: string
| {
|   "0" : {
|     "id" : "5",
|     "name" : "Output",
|     "descriptor" : {
|       "columnNames" : "[_col0]"
|     },
|     "outputs" : [ {
|       "type" : "integer",
|       "name" : "expr"
|     } ],
|     "details" : [ "_col0 := expr" ],
|     "estimates" : [ {
|       "outputRowCount" : 1.0,
|       "outputSizeInBytes" : 5.0,
|       "cpuCost" : 0.0,
|       "memoryCost" : 0.0,
|       "networkCost" : 0.0
|     } ],
|     "children" : [ {
|       "id" : "0",
|       "name" : "Values",
|       "descriptor" : { },
|       "outputs" : [ {
|         "type" : "integer",
|         "name" : "expr"
|       } ],
|       "details" : [ "(integer '1')" ],
|       "estimates" : [ {
|         "outputRowCount" : 1.0,
|         "outputSizeInBytes" : 5.0,
|         "cpuCost" : 0.0,
|         "memoryCost" : 0.0,
|         "networkCost" : 0.0
|       } ],
|       "children" : [ ]
|     } ]
|   }
| }

== TestQuery/explain_analyze_is_refused default#1

== TestQuery/failure_stops_the_script default#1
> query
| SELECT 1
< result
< column "_col0" "INTEGER"
< row
: int64 1
> query
| SELECT * FROM memory.default.no_such_table
< error
| {"Status":200,"Type":"USER_ERROR","Message":"line 1:15: Table 'memory.default.no_such_table' does not exist","SqlState":"","Name":"TABLE_NOT_FOUND","Code":46,"Line":1,"Column":15}

== TestQuery/limit default#1
> query
| SELECT x FROM UNNEST(SEQUENCE(1, 5)) AS t(x) LIMIT 2
< result
< column "x" "BIGINT"
< row
: int64 1
< row
: int64 2
> query
| SELECT x FROM UNNEST(SEQUENCE(1, 5)) AS t(x) LIMIT 3
< result
< column "x" "BIGINT"
< row
: int64 1
< row
: int64 2
< row
: int64 3

== TestQuery/limit_goes_before_a_trailing_comment default#1
> query
| SELECT x FROM UNNEST(SEQUENCE(1, 5)) AS t(x) LIMIT 2 -- note
< result
< column "x" "BIGINT"
< row
: int64 1
< row
: int64 2

== TestQuery/ping default#1
> query
| SELECT 1
< result
< column "_col0" "INTEGER"
< row
: int64 1

== TestQuery/result_size_cap default#1
> query
| SELECT x FROM UNNEST(SEQUENCE(1, 100)) AS t(x)
< result
< column "x" "BIGINT"
< row
: int64 1
< row
: int64 2
< row
: int64 3
< row
: int64 4
< row
: int64 5
< row
: int64 6
< row
: int64 7
< row
: int64 8
< row
: int64 9
< row
: int64 10
< row
: int64 11
< row
: int64 12
< row
: int64 13
< row
: int64 14
< row
: int64 15
< row
: int64 16
< row
: int64 17
< row
: int64 18
< row
: int64 19
< row
: int64 20
< row
: int64 21
< row
: int64 22
< row
: int64 23
< row
: int64 24
< row
: int64 25
< row
: int64 26
< row
: int64 27
< row
: int64 28
< row
: int64 29
< row
: int64 30
< row
: int64 31
< row
: int64 32
< row
: int64 33
< row
: int64 34
< row
: int64 35
< row
: int64 36
< row
: int64 37
< row
: int64 38
< row
: int64 39
< row
: int64 40
< row
: int64 41
< row
: int64 42
< row
: int64 43
< row
: int64 44
< row
: int64 45
< row
: int64 46
< row
: int64 47
< row
: int64 48
< row
: int64 49
< row
: int64 50
< row
: int64 51
< row
: int64 52
< row
: int64 53
< row
: int64 54
< row
: int64 55
< row
: int64 56
< row
: int64 57
< row
: int64 58
< row
: int64 59
< row
: int64 60
< row
: int64 61
< row
: int64 62
< row
: int64 63
< row
: int64 64
< row
: int64 65
< row
: int64 66
< row
: int64 67
< row
: int64 68
< row
: int64 69
< row
: int64 70
< row
: int64 71
< row
: int64 72
< row
: int64 73
< row
: int64 74
< row
: int64 75
< row
: int64 76
< row
: int64 77
< row
: int64 78
< row
: int64 79
< row
: int64 80
< row
: int64 81
< row
: int64 82
< row
: int64 83
< row
: int64 84
< row
: int64 85
< row
: int64 86
< row
: int64 87
< row
: int64 88
< row
: int64 89
< row
: int64 90
< row
: int64 91
< row
: int64 92
< row
: int64 93
< row
: int64 94
< row
: int64 95
< row
: int64 96
< row
: int64 97
< row
: int64 98
< row
: int64 99
< row
: int64 100

== TestQuery/schema default#1
> exec
| CREATE SCHEMA memory.omni_query_scoped
< ok 0 -
> exec
| CREATE TABLE memory.omni_query_scoped.scoped (id integer)
< ok 0 -
> exec
| INSERT INTO memory.omni_query_scoped.scoped VALUES (7)
< ok 1 -
> exec
| USE "omni_query_scoped"
< ok 0 -
> query
| SELECT id FROM scoped
< result
< column "id" "INTEGER"
< row
: int64 7

== TestQuery/the_SELECT_family_returns_rows default#1
> exec
| CREATE SCHEMA memory.omni_query_family
< ok 0 -
> exec
| CREATE TABLE memory.omni_query_family.items (id integer, note varchar)
< ok 0 -
> exec
| INSERT INTO memory.omni_query_family.items VALUES (1, 'a'), (2, 'b')
< ok 2 -
> query
| WITH t AS (SELECT 1 AS x) SELECT * FROM t
< result
< column "x" "INTEGER"
< row
: int64 1
> query
| TABLE memory.omni_query_family.items
< result
< column "id" "INTEGER"
< column "note" "VARCHAR"
< row
: int64 1
: string
| a
< row
: int64 2
: string
| b
> query
| VALUES 1, 2
< result
< column "_col0" "INTEGER"
< row
: int64 1
< row
: int64 2
> query
| (SELECT 1)
< result
< column "_col0" "INTEGER"
< row
: int64 1
> query
| DESC memory.omni_query_family.items
< result
< column "Column" "VARCHAR"
< column "Type" "VARCHAR"
< column "Extra" "VARCHAR"
< column "Comment" "VARCHAR"
< row
: string
| id
: string
| integer
: string
|
: string
|
< row
: string
| note
: string
| varchar
: string
|
: string
|

== TestQuery/the_limit_follows_a_top-level_OFFSET default#1
> query
| SELECT x FROM UNNEST(SEQUENCE(1, 5)) AS t(x) OFFSET 1 LIMIT 2
< result
< column "x" "BIGINT"
< row
: int64 2
< row
: int64 3

== TestQuery/the_limit_is_measured_by_tokens default#1
> query
| SELECT ' LIMIT ' AS v FROM UNNEST(SEQUENCE(1, 5)) AS t(x) LIMIT 2
< result
< column "v" "VARCHAR"
< row
: string
|  LIMIT 
< row
: string
|  LIMIT 
> query
| SELECT x FROM UNNEST(SEQUENCE(1, 5)) AS t(x) LIMIT
| 3
< result
< column "x" "BIGINT"
< row
: int64 1
< row
: int64 2
< row
: int64 3
> query
| SELECT (SELECT 1 LIMIT 1) AS v LIMIT 2
< result
< column "v" "INTEGER"
< row
: int64 1
> query
| SELECT t.limit FROM (VALUES 1, 2, 3) AS t(limit) LIMIT 2
< result
< column "limit" "INTEGER"
< row
: int64 1
< row
: int64 2
> query
| SELECT x FROM UNNEST(SEQUENCE(1, 5)) AS t(x) LIMIT ALL
< result
< column "x" "BIGINT"
< row
: int64 1
< row
: int64 2
< row
: int64 3
< row
: int64 4
< row
: int64 5
> query
| SELECT x FROM UNNEST(SEQUENCE(1, 5)) AS t(x) FETCH FIRST 2 ROWS ONLY
< result
< column "x" "BIGINT"
< row
: int64 1
< row
: int64 2
> query
| SELECT 1 AS x UNION CORRESPONDING SELECT 1 AS x LIMIT 1
< result
< column "x" "INTEGER"
< row
: int64 1
> query
| SELECT 1 AS x UNION CORRESPONDING SELECT 1 AS x LIMIT 2
< result
< column "x" "INTEGER"
< row
: int64 1

== TestQuery/types default#1
> query
| SELECT
| 			CAST(42 AS INTEGER) AS c_integer,
| 			CAST(42 AS BIGINT) AS c_bigint,
| 			true AS c_boolean,
| 			CAST('hi' AS VARCHAR) AS c_varchar,
| 			CAST('ab' AS CHAR(5)) AS c_char,
| 			CAST(1.25 AS DECIMAL(5,2)) AS c_decimal,
| 			CAST(1.5 AS REAL) AS c_real,
| 			CAST(2.5 AS DOUBLE) AS c_double,
| 			DATE '2026-01-02' AS c_date,
| 			TIME '12:34:56.789' AS c_time,
| 			TIME '12:34:56 +02:00' AS c_time_tz,
| 			TIMESTAMP '2026-01-02 12:34:56' AS c_timestamp,
| 			CAST('2026-01-02 12:34:56.123 +02:00' AS TIMESTAMP(3) WITH TIME ZONE) AS c_timestamp_tz,
| 			X'0102' AS c_varbinary,
| 			X'' AS c_varbinary_empty,
| 			UUID '12151fd2-7586-11e9-8f9e-2a86e4085a59' AS c_uuid,
| 			CAST('10.0.0.1' AS IPADDRESS) AS c_ipaddress,
| 			INTERVAL '1' DAY AS c_interval,
| 			CAST('{"a":1}' AS JSON) AS c_json,
| 			CAST(NULL AS VARCHAR) AS c_null
< result
< column "c_integer" "INTEGER"
< column "c_bigint" "BIGINT"
< column "c_boolean" "BOOLEAN"
< column "c_varchar" "VARCHAR"
< column "c_char" "CHAR"
< column "c_decimal" "DECIMAL" precision 5 scale 2
< column "c_real" "REAL"
< column "c_double" "DOUBLE"
< column "c_date" "DATE"
< column "c_time" "TIME" precision 3 scale 0
< column "c_time_tz" "TIME WITH TIME ZONE" precision 0 scale 0
< column "c_timestamp" "TIMESTAMP" precision 0 scale 0
< column "c_timestamp_tz" "TIMESTAMP WITH TIME ZONE" precision 3 scale 0
< column "c_varbinary" "VARBINARY" "[]uint8"
< column "c_varbinary_empty" "VARBINARY" "[]uint8"
< column "c_uuid" "UUID"
< column "c_ipaddress" "IPADDRESS"
< column "c_interval" "INTERVAL DAY TO SECOND"
< column "c_json" "JSON"
< column "c_null" "VARCHAR"
< row
: int64 42
: int64 42
: bool true
: string
| hi
: string
| ab   
: string
| 1.25
: float64 1.5
: float64 2.5
: time 2026-01-02T00:00:00Z
: time 0000-01-01T12:34:56.789Z
: time 0000-01-01T12:34:56+02:00
: time 2026-01-02T12:34:56Z
: time 2026-01-02T12:34:56.123+02:00
: bytes
|q "\x01\x02"
: bytes
|
: string
| 12151fd2-7586-11e9-8f9e-2a86e4085a59
: string
| 10.0.0.1
: string
| 1 00:00:00.000
: string
| "{\"a\":1}"
: null
