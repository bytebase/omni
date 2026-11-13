-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestStatementSeparatorFormsMatchTiDB default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| select 1 select 2
< error
| 1064 42000 …line 1 column 15 near "select 2" 
> exec ^
: string
| select 1 update t set a=1
< error
| 1064 42000 …line 1 column 15 near "update t set a=1" 
> exec ^
: string
| create table t (a int) create table u (a int)
< error
| 1064 42000 …line 1 column 29 near "create table u (a int)" 
> exec ^
: string
| begin commit
< error
| 1064 42000 …line 1 column 12 near "commit" 
> exec ^
: string
| commit begin
< error
| 1064 42000 …line 1 column 12 near "begin" 
> exec ^
: string
| rollback select 1
< error
| 1064 42000 …line 1 column 15 near "select 1" 
> exec ^
: string
| set @a = 1 select 1
< error
| 1064 42000 …line 1 column 17 near "select 1" 
> exec ^
: string
| set @a = 1 set @b = 2
< error
| 1064 42000 …line 1 column 14 near "set @b = 2" 
> exec ^
: string
| select 1 from dual where 1 = 1 select 2
< error
| 1064 42000 …line 1 column 37 near "select 2" 
> exec ^
: string
| select 1 union select 2 select 3
< error
| 1064 42000 …line 1 column 30 near "select 3" 
> exec ^
: string
| select 1 for update select 2
< error
| 1064 42000 …line 1 column 26 near "select 2" 
> exec ^
: string
| with c as (select 1) select * from c select 2
< error
| 1064 42000 …line 1 column 43 near "select 2" 
> exec ^
: string
| delete from t where a = 1 delete from t
< error
| 1064 42000 …line 1 column 32 near "delete from t" 
> exec ^
: string
| insert into t values (1) insert into t values (2)
< error
| 1064 42000 …line 1 column 31 near "insert into t values (2)" 
> exec ^
: string
| update t set a = 1 update t set a = 2
< error
| 1064 42000 …line 1 column 25 near "update t set a = 2" 
> exec ^
: string
| show tables show databases
< error
| 1064 42000 …line 1 column 16 near "show databases" 
> exec ^
: string
| drop table t drop table u
< error
| 1064 42000 …line 1 column 17 near "drop table u" 
> exec ^
: string
| alter table t add column b int alter table t drop column b
< error
| 1064 42000 …line 1 column 36 near "alter table t drop column b" 
> exec ^
: string
| create view v as select 1 select 2
< error
| 1064 42000 …line 1 column 32 near "select 2" 
> exec ^
: string
| explain select 1 select 2
< error
| 1064 42000 …line 1 column 23 near "select 2" 
> exec ^
: string
| do 1 do 2
< error
| 1064 42000 …line 1 column 7 near "do 2" 
> exec ^
: string
| handler t open handler t close
< error
| 1064 42000 …line 1 column 7 near "handler t open handler t close" 
> exec ^
: string
| prepare s from 'select 1' execute s
< error
| 1064 42000 …line 1 column 33 near "execute s" 
> exec ^
: string
| lock tables t read unlock tables
< error
| 1064 42000 …line 1 column 25 near "unlock tables" 
> exec ^
: string
| xa start 'x' xa end 'x'
< error
| 1064 42000 …line 1 column 2 near "xa start 'x' xa end 'x'" 
> exec ^
: string
| signal sqlstate '45000' select 1
< error
| 1064 42000 …line 1 column 6 near "signal sqlstate '45000' select 1" 
> exec ^
: string
| table t table u
< error
| 1064 42000 …line 1 column 13 near "table u" 
> exec ^
: string
| values row(1) values row(2)
< error
| 1064 42000 …line 1 column 20 near "values row(2)" 
> exec ^
: string
| select 1 = any select 1
< error
| 1064 42000 …line 1 column 21 near "select 1" 
> exec ^
: string
| select 1 -- c
| select 2
< error
| 1064 42000 …line 2 column 7 near "select 2" 
> exec ^
: string
| select 1 # c
| select 2
< error
| 1064 42000 …line 2 column 7 near "select 2" 
> exec ^
: string
| select 1 /* c */ select 2
< error
| 1064 42000 …line 1 column 23 near "select 2" 
> exec ^
: string
| select 1 /*! select 2 */
< error
| 1064 42000 …line 1 column 19 near "select 2 */" 
> exec ^
: string
| select 1 /*!80000 select 2 */
< error
| 1064 42000 …line 1 column 24 near "select 2 */" 
> exec ^
: string
| select 1 /*!100000 select 2 */
< error
| 1064 42000 …line 1 column 18 near "0 select 2 */" 
> exec ^
: string
| create procedure p() begin select 1; end select 2
< error
| 1064 42000 …line 1 column 47 near "select 2" 
> exec ^
: string
| create procedure p() select 1 select 2
< error
| 1064 42000 …line 1 column 36 near "select 2" 
> exec ^
: string
| create trigger tr before insert on t for each row begin set new.a = 1; end insert into t values (1)
< error
| 1064 42000 …line 1 column 14 near "trigger tr before insert on t for each row begin set new.a = 1; end insert into t values (1)" 
> exec ^
: string
| create trigger tr before insert on t for each row set new.a = 1 insert into t values (1)
< error
| 1064 42000 …line 1 column 14 near "trigger tr before insert on t for each row set new.a = 1 insert into t values (1)" 
> exec ^
: string
| create procedure p() begin select 1 select 2; end
< error
| 1064 42000 …line 1 column 42 near "select 2; end" 
> exec ^
: string
| create procedure p() begin select 1; select 2 end
< error
| 1064 42000 …line 1 column 49 near "" 
> exec ^
: string
| create procedure p() begin if 1 then select 1 end if; end
< error
| 1064 42000 …line 1 column 52 near "if; end" 
> exec ^
: string
| create procedure p() begin declare a int declare b int; end
< error
| 1064 42000 …line 1 column 48 near "declare b int; end" 
> exec ^
: string
| create procedure p() begin begin select 1; end end
< error
| 1064 42000 …line 1 column 50 near "end" 
> exec ^
: string
| create procedure p() begin while 1 do select 1 end while; end
< error
| 1064 42000 …line 1 column 56 near "while; end" 
> exec ^
: string
| create procedure p() begin select 1; select 2; end
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| create procedure p() begin select 1; select 2; end
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create procedure p() begin end
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| create procedure p() begin end
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create procedure p() begin if 1 then select 1; end if; end
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| create procedure p() begin if 1 then select 1; end if; end
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create procedure p() begin declare a int; declare b int; end
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| create procedure p() begin declare a int; declare b int; end
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create procedure p() begin begin select 1; end; end
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| create procedure p() begin begin select 1; end; end
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
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
< ok 0 0
> exec ^
: string
| begin work
< error
| 1064 42000 …line 1 column 10 near "work" 
> exec ^
: string
| create table t (a int) select 1
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) select 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) (select 1)
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) (select 1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) (select 1) union (select 2)
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) (select 1) union (select 2)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) (select 1) order by 1
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) (select 1) order by 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) ((select 1))
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) ((select 1))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) (with c as (select 1) select * from c)
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) (with c as (select 1) select * from c)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) with c as (select 1) select * from c
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) with c as (select 1) select * from c
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) with c as (select 1) table c
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) with c as (select 1) table c
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) table u
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) table u
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) (table u)
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) (table u)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) table u order by a limit 1
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) table u order by a limit 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) table u union select 1
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) table u union select 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) values row(1)
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) values row(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) (values row(1))
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) (values row(1))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) values row(1) order by 1
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) values row(1) order by 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) values row(1) union table u
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) values row(1) union table u
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) values (1)
< error
| 1064 42000 …line 1 column 31 near "(1)" 
> exec ^
: string
| create table t select 1
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t select 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t table u
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t table u
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t values row(1)
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t values row(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t values row(1), row(2)
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t values row(1), row(2)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t values (1)
< error
| 1064 42000 …line 1 column 23 near "(1)" 
> exec ^
: string
| create table t with c as (select 1) select * from c
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t with c as (select 1) select * from c
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (select 1)
< error
| 1064 42000 …line 1 column 22 near "select 1)" 
> exec ^
: string
| create table t (select 1) union (select 2)
< error
| 1064 42000 …line 1 column 22 near "select 1) union (select 2)" 
> exec ^
: string
| create table t (select 1) order by 1
< error
| 1064 42000 …line 1 column 22 near "select 1) order by 1" 
> exec ^
: string
| create table t ((select 1))
< error
| 1064 42000 …line 1 column 17 near "(select 1))" 
> exec ^
: string
| create table t ((select 1)) union (select 2)
< error
| 1064 42000 …line 1 column 17 near "(select 1)) union (select 2)" 
> exec ^
: string
| create table t (with c as (select 1) select * from c)
< error
| 1064 42000 …line 1 column 20 near "with c as (select 1) select * from c)" 
> exec ^
: string
| create table t (table u)
< error
| 1064 42000 …line 1 column 21 near "table u)" 
> exec ^
: string
| create table t (values row(1))
< error
| 1064 42000 …line 1 column 22 near "values row(1))" 
> exec ^
: string
| create table t (a int) engine=innodb (select 1)
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) engine=innodb (select 1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) engine=innodb table u
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) engine=innodb table u
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) engine=innodb values row(1)
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) engine=innodb values row(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) engine=innodb with c as (select 1) select * from c
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) engine=innodb with c as (select 1) select * from c
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t engine=innodb table u
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t engine=innodb table u
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) partition by hash(a) partitions 2 table u
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) partition by hash(a) partitions 2 table u
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) partition by hash(a) partitions 2 values row(1)
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) partition by hash(a) partitions 2 values row(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) partition by hash(a) partitions 2 with c as (select 1) select * from c
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) partition by hash(a) partitions 2 with c as (select 1) select * from c
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) partition by hash(a) (partition p0) (select 1)
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) partition by hash(a) (partition p0) (select 1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) partition by hash(a) partitions 2 (select 1)
< error
| 1064 42000 …line 1 column 64 near "select 1)" 
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
| 1064 42000 …line 1 column 37 near "" 
> exec ^
: string
| create table t (a int) engine=innodb, select 1
< error
| 1064 42000 …line 1 column 44 near "select 1" 
> exec ^
: string
| create table t (a int) engine=innodb, (select 1)
< error
| 1064 42000 …line 1 column 39 near "(select 1)" 
> exec ^
: string
| create table t (a int) engine=innodb, partition by hash(a)
< error
| 1064 42000 …line 1 column 47 near "partition by hash(a)" 
> exec ^
: string
| create table t (a int) engine=innodb,, comment 'x'
< error
| 1064 42000 …line 1 column 38 near ", comment 'x'" 
> exec ^
: string
| create table t engine=innodb, select 1
< error
| 1064 42000 …line 1 column 36 near "select 1" 
> exec ^
: string
| create table t (a int) ignore (select 1)
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) ignore (select 1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) ignore table u
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) ignore table u
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) ignore with c as (select 1) select * from c
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) ignore with c as (select 1) select * from c
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) replace table u
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) replace table u
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) replace values row(1)
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) replace values row(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) partition by hash(a) partitions 2 ignore (select 1)
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) partition by hash(a) partitions 2 ignore (select 1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t ignore select 1
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t ignore select 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) ignore as select 1
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) ignore as select 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) ignore as table u
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) ignore as table u
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) as ignore select 1
< error
| 1064 42000 …line 1 column 32 near "ignore select 1" 
> exec ^
: string
| create table t (a int) ignore
< ok 0 0
> exec ^
: string
| create table t (a int) replace
< ok 0 0
> exec ^
: string
| create table t ignore
< error
| 1113 42000 A table must have at least 1 column
> exec ^
: string
| create table t ignore
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t as table u
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t as table u
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t as values row(1)
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t as values row(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t as (select 1)
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t as (select 1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) as (select 1)
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) as (select 1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) as table u
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) as table u
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) as with c as (select 1) select * from c
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) as with c as (select 1) select * from c
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) as
< ok 0 0
> exec ^
: string
| create table t (a int) ignore as
< ok 0 0
> exec ^
: string
| create table t (a int) replace as
< ok 0 0
> exec ^
: string
| create table t as
< error
| 1113 42000 A table must have at least 1 column
> exec ^
: string
| create table t as
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t ignore as
< error
| 1113 42000 A table must have at least 1 column
> exec ^
: string
| create table t ignore as
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t replace as
< error
| 1113 42000 A table must have at least 1 column
> exec ^
: string
| create table t replace as
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create temporary table t (a int) as
< ok 0 0
> exec ^
: string
| create table t (a int) engine=innodb as
< ok 0 0
> exec ^
: string
| create table t engine=innodb as
< error
| 1113 42000 A table must have at least 1 column
> exec ^
: string
| create table t engine=innodb as
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) partition by hash(a) partitions 2 as
< ok 0 0
> exec ^
: string
| create table t (a int) engine=innodb ignore as
< ok 0 0
> exec ^
: string
| create table t (a int) as;
< ok 0 0
> exec ^
: string
| create global temporary table t (a int) as on commit delete rows
< ok 0 0
> exec ^
: string
| create global temporary table t (a int) ignore as on commit delete rows
< ok 0 0
> exec ^
: string
| create global temporary table t as on commit delete rows
< error
| 1113 42000 A table must have at least 1 column
> exec ^
: string
| create global temporary table t as on commit delete rows
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create global temporary table t (a int) as
< error
| 1064 42000 …line 1 column 42 near ""GLOBAL TEMPORARY and ON COMMIT DELETE ROWS must appear together 
> exec ^
: string
| create table t (a int) as on commit delete rows
< error
| 1064 42000 …line 1 column 47 near ""GLOBAL TEMPORARY and ON COMMIT DELETE ROWS must appear together 
> exec ^
: string
| create table t as ignore
< error
| 1064 42000 …line 1 column 24 near "ignore" 
> exec ^
: string
| create table t (a int) as replace
< error
| 1064 42000 …line 1 column 33 near "replace" 
> exec ^
: string
| create table t (a int) as ignore as
< error
| 1064 42000 …line 1 column 32 near "ignore as" 
> exec ^
: string
| create table t as as
< error
| 1064 42000 …line 1 column 20 near "as" 
> exec ^
: string
| create table t as like u
< error
| 1064 42000 …line 1 column 22 near "like u" 
> exec ^
: string
| create table t as foo
< error
| 1064 42000 …line 1 column 21 near "foo" 
> exec ^
: string
| create table t (a int) as engine=innodb
< error
| 1064 42000 …line 1 column 32 near "engine=innodb" 
> exec ^
: string
| create table t (a int) as partition by hash(a) partitions 2
< error
| 1064 42000 …line 1 column 35 near "partition by hash(a) partitions 2" 
> exec ^
: string
| create table t (a int) as (a int)
< error
| 1064 42000 …line 1 column 28 near "a int)" 
> exec ^
: string
| create table t as select
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| create table t (a int) select 1 for update
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) select 1 for update
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) (select 1) for update
< error
| 1064 42000 …line 1 column 37 near "for update" 
> exec ^
: string
| create table t (a int) table u for update
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) table u for update
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) (select 1) union (select 2) for update
< error
| 1064 42000 …line 1 column 54 near "for update" 
> exec ^
: string
| create table t (a int) select 1 union select 2 for update
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (a int) select 1 union select 2 for update
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) (select 1) order by 1 for update
< error
| 1064 42000 …line 1 column 48 near "for update" 
> exec ^
: string
| create table t as (select 1) for update
< error
| 1064 42000 …line 1 column 32 near "for update" 
> exec ^
: string
| create table t (a int) as (select 1) for update
< error
| 1064 42000 …line 1 column 40 near "for update" 
> exec ^
: string
| create table t (a int) as select 1 select 2
< error
| 1064 42000 …line 1 column 41 near "select 2" 
> exec ^
: string
| create table t (a int) table u table v
< error
| 1064 42000 …line 1 column 36 near "table v" 
> exec ^
: string
| create table t (a int) values row(1) select 1
< error
| 1064 42000 …line 1 column 43 near "select 1" 
> exec ^
: string
| create table t (a int) (select 1) (select 2)
< error
| 1064 42000 …line 1 column 35 near "(select 2)" 
> exec ^
: string
| create table t (a int) select 1 into @x
< error
| 1064 42000 …line 1 column 39 near "@x" 
