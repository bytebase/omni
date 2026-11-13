-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestTableNameClausesOceanBaseOracle default#1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select * from t sample (1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from t sample block (1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from t sample hybrid all (1.5) seed (1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from t sample base (1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from t sample incr (1) a
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from t sample (1) seed (1) as a use key (i)
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from t partition (p) sample (1) seed (1) as of snapshot 1 as a use key (i)
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from t as of snapshot 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from t as of snapshot 1 + 1 a
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from t as of proctime()
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from t partition (p) as of snapshot 1 a
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from t partition (a = 's')
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from t partition (a = 's', b = 1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from t use key ()
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from t use key (primary)
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from t use key (`primary_zone`)
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from t force key (i) ignore index for order by (i, j)
< error
| 1046 3D000 No database selected
> exec ^
: string
| update t sample (1) set a = 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| delete from t partition (p) a where a.id = 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| delete from t sample (1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| delete from t as of snapshot 1
< error
| 32491 HY000 snapshot expression not allowed here
> exec ^
: string
| delete from t use key (i)
< error
| 1046 3D000 No database selected
> exec ^
: string
| delete from t a use key (i) where a.id = 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from t sample
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from t sample use key (i)
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from t seed
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from t as of
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from t seed (1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| select * from t a sample (1)
< error
| 1064 42000 …near 'sample (1)' at line 1
> exec ^
: string
| select * from t a partition (p)
< error
| 1064 42000 …near 'partition (p)' at line 1
> exec ^
: string
| select * from t sample (1) partition (p)
< error
| 1064 42000 …near 'partition (p)' at line 1
> exec ^
: string
| select * from t as of snapshot 1 sample (1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| select * from t use key (i) a
< error
| 1064 42000 …near 'a' at line 1
> exec ^
: string
| delete from t a partition (p)
< error
| 1064 42000 …near 'partition (p)' at line 1
> exec ^
: string
| delete from t as a partition (p)
< error
| 1064 42000 …near 'partition (p)' at line 1
> exec ^
: string
| select * from t as of snapshot 1 use key (i)
< error
| 1064 42000 …near 'use key (i)' at line 1
> exec ^
: string
| select * from t as of snapshot 1 a use key (i)
< error
| 1064 42000 …near 'use key (i)' at line 1
> exec ^
: string
| select * from t sample (1) as of snapshot 1 use key (i)
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from t sample (1e1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from t sample ('1')
< error
| 1064 42000 …near ''1')' at line 1
> exec ^
: string
| select * from t sample (1) seed (1.5)
< error
| 1064 42000 …near '1.5)' at line 1
> exec ^
: string
| select * from t as of snapshot 1 = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| select * from t partition (a = 1 + 1)
< error
| 1064 42000 …near '+ 1)' at line 1
> exec ^
: string
| select * from t use key (primary_zone)
< error
| 1064 42000 …near 'primary_zone)' at line 1
> exec ^
: string
| select * from t use key (status)
< error
| 1064 42000 …near 'status)' at line 1
> exec ^
: string
| select * from t force key ()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select * from t 'a'
< error
| 1064 42000 …near ''a'' at line 1
> exec ^
: string
| select * from t as 'a'
< error
| 1064 42000 …near ''a'' at line 1
> exec ^
: string
| delete from t a using t as a
< error
| 1064 42000 …near 'using t as a' at line 1
> exec ^
: string
| delete from t partition (p) using t
< error
| 1064 42000 …near 'using t' at line 1
> exec ^
: string
| delete from t use key (i) using t
< error
| 1064 42000 …near 'using t' at line 1
> exec ^
: string
| delete from t using t
< error
| 1046 3D000 No database selected
