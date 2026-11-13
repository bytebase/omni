-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestColumnRefFormsOceanBaseOracle default#1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select t . key from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select key . key from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select d . t . key from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select d . key . key from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select . t . a from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select . t . key from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select . key . key from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select . t.a from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select t . key + 1 from t where t . key > 0 order by t . key
< error
| 1046 3D000 No database selected
> exec ^
: string
| update t set t . key = 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| update t set . t . a = 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| insert into t set key . key = 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| select not . not from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select key . a from t
< error
| 1064 42000 …near 'a from t' at line 1
> exec ^
: string
| select key . key . key from t
< error
| 1064 42000 …near '. key from t' at line 1
> exec ^
: string
| select t . key . a from t
< error
| 1064 42000 …near 'a from t' at line 1
> exec ^
: string
| select . key . a from t
< error
| 1064 42000 …near 'a from t' at line 1
> exec ^
: string
| select d . t . key . a from t
< error
| 1064 42000 …near '. a from t' at line 1
> exec ^
: string
| select . a from t
< error
| 1064 42000 …near 'from t' at line 1
> exec ^
: string
| select . t . * from t
< error
| 1064 42000 …near '* from t' at line 1
> exec ^
: string
| select t . select from t
< error
| 1064 42000 …near 'select from t' at line 1
> exec ^
: string
| select t . key (1) from t
< error
| 1064 42000 …near '(1) from t' at line 1
> exec ^
: string
| select . t . a ->> '$' from t
< error
| 1064 42000 …near '->> '$' from t' at line 1
> exec ^
: string
| select t . key -> '$' from t
< error
| 1064 42000 …near '-> '$' from t' at line 1
> exec ^
: string
| select t . a -> '$' from t
< error
| 1046 3D000 No database selected
