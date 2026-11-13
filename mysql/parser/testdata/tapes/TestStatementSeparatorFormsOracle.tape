-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestStatementSeparatorFormsOracle default#1
> exec
| PREPARE omni_db_option_stmt FROM ?
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
| 1064 42000 …near 'commit' at line 1
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
| 1064 42000 …near '2' at line 1
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
| 1064 42000 …near 'select 1' at line 1
> exec ^
: string
| table t table u
< error
| 1064 42000 …near 'table u' at line 1
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
< error
| 1064 42000 …near 'select 2 */' at line 1
> exec ^
: string
| select 1 /*!80000 select 2 */
< error
| 1064 42000 …near 'select 2 */' at line 1
> exec ^
: string
| select 1 /*!100000 select 2 */
< error
| 1064 42000 …near '0 select 2 */' at line 1
> exec ^
: string
| create procedure p() begin select 1; end select 2
< error
| 1064 42000 …near 'select 2' at line 1
> exec ^
: string
| create procedure p() select 1 select 2
< error
| 1064 42000 …near 'select 2' at line 1
> exec ^
: string
| create trigger tr before insert on t for each row begin set new.a = 1; end insert into t values (1)
< error
| 1064 42000 …near 'insert into t values (1)' at line 1
> exec ^
: string
| create trigger tr before insert on t for each row set new.a = 1 insert into t values (1)
< error
| 1064 42000 …near 'insert into t values (1)' at line 1
> exec ^
: string
| create procedure p() begin select 1 select 2; end
< error
| 1064 42000 …near 'select 2; end' at line 1
> exec ^
: string
| create procedure p() begin select 1; select 2 end
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| create procedure p() begin if 1 then select 1 end if; end
< error
| 1064 42000 …near 'if; end' at line 1
> exec ^
: string
| create procedure p() begin declare a int declare b int; end
< error
| 1064 42000 …near 'declare b int; end' at line 1
> exec ^
: string
| create procedure p() begin begin select 1; end end
< error
| 1064 42000 …near 'end' at line 1
> exec ^
: string
| create procedure p() begin while 1 do select 1 end while; end
< error
| 1064 42000 …near 'while; end' at line 1
> exec ^
: string
| create procedure p() begin select 1; select 2; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin select 1; select 2; end
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() begin end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin end
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() begin if 1 then select 1; end if; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin if 1 then select 1; end if; end
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() begin declare a int; declare b int; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin declare a int; declare b int; end
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() begin begin select 1; end; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin begin select 1; end; end
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
< ok 0 0
> exec ^
: string
| begin
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| begin
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| begin work
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| begin work
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create table t (a int) select 1
< ok 0 0
> exec ^
: string
| create table t (a int) (select 1)
< ok 0 0
> exec ^
: string
| create table t (a int) (select 1) union (select 2)
< ok 0 0
> exec ^
: string
| create table t (a int) (select 1) order by 1
< ok 0 0
> exec ^
: string
| create table t (a int) ((select 1))
< ok 0 0
> exec ^
: string
| create table t (a int) (with c as (select 1) select * from c)
< ok 0 0
> exec ^
: string
| create table t (a int) with c as (select 1) select * from c
< ok 0 0
> exec ^
: string
| create table t (a int) with c as (select 1) table c
< ok 0 0
> exec ^
: string
| create table t (a int) table u
< error
| 1146 42S02 Table 'omni_t_mysql_parser.u' doesn't exist
> exec ^
: string
| create table t (a int) table u
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create table t (a int) (table u)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.u' doesn't exist
> exec ^
: string
| create table t (a int) (table u)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create table t (a int) table u order by a limit 1
< error
| 1146 42S02 Table 'omni_t_mysql_parser.u' doesn't exist
> exec ^
: string
| create table t (a int) table u order by a limit 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create table t (a int) table u union select 1
< error
| 1146 42S02 Table 'omni_t_mysql_parser.u' doesn't exist
> exec ^
: string
| create table t (a int) table u union select 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create table t (a int) values row(1)
< ok 0 0
> exec ^
: string
| create table t (a int) (values row(1))
< ok 0 0
> exec ^
: string
| create table t (a int) values row(1) order by 1
< ok 0 0
> exec ^
: string
| create table t (a int) values row(1) union table u
< error
| 1146 42S02 Table 'omni_t_mysql_parser.u' doesn't exist
> exec ^
: string
| create table t (a int) values row(1) union table u
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create table t (a int) values (1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| create table t select 1
< ok 0 0
> exec ^
: string
| create table t table u
< error
| 1146 42S02 Table 'omni_t_mysql_parser.u' doesn't exist
> exec ^
: string
| create table t table u
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create table t values row(1)
< ok 0 0
> exec ^
: string
| create table t values row(1), row(2)
< ok 0 0
> exec ^
: string
| create table t values (1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| create table t with c as (select 1) select * from c
< ok 0 0
> exec ^
: string
| create table t (select 1)
< ok 0 0
> exec ^
: string
| create table t (select 1) union (select 2)
< ok 0 0
> exec ^
: string
| create table t (select 1) order by 1
< ok 0 0
> exec ^
: string
| create table t ((select 1))
< ok 0 0
> exec ^
: string
| create table t ((select 1)) union (select 2)
< ok 0 0
> exec ^
: string
| create table t (with c as (select 1) select * from c)
< ok 0 0
> exec ^
: string
| create table t (table u)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.u' doesn't exist
> exec ^
: string
| create table t (table u)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create table t (values row(1))
< ok 0 0
> exec ^
: string
| create table t (a int) engine=innodb (select 1)
< ok 0 0
> exec ^
: string
| create table t (a int) engine=innodb table u
< error
| 1146 42S02 Table 'omni_t_mysql_parser.u' doesn't exist
> exec ^
: string
| create table t (a int) engine=innodb table u
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create table t (a int) engine=innodb values row(1)
< ok 0 0
> exec ^
: string
| create table t (a int) engine=innodb with c as (select 1) select * from c
< ok 0 0
> exec ^
: string
| create table t engine=innodb table u
< error
| 1146 42S02 Table 'omni_t_mysql_parser.u' doesn't exist
> exec ^
: string
| create table t engine=innodb table u
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create table t (a int) partition by hash(a) partitions 2 table u
< error
| 1146 42S02 Table 'omni_t_mysql_parser.u' doesn't exist
> exec ^
: string
| create table t (a int) partition by hash(a) partitions 2 table u
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create table t (a int) partition by hash(a) partitions 2 values row(1)
< ok 0 0
> exec ^
: string
| create table t (a int) partition by hash(a) partitions 2 with c as (select 1) select * from c
< ok 0 0
> exec ^
: string
| create table t (a int) partition by hash(a) (partition p0) (select 1)
< ok 0 0
> exec ^
: string
| create table t (a int) partition by hash(a) partitions 2 (select 1)
< error
| 1064 42000 …near 'select 1)' at line 1
> exec ^
: string
| create table t (a int) engine=innodb comment 'x'
< ok 0 0
> exec ^
: string
| create table t (a int) engine=innodb, comment 'x'
< ok 0 0
> exec ^
: string
| create table t (a int) engine=innodb,
< error
| 1064 42000 …near '' at line 1
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
< ok 0 0
> exec ^
: string
| create table t (a int) ignore table u
< error
| 1146 42S02 Table 'omni_t_mysql_parser.u' doesn't exist
> exec ^
: string
| create table t (a int) ignore table u
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create table t (a int) ignore with c as (select 1) select * from c
< ok 0 0
> exec ^
: string
| create table t (a int) replace table u
< error
| 1146 42S02 Table 'omni_t_mysql_parser.u' doesn't exist
> exec ^
: string
| create table t (a int) replace table u
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create table t (a int) replace values row(1)
< ok 0 0
> exec ^
: string
| create table t (a int) partition by hash(a) partitions 2 ignore (select 1)
< ok 0 0
> exec ^
: string
| create table t ignore select 1
< ok 0 0
> exec ^
: string
| create table t (a int) ignore as select 1
< ok 0 0
> exec ^
: string
| create table t (a int) ignore as table u
< error
| 1146 42S02 Table 'omni_t_mysql_parser.u' doesn't exist
> exec ^
: string
| create table t (a int) ignore as table u
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create table t (a int) as ignore select 1
< error
| 1064 42000 …near 'ignore select 1' at line 1
> exec ^
: string
| create table t (a int) ignore
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| create table t (a int) replace
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| create table t ignore
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| create table t as table u
< error
| 1146 42S02 Table 'omni_t_mysql_parser.u' doesn't exist
> exec ^
: string
| create table t as table u
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create table t as values row(1)
< ok 0 0
> exec ^
: string
| create table t as (select 1)
< ok 0 0
> exec ^
: string
| create table t (a int) as (select 1)
< ok 0 0
> exec ^
: string
| create table t (a int) as table u
< error
| 1146 42S02 Table 'omni_t_mysql_parser.u' doesn't exist
> exec ^
: string
| create table t (a int) as table u
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create table t (a int) as with c as (select 1) select * from c
< ok 0 0
> exec ^
: string
| create table t (a int) as
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| create table t (a int) ignore as
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| create table t as
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| create table t (a int) select 1 for update
< ok 0 0
> exec ^
: string
| create table t (a int) (select 1) for update
< ok 0 0
> exec ^
: string
| create table t (a int) table u for update
< error
| 1146 42S02 Table 'omni_t_mysql_parser.u' doesn't exist
> exec ^
: string
| create table t (a int) table u for update
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create table t (a int) (select 1) union (select 2) for update
< ok 0 0
> exec ^
: string
| create table t (a int) select 1 union select 2 for update
< ok 0 0
> exec ^
: string
| create table t (a int) (select 1) order by 1 for update
< ok 0 0
> exec ^
: string
| create table t as (select 1) for update
< ok 0 0
> exec ^
: string
| create table t (a int) as (select 1) for update
< ok 0 0
> exec ^
: string
| create table t (a int) as select 1 select 2
< error
| 1064 42000 …near 'select 2' at line 1
> exec ^
: string
| create table t (a int) table u table v
< error
| 1064 42000 …near 'table v' at line 1
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
< ok 0 0
