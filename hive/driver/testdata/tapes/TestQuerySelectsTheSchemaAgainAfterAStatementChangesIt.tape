-- pin: default apache/hive:4.0.1 sha256:5194161ef50b80875f937dff04936df047cbe894e9727ab8f0606349345b1bd3 rel/release-4.0.1

== TestQuerySelectsTheSchemaAgainAfterAStatementChangesIt default#1
> exec
| CREATE TABLE omni_schema_a.t (tag STRING) STORED AS PARQUET
< ok - -
> exec
| ALTER TABLE omni_schema_a.t SET TBLPROPERTIES('EXTERNAL'='FALSE')
< ok - -
> exec
| INSERT INTO omni_schema_a.t VALUES ('a')
< ok - -
> exec
| CREATE TABLE omni_schema_b.t (tag STRING) STORED AS PARQUET
< ok - -
> exec
| ALTER TABLE omni_schema_b.t SET TBLPROPERTIES('EXTERNAL'='FALSE')
< ok - -
> exec
| INSERT INTO omni_schema_b.t VALUES ('b')
< ok - -
> exec
| USE `omni_schema_a`
< ok - -
> query
| SELECT tag FROM t
< result
< column "tag" "STRING"
< row
: string
| a
> query
| USE omni_schema_b
< error
| -
| TStatus({StatusCode:ERROR_STATUS InfoMessages:[Server-side error; please check HS2 logs.] SqlState:<nil> ErrorCode:<nil> ErrorMessage:0x})
> query
| SELECT tag FROM t
< result
< column "tag" "STRING"
< row
: string
| b
> exec
| USE `omni_schema_a`
< ok - -
> query
| SELECT tag FROM t
< result
< column "tag" "STRING"
< row
: string
| a
