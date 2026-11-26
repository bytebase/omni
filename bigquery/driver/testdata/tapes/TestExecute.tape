-- pin: driver gcp:bigquery:bytebase-dev:US sha256:92a85df45321648c0feb59ee0c222e8156cb9332d322ddb9f15eaf37c8c12533 2026-11-26

== TestExecute driver#1
> raw
| {"op":"ping"}
< raw
| {}
> raw
| {"op":"run","dataset":"omni_record","statement":"CREATE TABLE omni_execute_probe (a INT64)"}
< raw
| {"jobID":"job-1"}
> raw
| {"op":"jobWait","jobID":"job-1"}
< raw
| {"status":{"stats":{"statementType":"CREATE_TABLE"}}}
> raw
| {"op":"run","dataset":"omni_record","statement":"INSERT INTO omni_execute_probe (a) VALUES (1), (2)"}
< raw
| {"jobID":"job-2"}
> raw
| {"op":"jobWait","jobID":"job-2"}
< raw
| {"status":{"stats":{"statementType":"INSERT","numDmlAffectedRows":2,"queryPlan":[{"id":0,"name":"S00: Input","status":"COMPLETE","steps":[{"kind":"COMPUTE","substeps":["$10 := 1"]},{"kind":"COMPUTE","substeps":["$20 := 2"]},{"kind":"WRITE","substeps":["$1","TO __stage00_output"]}],"recordsWritten":2,"parallelInputs":1,"completedParallelInputs":1,"shuffleOutputBytes":18},{"id":1,"name":"S01: Coalesce","status":"COMPLETE","inputStages":[0],"steps":[{"kind":"READ","substeps":["FROM __stage00_output"]}],"recordsRead":2,"recordsWritten":2,"parallelInputs":50,"completedParallelInputs":50,"shuffleOutputBytes":18},{"id":2,"name":"S02: Output","status":"COMPLETE","inputStages":[1],"steps":[{"kind":"READ","substeps":["$1","FROM __stage01_output"]},{"kind":"EXPORT","substeps":["$1:a","TO omni_record.omni_execute_probe"]},{"kind":"WRITE","substeps":["$1","TO __stage02_output"]}],"recordsRead":2,"recordsWritten":2,"parallelInputs":1,"completedParallelInputs":1}]}}}
> raw
| {"op":"run","dataset":"omni_record","statement":"INSERT INTO omni_execute_probe (a) VALUES (3); INSERT INTO omni_execute_probe (a) VALUES ('not an int')"}
< raw
| {"jobID":"job-3"}
> raw
| {"op":"jobWait","jobID":"job-3"}
< raw
| {"error":{"kind":"googleapi","code":400,"message":"Query error: Value has type STRING which cannot be inserted into column a, which has type INT64 at [1:91]","errors":[{"reason":"invalidQuery","message":"Query error: Value has type STRING which cannot be inserted into column a, which has type INT64 at [1:91]"}]}}
> raw
| {"op":"read","dataset":"omni_record","statement":"SELECT COUNT(*) AS n FROM omni_execute_probe"}
< raw
| {"jobID":"job-4","schema":[{"name":"n","type":"INTEGER"}],"rows":[[{"k":"int64","i":3}]]}
> raw
| {"op":"jobStatus","jobID":"job-4"}
< raw
| {"status":{"stats":{"statementType":"SELECT","queryPlan":[{"id":0,"name":"S00: Output","status":"COMPLETE","steps":[{"kind":"READ","substeps":["","FROM omni_execute_probe"]},{"kind":"AGGREGATE","substeps":["$1 := SUM_OF_COUNTS($10)"]},{"kind":"AGGREGATE","substeps":["$10 := COUNT_STAR()"]},{"kind":"WRITE","substeps":["$1","TO __stage00_output"]}],"recordsRead":3,"recordsWritten":1,"parallelInputs":2,"completedParallelInputs":2}]}}}
> raw
| {"op":"run","dataset":"omni_record","statement":"DROP TABLE IF EXISTS omni_execute_probe"}
< raw
| {"jobID":"job-5"}
> raw
| {"op":"jobWait","jobID":"job-5"}
< raw
| {"status":{"stats":{"statementType":"DROP_TABLE"}}}
