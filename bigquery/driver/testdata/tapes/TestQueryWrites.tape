-- pin: driver gcp:bigquery:bytebase-dev:US sha256:a38c0d3800b3dfd8d1c6e3f3644af8fd02fc0f8ed110efc924159d3415533cd7 2026-11-21

== TestQueryWrites driver#1
> raw
| {"op":"ping"}
< raw
| {}
> raw
| {"op":"run","dataset":"omni_record","statement":"CREATE TABLE omni_query_probe (a INT64)"}
< raw
| {"jobID":"job-1"}
> raw
| {"op":"jobWait","jobID":"job-1"}
< raw
| {"status":{"stats":{"statementType":"CREATE_TABLE"}}}
> raw
| {"op":"jobRead","jobID":"job-1"}
< raw
| {"jobID":"job-1","schema":[{"name":"a","type":"INTEGER"}]}
> raw
| {"op":"dryRun","dataset":"omni_record","statement":"SELECT COUNT(*) FROM omni_query_probe"}
< raw
| {"status":{"stats":{"statementType":"SELECT"}}}
> raw
| {"op":"read","dataset":"omni_record","statement":"SELECT * FROM omni_query_probe_no_such_table"}
< raw
| {"error":{"kind":"googleapi","code":404,"message":"Not found: Table bytebase-dev:omni_record.omni_query_probe_no_such_table was not found in location US","errors":[{"reason":"notFound","message":"Not found: Table bytebase-dev:omni_record.omni_query_probe_no_such_table was not found in location US"}]}}
> raw
| {"op":"run","dataset":"omni_record","statement":"DROP TABLE IF EXISTS omni_query_probe"}
< raw
| {"jobID":"job-2"}
> raw
| {"op":"jobWait","jobID":"job-2"}
< raw
| {"status":{"stats":{"statementType":"DROP_TABLE"}}}
> raw
| {"op":"jobRead","jobID":"job-2"}
< raw
| {"jobID":"job-2"}
