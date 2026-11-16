-- pin: default apache/hive:4.0.1 sha256:5194161ef50b80875f937dff04936df047cbe894e9727ab8f0606349345b1bd3 rel/release-4.0.1

== TestExecuteStartsFromTheConfiguredDatabase default#1
> exec
| USE `omni_execute_home`
< ok - -
> query
| USE omni_execute_other
< error
| -
| TStatus({StatusCode:ERROR_STATUS InfoMessages:[Server-side error; please check HS2 logs.] SqlState:<nil> ErrorCode:<nil> ErrorMessage:0x})
> exec
| USE `omni_execute_home`
< ok - -
> exec
| CREATE TABLE exec_home (id INT) STORED AS PARQUET
< ok - -
> exec
| USE `omni_execute_home`
< ok - -
> exec
| USE `omni_execute_other`
< ok - -
> query
| SHOW TABLES
< result
< column "tab_name" "STRING"
> exec
| USE `omni_execute_home`
< ok - -
> query
| SHOW TABLES
< result
< column "tab_name" "STRING"
< row
: string
| exec_home
