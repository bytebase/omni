-- pin: default apache/hive:4.0.1 sha256:5194161ef50b80875f937dff04936df047cbe894e9727ab8f0606349345b1bd3 rel/release-4.0.1

== TestSessionStartsFromTheConfiguredDatabase default#1
> exec
| CREATE TABLE omni_session_home.t (tag STRING) STORED AS PARQUET
< ok - -
> exec
| ALTER TABLE omni_session_home.t SET TBLPROPERTIES('EXTERNAL'='FALSE')
< ok - -
> exec
| INSERT INTO omni_session_home.t VALUES ('home')
< ok - -
> exec
| CREATE TABLE omni_session_other.t (tag STRING) STORED AS PARQUET
< ok - -
> exec
| ALTER TABLE omni_session_other.t SET TBLPROPERTIES('EXTERNAL'='FALSE')
< ok - -
> exec
| INSERT INTO omni_session_other.t VALUES ('other')
< ok - -
> exec
| USE `omni_session_home`
< ok - -
> query
| SELECT tag FROM t
< result
< column "tag" "STRING"
< row
: string
| home
> query
| USE omni_session_other
< error
| -
| TStatus({StatusCode:ERROR_STATUS InfoMessages:[Server-side error; please check HS2 logs.] SqlState:<nil> ErrorCode:<nil> ErrorMessage:0x})
> query
| SELECT tag FROM t
< result
< column "tag" "STRING"
< row
: string
| other
> exec
| USE `omni_session_home`
< ok - -
> query
| SELECT tag FROM t
< result
< column "tag" "STRING"
< row
: string
| home
> exec
| USE `default`
< ok - -
> query
| SELECT current_database() AS db
< result
< column "db" "STRING"
< row
: string
| default
> query
| USE omni_session_other
< error
| -
| TStatus({StatusCode:ERROR_STATUS InfoMessages:[Server-side error; please check HS2 logs.] SqlState:<nil> ErrorCode:<nil> ErrorMessage:0x})
> query
| SELECT current_database() AS db
< result
< column "db" "STRING"
< row
: string
| omni_session_other
> exec
| USE `default`
< ok - -
> query
| SELECT current_database() AS db
< result
< column "db" "STRING"
< row
: string
| default
