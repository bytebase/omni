-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestExecFallbackNestingLimit_OceanBase default#1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
|# sha256:dcaab7b20db9fb65f43fd1312a627726fee74e6804ff808ce43dd34ac4407278 500077
< ok 0 0
> query
|# sha256:dcaab7b20db9fb65f43fd1312a627726fee74e6804ff808ce43dd34ac4407278 500077
< result
< column "0" "BIGINT"
< row
: int64 0
