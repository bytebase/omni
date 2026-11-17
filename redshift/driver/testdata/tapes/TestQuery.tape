-- pin: driver aws:redshift:1.0.539507 sha256:afc4041fe6e35ce273970308e430eb8245ea9c4bbf597c497bb8eb8bb00810d2 2026-11-17

== TestQuery driver#1
> query
| SELECT 1 AS n
< result
< column "n" "INT4" "int32"
< row
: int64 1
> query
| SELECT 1 AS n UNION ALL SELECT 2 LIMIT 1
< result
< column "n" "INT4" "int32"
< row
: int64 1
> exec
| CREATE TEMP TABLE omni_query_probe (a int)
< ok 0 -
> query
| EXPLAIN SELECT COUNT(*) FROM omni_query_probe
< result
< column "QUERY PLAN" "VARCHAR" "string"
< row
: string
| XN Aggregate  (cost=17.50..17.50 rows=1 width=0)
< row
: string
|   ->  XN Seq Scan on omni_query_probe  (cost=0.00..14.00 rows=1400 width=0)
