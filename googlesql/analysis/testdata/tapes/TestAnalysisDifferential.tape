-- pin: default gcr.io/cloud-spanner-emulator/emulator:1.5.54 sha256:caf1bd24c081e005837b5977bae5a250e25cb4da9f25ec1abc91936ad67e4de2 v1.5.54

== TestAnalysisDifferential default#1
> query
| SELECT a FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:15] SELECT a FROM t               ^"}
> query
| SELECT * FROM t1 JOIN t2 ON t1.id = t2.id
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t1 [at 1:15] SELECT * FROM t1 JOIN t2 ON t1.id = t2.id               ^"}
> query
| WITH c AS (SELECT a FROM t) SELECT a FROM c
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:26] WITH c AS (SELECT a FROM t) SELECT a FROM c                          ^"}
> query
| SELECT a FROM t WHERE b IN (SELECT b FROM t2)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:15] SELECT a FROM t WHERE b IN (SELECT b FROM t2)               ^"}
> query
| SELECT a FROM t UNION ALL SELECT a2 FROM t2
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:15] SELECT a FROM t UNION ALL SELECT a2 FROM t2               ^"}
> query
| SELECT * FROM INFORMATION_SCHEMA.TABLES
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT * FROM SPANNER_SYS.QUERY_STATS_TOP_MINUTE
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: SPANNER_SYS.QUERY_STATS_TOP_MINUTE [at 1:15] SELECT * FROM SPANNER_SYS.QUERY_STATS_TOP_MINUTE               ^"}
> query
| INSERT INTO t (a) VALUES (1)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"dml","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:13] INSERT INTO t (a) VALUES (1)             ^"}
> query
| DELETE FROM t WHERE a = 1
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"dml","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:13] DELETE FROM t WHERE a = 1             ^"}
> query
| UPDATE t SET a = 2 WHERE a = 1
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"dml","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:8] UPDATE t SET a = 2 WHERE a = 1        ^"}
> query
| CREATE TABLE t3 (a INT64) PRIMARY KEY (a)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"ok","code":"OK"}
> query
| DROP TABLE t3
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Table not found: t3"}
> query
| ALTER TABLE t ADD COLUMN c INT64
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Table not found: t"}
