-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_OceanBaseCheckExpr/ALTER_TABLE default#1
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, CONSTRAINT k CHECK (c > 0))
< ok 0 0
> exec
| ALTER TABLE a ADD CHECK (c)
< error
| 3812 HY000 An expression of non-boolean type specified to a check constraint 'a_OBCHECK_0000000000000001'
> exec
| ALTER TABLE a ADD CHECK ('s')
< error
| 3812 HY000 An expression of non-boolean type specified to a check constraint 'a_OBCHECK_0000000000000002'
> exec
| ALTER TABLE a ADD CHECK (c + 1)
< error
| 3812 HY000 An expression of non-boolean type specified to a check constraint 'a_OBCHECK_0000000000000003'
> exec
| ALTER TABLE a ADD CHECK (c > zz)
< error
| 3820 HY000 Check constraint 'a_OBCHECK_0000000000000004' refers to non-existing column 'zz'.
> exec
| ALTER TABLE a ADD CONSTRAINT k CHECK (zz)
< error
| 3820 HY000 Check constraint 'k' refers to non-existing column 'zz'.
> exec
| ALTER TABLE a ADD CONSTRAINT k CHECK (c)
< error
| 3812 HY000 An expression of non-boolean type specified to a check constraint 'k'
> exec
| ALTER TABLE a ADD CONSTRAINT k CHECK (c > 1)
< error
| 3822 HY000 Duplicate check constraint name 'k'.
> exec
| ALTER TABLE a ADD CHECK (1)
< ok 0 0

== TestOracle_OceanBaseCheckExpr/CREATE_TABLE default#1
> exec
| CREATE TABLE t2 (b INT PRIMARY KEY, CHECK (b))
< error
| 3812 HY000 An expression of non-boolean type specified to a check constraint 't2_OBCHECK_0000000000000005'
> exec
| CREATE TABLE t2 (b INT PRIMARY KEY, CHECK (zz > 0))
< error
| 3820 HY000 Check constraint 't2_OBCHECK_0000000000000006' refers to non-existing column 'zz'.
