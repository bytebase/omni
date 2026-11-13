-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestExplainAnalyzeOfAnAlterTableRuns default#1
> exec
| CREATE DATABASE IF NOT EXISTS explain_forms
< ok 0 0
> exec
| DROP TABLE IF EXISTS explain_forms.ddl
< ok 0 0
> exec
| CREATE TABLE explain_forms.ddl (id int)
< ok 0 0
> exec
| EXPLAIN ALTER TABLE explain_forms.ddl ADD COLUMN plain int
< ok 0 0
> query
| SELECT column_name FROM information_schema.columns WHERE table_schema = 'explain_forms' AND table_name = 'ddl' ORDER BY column_name
< result
< column "column_name" "VARCHAR"
< row
: bytes
| id
> exec
| EXPLAIN ANALYZE ALTER TABLE explain_forms.ddl ADD COLUMN analyzed int
< ok 0 0
> query
| SELECT column_name FROM information_schema.columns WHERE table_schema = 'explain_forms' AND table_name = 'ddl' ORDER BY column_name
< result
< column "column_name" "VARCHAR"
< row
: bytes
| analyzed
< row
: bytes
| id
