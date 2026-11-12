-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestLimitParamMarkerOracle default#1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| SELECT * FROM t LIMIT ?
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT * FROM t LIMIT ?, ?
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT * FROM t LIMIT ? OFFSET ?
< error
| 1046 3D000 No database selected
