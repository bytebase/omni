-- pin: driver gcp:bigquery:bytebase-dev:US sha256:92a85df45321648c0feb59ee0c222e8156cb9332d322ddb9f15eaf37c8c12533 2026-11-26

== TestQuery driver#1
> raw
| {"op":"ping"}
< raw
| {}
> raw
| {"op":"read","dataset":"omni_record","statement":"SELECT 1 AS n"}
< raw
| {"jobID":"job-1","schema":[{"name":"n","type":"INTEGER"}],"rows":[[{"k":"int64","i":1}]]}
> raw
| {"op":"jobStatus","jobID":"job-1"}
< raw
| {"status":{"stats":{"statementType":"SELECT","queryPlan":[{"id":0,"name":"S00: Output","status":"COMPLETE","steps":[{"kind":"COMPUTE","substeps":["$1 := 1"]},{"kind":"WRITE","substeps":["$1","TO __stage00_output"]}],"recordsWritten":1}]}}}
> raw
| {"op":"read","dataset":"omni_record","statement":"SELECT\n\t\tNUMERIC '12.50' AS num,\n\t\tBIGNUMERIC '-7.25' AS bn,\n\t\tDATE '2026-01-02' AS d,\n\t\tDATETIME '2026-01-02 03:04:05.123' AS dt,\n\t\tTIME '03:04:05' AS t,\n\t\tTIMESTAMP '2026-01-02 03:04:05.120000 UTC' AS ts,\n\t\tJSON '{\"a\":1}' AS j,\n\t\t[1, 2] AS arr,\n\t\tSTRUCT(1 AS a, 'x' AS b) AS st,\n\t\tST_GEOGPOINT(1, 2) AS g,\n\t\tINTERVAL 1 DAY AS iv,\n\t\tRANGE\u003cDATE\u003e '[2026-01-01, 2026-02-01)' AS r,\n\t\tb'ab' AS bin,\n\t\tTRUE AS flag,\n\t\t1.5 AS flt"}
< raw
| {"jobID":"job-2","schema":[{"name":"num","type":"NUMERIC"},{"name":"bn","type":"BIGNUMERIC"},{"name":"d","type":"DATE"},{"name":"dt","type":"DATETIME"},{"name":"t","type":"TIME"},{"name":"ts","type":"TIMESTAMP"},{"name":"j","type":"JSON"},{"name":"arr","type":"INTEGER"},{"name":"st","type":"RECORD"},{"name":"g","type":"GEOGRAPHY"},{"name":"iv","type":"INTERVAL"},{"name":"r","type":"RANGE"},{"name":"bin","type":"BYTES"},{"name":"flag","type":"BOOLEAN"},{"name":"flt","type":"FLOAT"}],"rows":[[{"k":"rat","s":"25/2"},{"k":"rat","s":"-29/4"},{"k":"date","s":"2026-01-02"},{"k":"dateTime","s":"2026-01-02T03:04:05.123000000"},{"k":"civilTime","s":"03:04:05"},{"k":"time","s":"2026-01-02T03:04:05.12Z"},{"k":"string","s":"{\"a\":1}"},{"k":"list","l":[{"k":"int64","i":1},{"k":"int64","i":2}]},{"k":"list","l":[{"k":"int64","i":1},{"k":"string","s":"x"}]},{"k":"string","s":"POINT(1 2)"},{"k":"interval","iv":{"d":1}},{"k":"range","rg":{"start":{"k":"date","s":"2026-01-01"},"end":{"k":"date","s":"2026-02-01"}}},{"k":"bytes","s":"YWI="},{"k":"bool","b":true},{"k":"float64","f":1.5}]]}
> raw
| {"op":"jobStatus","jobID":"job-2"}
< raw
| {"status":{"stats":{"statementType":"SELECT","queryPlan":[{"id":0,"name":"S00: Output","status":"COMPLETE","steps":[{"kind":"COMPUTE","substeps":["$1 := STRUCT\u003c...\u003e","$2 := NUMERIC\u003c...\u003e","$3 := BIGNUMERIC\u003c...\u003e","$4 := 20455","$5 := 2026-01-02 03:04:05.123","$6 := 03:04:05","$7 := 1767323045.120000000","$8 := UNKNOWN_TYPE\u003c...\u003e","$9 := ARRAY\u003c...\u003e","$10 := GEOGRAPHY\u003c...\u003e","$11 := INTERVAL\u003c...\u003e","$12 := RANGE\u003c...\u003e","$13 := BYTES\u003c...\u003e","$14 := 1.5","$15 := 1"]},{"kind":"WRITE","substeps":["$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15","TO __stage00_output"]}],"recordsWritten":1,"parallelInputs":1,"completedParallelInputs":1,"shuffleOutputBytes":225}]}}}
> raw
| {"op":"read","dataset":"omni_record","statement":"SELECT 1 AS n |\u003e WHERE n = 1"}
< raw
| {"jobID":"job-3","schema":[{"name":"n","type":"INTEGER"}],"rows":[[{"k":"int64","i":1}]]}
> raw
| {"op":"jobStatus","jobID":"job-3"}
< raw
| {"status":{"stats":{"statementType":"SELECT","queryPlan":[{"id":0,"name":"S00: Output","status":"COMPLETE","steps":[{"kind":"COMPUTE","substeps":["$1 := 1"]},{"kind":"WRITE","substeps":["$1","TO __stage00_output"]}],"recordsWritten":1}]}}}
> raw
| {"op":"read","dataset":"omni_record","statement":"WITH result AS (SELECT n FROM UNNEST([1, 2, 3]) AS n |\u003e WHERE n \u003e 1) SELECT * FROM result LIMIT 1;"}
< raw
| {"jobID":"job-4","schema":[{"name":"n","type":"INTEGER"}],"rows":[[{"k":"int64","i":2}]]}
> raw
| {"op":"jobStatus","jobID":"job-4"}
< raw
| {"status":{"stats":{"statementType":"SELECT","queryPlan":[{"id":0,"name":"S00: Output","status":"COMPLETE","steps":[{"kind":"COMPUTE","substeps":["$1 := ARRAY\u003c...\u003e"]},{"kind":"LIMIT","substeps":["1"]},{"kind":"FILTER","substeps":["greater($10, 1)"]},{"kind":"WRITE","substeps":["$10","TO __stage00_output"]}],"recordsWritten":1}]}}}
> raw
| {"op":"read","dataset":"omni_record","statement":"WITH result AS (SELECT 1 AS n UNION ALL SELECT 2) SELECT * FROM result LIMIT 1;"}
< raw
| {"jobID":"job-5","schema":[{"name":"n","type":"INTEGER"}],"rows":[[{"k":"int64","i":1}]]}
> raw
| {"op":"jobStatus","jobID":"job-5"}
< raw
| {"status":{"stats":{"statementType":"SELECT","queryPlan":[{"id":0,"name":"S00: Output","status":"COMPLETE","steps":[{"kind":"LIMIT","substeps":["1"]},{"kind":"COMPUTE","substeps":["$10 := 1"]},{"kind":"COMPUTE","substeps":["$20 := 2"]},{"kind":"WRITE","substeps":["$1","TO __stage00_output"]}],"recordsWritten":1}]}}}
> raw
| {"op":"read","dataset":"omni_record","statement":"WITH result AS (SELECT 3 AS n UNION ALL SELECT 4) SELECT * FROM result LIMIT 1;"}
< raw
| {"jobID":"job-6","schema":[{"name":"n","type":"INTEGER"}],"rows":[[{"k":"int64","i":3}]]}
> raw
| {"op":"jobStatus","jobID":"job-6"}
< raw
| {"status":{"stats":{"statementType":"SELECT","queryPlan":[{"id":0,"name":"S00: Output","status":"COMPLETE","steps":[{"kind":"LIMIT","substeps":["1"]},{"kind":"COMPUTE","substeps":["$10 := 3"]},{"kind":"COMPUTE","substeps":["$20 := 4"]},{"kind":"WRITE","substeps":["$1","TO __stage00_output"]}],"recordsWritten":1}]}}}
> raw
| {"op":"read","dataset":"omni_record","statement":"SELECT REPEAT('x', 4096) AS s"}
< raw
| {"jobID":"job-7","schema":[{"name":"s","type":"STRING"}],"rows":[[{"k":"string","s":"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"}]]}
> raw
| {"op":"jobStatus","jobID":"job-7"}
< raw
| {"status":{"stats":{"statementType":"SELECT","queryPlan":[{"id":0,"name":"S00: Output","status":"COMPLETE","steps":[{"kind":"COMPUTE","substeps":["$1 := 'xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'"]},{"kind":"WRITE","substeps":["$1","TO __stage00_output"]}],"recordsWritten":1}]}}}
> raw
| {"op":"read","dataset":"omni_record","statement":"SELECT 1 AS n |\u003e WHERE n = 2"}
< raw
| {"jobID":"job-8","schema":[{"name":"n","type":"INTEGER"}]}
> raw
| {"op":"jobStatus","jobID":"job-8"}
< raw
| {"status":{"stats":{"statementType":"SELECT","queryPlan":[{"id":0,"name":"S00: Output","status":"COMPLETE","steps":[{"kind":"COMPUTE","substeps":["$1 := NULL"]},{"kind":"FILTER","substeps":["0"]},{"kind":"WRITE","substeps":["$1","TO __stage00_output"]}]}]}}}
> raw
| {"op":"dryRun","dataset":"omni_record","statement":"SELECT COUNT(*) AS n FROM UNNEST([1, 2, 3])"}
< raw
| {"status":{"stats":{"statementType":"SELECT"}}}
> raw
| {"op":"run","dataset":"omni_record","statement":"SELEKT 1"}
< raw
| {"jobID":"job-9"}
> raw
| {"op":"jobWait","jobID":"job-9"}
< raw
| {"error":{"kind":"googleapi","code":400,"message":"Syntax error: Unexpected identifier \"SELEKT\" at [1:1]","errors":[{"reason":"invalidQuery","message":"Syntax error: Unexpected identifier \"SELEKT\" at [1:1]"}]}}
