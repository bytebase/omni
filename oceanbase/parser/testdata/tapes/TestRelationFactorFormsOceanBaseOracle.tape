-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestRelationFactorFormsOceanBaseOracle default#1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select * from .t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from . t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from . table
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from .table
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from d . table
< error
| 1146 42S02 Table 'd.table' doesn't exist
> exec ^
: string
| select * from d . key as k
< error
| 1146 42S02 Table 'd.key' doesn't exist
> exec ^
: string
| select * from .t join . u on true
< error
| 1046 3D000 No database selected
> exec ^
: string
| explain . linear
< error
| 1046 3D000 No database selected
> exec ^
: string
| desc . linear a
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table .t (a int)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table d . key (a int)
< error
| 1049 42000 Unknown database 'd'
> exec ^
: string
| create table u like .t
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table c (a int, foreign key (a) references .t (a))
< error
| 1046 3D000 No database selected
> exec ^
: string
| drop table .t, d . index
< error
| 1046 3D000 No database selected
> exec ^
: string
| truncate table .t
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table .t add column b int
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t rename to .u
< error
| 1046 3D000 No database selected
> exec ^
: string
| rename table .t to .u
< error
| 1046 3D000 No database selected
> exec ^
: string
| insert into .t values (1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| update .t set a = 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| delete from .t
< error
| 1046 3D000 No database selected
> exec ^
: string
| create index i on .t (a)
< error
| 1046 3D000 No database selected
> exec ^
: string
| drop index i on .t
< error
| 1046 3D000 No database selected
> exec ^
: string
| create view .v as select 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| show create table .t
< error
| 1046 3D000 No database selected
> exec ^
: string
| show columns from .t
< error
| 1046 3D000 No database selected
> exec ^
: string
| show index from .t
< error
| 1046 3D000 No database selected
> exec ^
: string
| optimize table .t
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t drop tablegroup
< error
| 1046 3D000 No database selected
> exec ^
: string
| explain . select
< error
| 1064 42000 …near 'select' at line 1
> exec ^
: string
| explain . all
< error
| 1064 42000 …near 'all' at line 1
> exec ^
: string
| select * from . from
< error
| 1064 42000 …near 'from' at line 1
> exec ^
: string
| select * from d . distinct
< error
| 1064 42000 …near 'distinct' at line 1
> exec ^
: string
| select * from table . t
< error
| 1064 42000 …near '. t' at line 1
> exec ^
: string
| explain linear
< error
| 1064 42000 …near 'linear' at line 1
> exec ^
: string
| select * from ..t
< error
| 1064 42000 …near '.t' at line 1
> exec ^
: string
| select * from .
< error
| 1064 42000 …near '.' at line 1
> exec ^
: string
| select * from . 1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| create index .i on t (a)
< error
| 1064 42000 …near '.i on t (a)' at line 1
> exec ^
: string
| select * from tablegroup
< error
| 1064 42000 …near 'tablegroup' at line 1
> exec ^
: string
| select * from . tablegroup
< error
| 1064 42000 …near 'tablegroup' at line 1
> exec ^
: string
| select read_consistency from t
< error
| 1064 42000 …near 'read_consistency from t' at line 1
> exec ^
: string
| select no_parallel from t
< error
| 1064 42000 …near 'no_parallel from t' at line 1
> exec ^
: string
| select no_rewrite from t
< error
| 1064 42000 …near 'no_rewrite from t' at line 1
