-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestStatementSeparatorFormsOceanBaseOracle default#1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select 1 select 2
< error
| 1064 42000 …near 'select 2' at line 1
> exec ^
: string
| select 1 update t set a=1
< error
| 1064 42000 …near 'update t set a=1' at line 1
> exec ^
: string
| create table t (a int) create table u (a int)
< error
| 1064 42000 …near 'create table u (a int)' at line 1
> exec ^
: string
| begin commit
< error
| 1064 42000 …near 'begin' at line 1
> exec ^
: string
| commit begin
< error
| 1064 42000 …near 'begin' at line 1
> exec ^
: string
| rollback select 1
< error
| 1064 42000 …near 'select 1' at line 1
> exec ^
: string
| set @a = 1 select 1
< error
| 1064 42000 …near 'select 1' at line 1
> exec ^
: string
| set @a = 1 set @b = 2
< error
| 1064 42000 …near 'set @b = 2' at line 1
> exec ^
: string
| select 1 from dual where 1 = 1 select 2
< error
| 1064 42000 …near 'select 2' at line 1
> exec ^
: string
| select 1 union select 2 select 3
< error
| 1064 42000 …near 'select 3' at line 1
> exec ^
: string
| select 1 for update select 2
< error
| 1064 42000 …near 'select 2' at line 1
> exec ^
: string
| with c as (select 1) select * from c select 2
< error
| 1064 42000 …near 'select 2' at line 1
> exec ^
: string
| delete from t where a = 1 delete from t
< error
| 1064 42000 …near 'delete from t' at line 1
> exec ^
: string
| insert into t values (1) insert into t values (2)
< error
| 1064 42000 …near 'insert into t values (2)' at line 1
> exec ^
: string
| update t set a = 1 update t set a = 2
< error
| 1064 42000 …near 'update t set a = 2' at line 1
> exec ^
: string
| show tables show databases
< error
| 1064 42000 …near 'show databases' at line 1
> exec ^
: string
| drop table t drop table u
< error
| 1064 42000 …near 'drop table u' at line 1
> exec ^
: string
| alter table t add column b int alter table t drop column b
< error
| 1064 42000 …near 'alter table t drop column b' at line 1
> exec ^
: string
| create view v as select 1 select 2
< error
| 1064 42000 …near 'select 2' at line 1
> exec ^
: string
| explain select 1 select 2
< error
| 1064 42000 …near 'select 2' at line 1
> exec ^
: string
| do 1 do 2
< error
| 1064 42000 …near 'do' at line 1
> exec ^
: string
| handler t open handler t close
< error
| 1064 42000 …near 't close' at line 1
> exec ^
: string
| prepare s from 'select 1' execute s
< error
| 1064 42000 …near 'execute s' at line 1
> exec ^
: string
| lock tables t read unlock tables
< error
| 1064 42000 …near 'unlock tables' at line 1
> exec ^
: string
| xa start 'x' xa end 'x'
< error
| 1064 42000 …near 'xa end 'x'' at line 1
> exec ^
: string
| signal sqlstate '45000' select 1
< error
| 1064 42000 …near 'select' at line 1
> exec ^
: string
| table t table u
< error
| 1064 42000 …near 'table t table u' at line 1
> exec ^
: string
| values row(1) values row(2)
< error
| 1064 42000 …near 'values row(2)' at line 1
> exec ^
: string
| select 1 = any select 1
< error
| 1064 42000 …near 'select 1' at line 1
> exec ^
: string
| select 1 -- c
| select 2
< error
| 1064 42000 …near 'select 2' at line 2
> exec ^
: string
| select 1 # c
| select 2
< error
| 1064 42000 …near 'select 2' at line 2
> exec ^
: string
| select 1 /* c */ select 2
< error
| 1064 42000 …near 'select 2' at line 1
> exec ^
: string
| select 1 /*! select 2 */
< ok 0 0
> exec ^
: string
| select 1 /*!80000 select 2 */
< ok 0 0
> exec ^
: string
| select 1 /*!100000 select 2 */
< ok 0 0
> exec ^
: string
| create procedure p() begin select 1; end select 2
< error
| 1064 42000 …near 'select' at line 1
> exec ^
: string
| create procedure p() select 1 select 2
< error
| 1064 42000 …near 'select 1 select 2' at line 1
> exec ^
: string
| create trigger tr before insert on t for each row begin set new.a = 1; end insert into t values (1)
< error
| 1064 42000 …near 'insert' at line 1
> exec ^
: string
| create trigger tr before insert on t for each row set new.a = 1 insert into t values (1)
< error
| 1064 42000 …near 'set new.a = 1 insert into t values (1)' at line 1
> exec ^
: string
| create procedure p() begin select 1 select 2; end
< error
| 1064 42000 …near 'select 1 select 2' at line 1
> exec ^
: string
| create procedure p() begin select 1; select 2 end
< error
| 1064 42000 …near 'end' at line 1
> exec ^
: string
| create procedure p() begin if 1 then select 1 end if; end
< error
| 1064 42000 …near 'select 1 end if' at line 1
> exec ^
: string
| create procedure p() begin declare a int declare b int; end
< error
| 1064 42000 …near 'declare' at line 1
> exec ^
: string
| create procedure p() begin begin select 1; end end
< error
| 1064 42000 …near 'end' at line 1
> exec ^
: string
| create procedure p() begin while 1 do select 1 end while; end
< error
| 1064 42000 …near 'select 1 end while' at line 1
> exec ^
: string
| create procedure p() begin select 1; select 2; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin if 1 then select 1; end if; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin declare a int; declare b int; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin begin select 1; end; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| select 1;
< ok 0 0
> exec ^
: string
| select 1 ;
< ok 0 0
> exec ^
: string
| select 1; -- c
< error
| 1235 0A000 prepare multi-statement query not supported
> exec ^
: string
| begin
< ok 0 0
> exec ^
: string
| begin work
< ok 0 0
> exec ^
: string
| create table t (a int) select 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (a int) (select 1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (a int) (select 1) union (select 2)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (a int) (select 1) order by 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (a int) ((select 1))
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (a int) (with c as (select 1) select * from c)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (a int) with c as (select 1) select * from c
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (a int) with c as (select 1) table c
< error
| 1064 42000 …near 'table c' at line 1
> exec ^
: string
| create table t (a int) table u
< error
| 1064 42000 …near 'table u' at line 1
> exec ^
: string
| create table t (a int) (table u)
< error
| 1064 42000 …near 'table u)' at line 1
> exec ^
: string
| create table t (a int) table u order by a limit 1
< error
| 1064 42000 …near 'table u order by a limit 1' at line 1
> exec ^
: string
| create table t (a int) table u union select 1
< error
| 1064 42000 …near 'table u union select 1' at line 1
> exec ^
: string
| create table t (a int) values row(1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (a int) (values row(1))
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (a int) values row(1) order by 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (a int) values row(1) union table u
< error
| 1064 42000 …near 'table u' at line 1
> exec ^
: string
| create table t (a int) values (1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| create table t select 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t table u
< error
| 1064 42000 …near 'table u' at line 1
> exec ^
: string
| create table t values row(1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t values row(1), row(2)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t values (1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| create table t with c as (select 1) select * from c
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (select 1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (select 1) union (select 2)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (select 1) order by 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t ((select 1))
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t ((select 1)) union (select 2)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (with c as (select 1) select * from c)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (table u)
< error
| 1064 42000 …near 'table u)' at line 1
> exec ^
: string
| create table t (values row(1))
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (a int) engine=innodb (select 1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (a int) engine=innodb table u
< error
| 1064 42000 …near 'table u' at line 1
> exec ^
: string
| create table t (a int) engine=innodb values row(1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (a int) engine=innodb with c as (select 1) select * from c
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t engine=innodb table u
< error
| 1064 42000 …near 'table u' at line 1
> exec ^
: string
| create table t (a int) partition by hash(a) partitions 2 table u
< error
| 1064 42000 …near 'table u' at line 1
> exec ^
: string
| create table t (a int) partition by hash(a) partitions 2 values row(1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (a int) partition by hash(a) partitions 2 with c as (select 1) select * from c
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (a int) partition by hash(a) (partition p0) (select 1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (a int) partition by hash(a) partitions 2 (select 1)
< error
| 1064 42000 …near 'select 1)' at line 1
> exec ^
: string
| create table t (a int) engine=innodb comment 'x'
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (a int) engine=innodb, comment 'x'
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (a int) engine=innodb,
< error
| 1064 42000 …near ',' at line 1
> exec ^
: string
| create table t (a int) engine=innodb, select 1
< error
| 1064 42000 …near 'select 1' at line 1
> exec ^
: string
| create table t (a int) engine=innodb, (select 1)
< error
| 1064 42000 …near '(select 1)' at line 1
> exec ^
: string
| create table t (a int) engine=innodb, partition by hash(a)
< error
| 1064 42000 …near 'partition by hash(a)' at line 1
> exec ^
: string
| create table t (a int) engine=innodb,, comment 'x'
< error
| 1064 42000 …near ', comment 'x'' at line 1
> exec ^
: string
| create table t engine=innodb, select 1
< error
| 1064 42000 …near 'select 1' at line 1
> exec ^
: string
| create table t (a int) ignore (select 1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (a int) ignore table u
< error
| 1064 42000 …near 'table u' at line 1
> exec ^
: string
| create table t (a int) ignore with c as (select 1) select * from c
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (a int) replace table u
< error
| 1064 42000 …near 'table u' at line 1
> exec ^
: string
| create table t (a int) replace values row(1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (a int) partition by hash(a) partitions 2 ignore (select 1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t ignore select 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (a int) ignore as select 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (a int) ignore as table u
< error
| 1064 42000 …near 'table u' at line 1
> exec ^
: string
| create table t (a int) as ignore select 1
< error
| 1064 42000 …near 'ignore select 1' at line 1
> exec ^
: string
| create table t (a int) ignore
< error
| 1064 42000 …near 'ignore' at line 1
> exec ^
: string
| create table t (a int) replace
< error
| 1064 42000 …near 'replace' at line 1
> exec ^
: string
| create table t ignore
< error
| 1064 42000 …near 'ignore' at line 1
> exec ^
: string
| create table t as table u
< error
| 1064 42000 …near 'table u' at line 1
> exec ^
: string
| create table t as values row(1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t as (select 1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (a int) as (select 1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (a int) as table u
< error
| 1064 42000 …near 'table u' at line 1
> exec ^
: string
| create table t (a int) as with c as (select 1) select * from c
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (a int) as
< error
| 1064 42000 …near 'as' at line 1
> exec ^
: string
| create table t (a int) ignore as
< error
| 1064 42000 …near 'as' at line 1
> exec ^
: string
| create table t as
< error
| 1064 42000 …near 'as' at line 1
> exec ^
: string
| create table t (a int) select 1 for update
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (a int) (select 1) for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| create table t (a int) table u for update
< error
| 1064 42000 …near 'table u for update' at line 1
> exec ^
: string
| create table t (a int) (select 1) union (select 2) for update
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (a int) select 1 union select 2 for update
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (a int) (select 1) order by 1 for update
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t as (select 1) for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| create table t (a int) as (select 1) for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| create table t (a int) as select 1 select 2
< error
| 1064 42000 …near 'select 2' at line 1
> exec ^
: string
| create table t (a int) table u table v
< error
| 1064 42000 …near 'table u table v' at line 1
> exec ^
: string
| create table t (a int) values row(1) select 1
< error
| 1064 42000 …near 'select 1' at line 1
> exec ^
: string
| create table t (a int) (select 1) (select 2)
< error
| 1064 42000 …near '(select 2)' at line 1
> exec ^
: string
| create table t (a int) select 1 into @x
< error
| 1046 3D000 No database selected
