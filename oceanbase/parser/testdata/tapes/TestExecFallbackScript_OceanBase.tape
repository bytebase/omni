-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestExecFallbackScript_OceanBase default#1
> query
| select 1 /*! + 1 */; select 1 /*! from */
< result
< column "1 /*! + 1" "BIGINT"
< row
: int64 2
< result
< column "1" "BIGINT"
< row
: int64 1
> query
| select 1 /*! from */; select 1 /*! + 1 */
< result
< column "1" "BIGINT"
< row
: int64 1
< result
< column "1 /*! + 1" "BIGINT"
< row
: int64 2
> query
| select 1; /*! select from */; select 2
< result
< column "1" "BIGINT"
< row
: int64 1
< result
< column "2" "BIGINT"
< row
: int64 2
> query
| select 1 /*! from ; select 2 */
< result
< column "1" "BIGINT"
< row
: int64 1
