-- pin: default apache/hive:4.0.1 sha256:5194161ef50b80875f937dff04936df047cbe894e9727ab8f0606349345b1bd3 rel/release-4.0.1

== TestQuerySelectsTheSchema default#1
> exec
| CREATE TABLE omni_query_schema.t (id INT) STORED AS PARQUET
< ok - -
> exec
| USE `default`
< ok - -
> exec
| USE `omni_query_schema`
< ok - -
> query
| SELECT count(*) AS n FROM t
< result
< column "n" "BIGINT"
< row
: int64 0
