-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestLoadDataCharsetFormsOceanBaseOracle default#1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| load data infile '/x.csv' into table t character set binary
< error
| 1017 HY000 File not exist
> exec ^
: string
| load data infile '/x.csv' into table t character set 'utf8mb4'
< error
| 1017 HY000 File not exist
> exec ^
: string
| load data infile '/x.csv' into table t character set utf8mb4
< error
| 1017 HY000 File not exist
> exec ^
: string
| load data infile '/x.csv' into table t character set 1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| load data infile '/x.csv' into table t charset binary
< error
| 1064 42000 …near 'charset binary' at line 1
> exec ^
: string
| load data infile '/x.csv' into table t charset 'utf8mb4'
< error
| 1064 42000 …near 'charset 'utf8mb4'' at line 1
