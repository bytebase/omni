-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestExecFallback_OceanBase default#1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select 1 /*! select 2 */
< ok 0 0
> query
| select 1 /*! select 2 */
< result
< column "1" "BIGINT"
< row
: int64 1
> query
| select 1
< result
< column "1" "BIGINT"
< row
: int64 1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select 1 /*!50000 select 2 */
< ok 0 0
> query
| select 1 /*!50000 select 2 */
< result
< column "1" "BIGINT"
< row
: int64 1
> query
| select 1
< result
< column "1" "BIGINT"
< row
: int64 1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select 1 /*! from */
< ok 0 0
> query
| select 1 /*! from */
< result
< column "1" "BIGINT"
< row
: int64 1
> query
| select 1
< result
< column "1" "BIGINT"
< row
: int64 1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select 1 /*!50000 from */
< ok 0 0
> query
| select 1 /*!50000 from */
< result
< column "1" "BIGINT"
< row
: int64 1
> query
| select 1
< result
< column "1" "BIGINT"
< row
: int64 1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select 1 /*! from dual */ , 2
< ok 0 0
> query
| select 1 /*! from dual */ , 2
< result
< column "1" "BIGINT"
< column "2" "BIGINT"
< row
: int64 1
: int64 2
> query
| select 1 , 2
< result
< column "1" "BIGINT"
< column "2" "BIGINT"
< row
: int64 1
: int64 2
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select 1 /*! as */ , 2
< ok 0 0
> query
| select 1 /*! as */ , 2
< result
< column "1" "BIGINT"
< column "2" "BIGINT"
< row
: int64 1
: int64 2
> query
| select 1 , 2
< result
< column "1" "BIGINT"
< column "2" "BIGINT"
< row
: int64 1
: int64 2
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select 1 /*! + 1 */ /*! from */
< ok 0 0
> query
| select 1 /*! + 1 */ /*! from */
< result
< column "1" "BIGINT"
< row
: int64 1
> query
| select 1
< result
< column "1" "BIGINT"
< row
: int64 1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select 1 /*! from */ /*! + 1 */
< ok 0 0
> query
| select 1 /*! from */ /*! + 1 */
< result
< column "1" "BIGINT"
< row
: int64 1
> query
| select 1
< result
< column "1" "BIGINT"
< row
: int64 1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select 1 /*! + 1 */ + 1 /*! from */
< ok 0 0
> query
| select 1 /*! + 1 */ + 1 /*! from */
< result
< column "1 /*! + 1 */ + 1" "BIGINT"
< row
: int64 2
> query
| select 1 + 1
< result
< column "1 + 1" "BIGINT"
< row
: int64 2
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select 1 /*! + 10 */ + 100 /*! , 2 */ /*! from */
< ok 0 0
> query
| select 1 /*! + 10 */ + 100 /*! , 2 */ /*! from */
< result
< column "1 /*! + 10 */ + 100" "BIGINT"
< row
: int64 101
> query
| select 1 + 100
< result
< column "1 + 100" "BIGINT"
< row
: int64 101
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select 1 /*! + 1 /*! from */
< ok 0 0
> query
| select 1 /*! + 1 /*! from */
< result
< column "1" "BIGINT"
< row
: int64 1
> query
| select 1
< result
< column "1" "BIGINT"
< row
: int64 1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select 1 /*! from /*! + 1 */
< ok 0 0
> query
| select 1 /*! from /*! + 1 */
< result
< column "1" "BIGINT"
< row
: int64 1
> query
| select 1
< result
< column "1" "BIGINT"
< row
: int64 1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select 1 /*! + 1 */
< ok 0 0
> query
| select 1 /*! + 1 */
< result
< column "1 /*! + 1" "BIGINT"
< row
: int64 2
> query
| select 1 + 1
< result
< column "1 + 1" "BIGINT"
< row
: int64 2
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select 1 /*! , 2 */
< ok 0 0
> query
| select 1 /*! , 2 */
< result
< column "1" "BIGINT"
< column "2" "BIGINT"
< row
: int64 1
: int64 2
> query
| select 1 , 2
< result
< column "1" "BIGINT"
< column "2" "BIGINT"
< row
: int64 1
: int64 2
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select /*! 1 */
< ok 0 0
> query
| select /*! 1 */
< result
< column "1" "BIGINT"
< row
: int64 1
> query
| select 1
< result
< column "1" "BIGINT"
< row
: int64 1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select 1 /*! from dual */ where 1 = 1
< ok 0 0
> query
| select 1 /*! from dual */ where 1 = 1
< result
< column "1" "BIGINT"
< row
: int64 1
> query
| select 1 from dual where 1 = 1
< result
< column "1" "BIGINT"
< row
: int64 1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| /*!50000 select 1 */ /*!50000 + 1 */
< ok 0 0
> query
| /*!50000 select 1 */ /*!50000 + 1 */
< result
< column "1 */ /*!50000 + 1" "BIGINT"
< row
: int64 2
> query
| select 1 + 1
< result
< column "1 + 1" "BIGINT"
< row
: int64 2
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| /*! select 1 */
< ok 0 0
> query
| /*! select 1 */
< result
< column "1" "BIGINT"
< row
: int64 1
> query
| select 1
< result
< column "1" "BIGINT"
< row
: int64 1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| /*! select */ 1
< ok 0 0
> query
| /*! select */ 1
< result
< column "1" "BIGINT"
< row
: int64 1
> query
| select 1
< result
< column "1" "BIGINT"
< row
: int64 1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| /*! x */ select 1
< ok 0 0
> query
| /*! x */ select 1
< result
< column "1" "BIGINT"
< row
: int64 1
> query
| select 1
< result
< column "1" "BIGINT"
< row
: int64 1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| /*! select 1 from */
< ok 0 0
> query
| /*! select 1 from */
< result
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| /*! select from */
< ok 0 0
> query
| /*! select from */
< result
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| /*! select 1 */ /*! from */
< ok 0 0
> query
| /*! select 1 */ /*! from */
< result
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| /*!50000 SELECT FROM t */
< ok 0 0
> query
| /*!50000 SELECT FROM t */
< result
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| /*! SELECT 'abc */
< ok 0 0
> query
| /*! SELECT 'abc */
< result
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| /*!50100 SELECT 'x */
< ok 0 0
> query
| /*!50100 SELECT 'x */
< result
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| /*!80000 SELECT a FROM WHERE */
< ok 0 0
> query
| /*!80000 SELECT a FROM WHERE */
< result
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| /*! select 1 */ from
< error
| 1064 42000 …near 'from' at line 1
> exec ^
: string
| /*! select 1 */ + 1 /*! from */
< error
| 1064 42000 …near '+ 1 /*! from */' at line 1
> exec ^
: string
| select 1 /*! + 1
< error
| 1064 42000 …near 'select 1 /*! + 1' at line 2
> exec ^
: string
| select 1 /*! from
< error
| 1064 42000 …near 'm' at line 1
> exec ^
: string
| select 1 /*! + 'a*/ -- '
< error
| 1064 42000 …near ''a*/ -- '' at line 2
> exec ^
: string
| select 1 /*! from 'a*/ , 2 -- '
< ok 0 0
> query
| select 1 /*! from 'a*/ , 2 -- '
< result
< column "1" "BIGINT"
< column "2" "BIGINT"
< row
: int64 1
: int64 2
> query
| select 1 , 2
< result
< column "1" "BIGINT"
< column "2" "BIGINT"
< row
: int64 1
: int64 2
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select 1 /*! from -- c */
|
< ok 0 0
> query
| select 1 /*! from -- c */
|
< result
< column "1" "BIGINT"
< row
: int64 1
> query
| select 1
< result
< column "1" "BIGINT"
< row
: int64 1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select 1 /*! , 'abc */
< ok 0 0
> query
| select 1 /*! , 'abc */
< result
< column "1" "BIGINT"
< row
: int64 1
> query
| select 1
< result
< column "1" "BIGINT"
< row
: int64 1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select 1 /*! , 2 */ from
< error
| 1064 42000 …near 'from' at line 1
> exec ^
: string
| select 1 /*! from */ */
< error
| 1064 42000 …near '/' at line 1
> exec ^
: string
| select 1 /*! + 1 /*! + 2 */ from */
< error
| 1064 42000 …near '*/' at line 1
> exec ^
: string
| select 1 /*! table t */
< ok 0 0
> query
| select 1 /*! table t */
< result
< column "1" "BIGINT"
< row
: int64 1
> query
| select 1
< result
< column "1" "BIGINT"
< row
: int64 1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| /*! x */ table t
< error
| 1064 42000 …near 'table t' at line 1
> exec ^
: string
| (select 1) /*! for update */
< ok 0 0
> exec ^
: string
| (select 1) /*! from */ for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| (select 1) order by 1 /*! from */ for update
< ok 0 0
> exec ^
: string
| create table t (a int) /*! engine */
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (a int) /*!50000 select 2 */
< error
| 1046 3D000 No database selected
> exec ^
: string
| create database d /*!40100 default character set utf8mb4 */ /*!80016 default encryption='N' */
< ok 0 0
> exec ^
: string
| insert into t values (1) /*! from */
< error
| 1046 3D000 No database selected
> exec ^
: string
| update t set a = 2 /*! from */
< error
| 1046 3D000 No database selected
> exec ^
: string
| delete from t /*! from */
< error
| 1046 3D000 No database selected
> exec ^
: string
| set @v = 1 /*! from */
< ok 0 0
> exec ^
: string
| create view v as select 1 /*! from */
< error
| 1046 3D000 No database selected
> exec ^
: string
| create or replace view v as select 1 /*! from */
< error
| 1046 3D000 No database selected
> exec ^
: string
| create /*! definer=root@'%' */ view v as select 1 /*! from */
< error
| 1046 3D000 No database selected
> exec ^
: string
| drop table if exists t /*! from */
< error
| 1046 3D000 No database selected
> exec ^
: string
| begin /*! from */
< ok 0 0
> exec ^
: string
| begin work /*! from */
< ok 0 0
> exec ^
: string
| drop/**/function if exists f /*! from */
< ok 0 0
> exec ^
: string
| create aggregate function f returns int /*! from */ soname 'f.so'
< error
| 1235 0A000 creating loadable function not supported
> exec ^
: string
| /*! x */ get diagnostics @n = number
< ok 0 0
> exec ^
: string
| call p() /*! from */
< error
| 1064 42000 …near 'from' at line 1
> exec ^
: string
| call p() /*!50000 from t */
< error
| 1064 42000 …near 'from' at line 1
> exec ^
: string
| /* c */ call p() /*! from */
< error
| 1064 42000 …near 'from' at line 1
> exec ^
: string
| -- c
| call p() /*! from */
< error
| 1064 42000 …near 'from' at line 2
> exec ^
: string
| call p() /*! , 2 */
< error
| 1064 42000 …near ',' at line 1
> exec ^
: string
| signal sqlstate '45000' /*! from */
< error
| 1064 42000 …near 'from' at line 1
> exec ^
: string
| drop procedure if exists p /*! from */
< error
| 1064 42000 …near 'from' at line 1
> exec ^
: string
| alter procedure p /*! from */ comment 'x'
< error
| 1064 42000 …near 'from' at line 1
> exec ^
: string
| create procedure p() /*! from */ begin select 1; end
< error
| 1064 42000 …near 'from' at line 1
> exec ^
: string
| create or replace procedure p() /*! from */ select 1
< error
| 1064 42000 …near 'procedure' at line 1
> exec ^
: string
| create procedure p() begin declare x int /*! from */; select 1; end
< error
| 1064 42000 …near 'from' at line 1
> exec ^
: string
| create definer=root@'%' view v as select 1 /*! from */
< error
| 1064 42000 …near 'view' at line 1
> exec ^
: string
| create event e /*! from */ on schedule every 1 hour do select 1
< error
| 1064 42000 …near 'from' at line 1
> exec ^
: string
| create event e /*! if not exists */ on schedule every 1 hour do select 1
< error
| 1064 42000 …near 'if' at line 1
> exec ^
: string
| alter event e /*! from */ disable
< error
| 1064 42000 …near 'from' at line 1
> exec ^
: string
| drop event if exists e /*! from */
< error
| 1064 42000 …near 'from' at line 1
> exec ^
: string
| create function f returns int /*! from */ soname 'f.so'
< error
| 1235 0A000 creating loadable function not supported
> exec ^
: string
| create function f() returns int /*! from */ return 1
< error
| 1064 42000 …near '() returns int /*! from */ return 1' at line 1
> exec ^
: string
| create function f /*! from */ () returns int return 1
< error
| 1064 42000 …near '() returns int return 1' at line 1
> exec ^
: string
| do(1) /*! from */
< error
| 1064 42000 …near 'do(1) /*! from */' at line 2
> exec ^
: string
| call/**/p() /*! from */
< error
| 1064 42000 …near 'call/**/p() /*! from */' at line 1
> exec ^
: string
| drop/**/procedure if exists p /*! from */
< error
| 1064 42000 …near 'procedure if exists p /*! from */' at line 1
> exec ^
: string
| drop/**/trigger if exists tr /*! from */
< error
| 1064 42000 …near 'trigger if exists tr /*! from */' at line 1
> exec ^
: string
| create/**/procedure p() select 1 /*! from */
< error
| 1064 42000 …near 'procedure p() select 1 /*! from */' at line 1
> exec ^
: string
| /*! x */ call p()
< error
| 1064 42000 …near 'call p()' at line 1
> exec ^
: string
| /*! x */ do 1
< error
| 1064 42000 …near '/*! x */ do 1' at line 2
> exec ^
: string
| /*! x */ signal sqlstate '45000'
< error
| 1064 42000 …near 'signal sqlstate '45000'' at line 1
> exec ^
: string
| /*! x */ create procedure p() select 1
< error
| 1064 42000 …near 'procedure p() select 1' at line 1
> exec ^
: string
| /*! x */ create trigger tr before insert on t for each row set new.a = 1
< error
| 1064 42000 …near 'trigger tr before insert on t for each row set new.a = 1' at line 1
> exec ^
: string
| /*! x */ drop trigger if exists tr
< error
| 1064 42000 …near 'trigger if exists tr' at line 1
> exec ^
: string
| /*! x */ alter procedure p comment 'x'
< error
| 1064 42000 …near 'procedure p comment 'x'' at line 1
> exec ^
: string
| /*! x */ create event e on schedule every 1 hour do select 1
< error
| 1064 42000 …near 'event e on schedule every 1 hour do select 1' at line 1
> exec ^
: string
| select 1 /*!50000 + 1 */
< ok 0 0
> query
| select 1 /*!50000 + 1 */
< result
< column "1 /*!50000 + 1" "BIGINT"
< row
: int64 2
> query
| select 1 + 1
< result
< column "1 + 1" "BIGINT"
< row
: int64 2
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select 1 /*!99999 + 1 */
< ok 0 0
> query
| select 1 /*!99999 + 1 */
< result
< column "1 /*!99999 + 1" "BIGINT"
< row
: int64 2
> query
| select 1 + 1
< result
< column "1 + 1" "BIGINT"
< row
: int64 2
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select /*!99999 3 */
< ok 0 0
> query
| select /*!99999 3 */
< result
< column "3" "BIGINT"
< row
: int64 3
> query
| select 3
< result
< column "3" "BIGINT"
< row
: int64 3
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select /*!00000 1 */
< ok 0 0
> query
| select /*!00000 1 */
< result
< column "1" "BIGINT"
< row
: int64 1
> query
| select 1
< result
< column "1" "BIGINT"
< row
: int64 1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select /*!12345 */ 7
< ok 0 0
> query
| select /*!12345 */ 7
< result
< column "7" "BIGINT"
< row
: int64 7
> query
| select 7
< result
< column "7" "BIGINT"
< row
: int64 7
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select 1 /*!50000+1 */
< ok 0 0
> query
| select 1 /*!50000+1 */
< result
< column "1 /*!50000+1" "BIGINT"
< row
: int64 2
> query
| select 1 + 1
< result
< column "1 + 1" "BIGINT"
< row
: int64 2
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select 1 /*!50000*/
< ok 0 0
> query
| select 1 /*!50000*/
< result
< column "1" "BIGINT"
< row
: int64 1
> query
| select 1
< result
< column "1" "BIGINT"
< row
: int64 1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select /*!1 */
< ok 0 0
> query
| select /*!1 */
< result
< column "1" "BIGINT"
< row
: int64 1
> query
| select 1
< result
< column "1" "BIGINT"
< row
: int64 1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select /*!1234 */
< ok 0 0
> query
| select /*!1234 */
< result
< column "1234" "BIGINT"
< row
: int64 1234
> query
| select 1234
< result
< column "1234" "BIGINT"
< row
: int64 1234
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select /*!5000*/
< ok 0 0
> query
| select /*!5000*/
< result
< column "5000" "BIGINT"
< row
: int64 5000
> query
| select 5000
< result
< column "5000" "BIGINT"
< row
: int64 5000
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select 1 /*!5000*/
< ok 0 0
> query
| select 1 /*!5000*/
< result
< column "1" "BIGINT"
< row
: int64 1
> query
| select 1
< result
< column "1" "BIGINT"
< row
: int64 1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select /*!123456 */
< ok 0 0
> query
| select /*!123456 */
< result
< column "6" "BIGINT"
< row
: int64 6
> query
| select 6
< result
< column "6" "BIGINT"
< row
: int64 6
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select /*!500001 */
< ok 0 0
> query
| select /*!500001 */
< result
< column "1" "BIGINT"
< row
: int64 1
> query
| select 1
< result
< column "1" "BIGINT"
< row
: int64 1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select 2 * /*!100005 */
< ok 0 0
> query
| select 2 * /*!100005 */
< result
< column "2 * /*!100005" "BIGINT"
< row
: int64 10
> query
| select 2 * 5
< result
< column "2 * 5" "BIGINT"
< row
: int64 10
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select 1 + /*!1234567 */
< ok 0 0
> query
| select 1 + /*!1234567 */
< result
< column "1 + /*!1234567" "BIGINT"
< row
: int64 68
> query
| select 1 + 67
< result
< column "1 + 67" "BIGINT"
< row
: int64 68
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select 1 /*!100000 + 1 */
< ok 0 0
> query
| select 1 /*!100000 + 1 */
< result
< column "1" "BIGINT"
< row
: int64 1
> query
| select 1
< result
< column "1" "BIGINT"
< row
: int64 1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select 1 /*!80000 select 2 */
< ok 0 0
> query
| select 1 /*!80000 select 2 */
< result
< column "1" "BIGINT"
< row
: int64 1
> query
| select 1
< result
< column "1" "BIGINT"
< row
: int64 1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select 1 /*!100000 select 2 */
< ok 0 0
> query
| select 1 /*!100000 select 2 */
< result
< column "1" "BIGINT"
< row
: int64 1
> query
| select 1
< result
< column "1" "BIGINT"
< row
: int64 1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select 1 /*!M + 1 */
< ok 0 0
> query
| select 1 /*!M + 1 */
< result
< column "1" "BIGINT"
< row
: int64 1
> query
| select 1
< result
< column "1" "BIGINT"
< row
: int64 1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select /*!*/ 1
< ok 0 0
> query
| select /*!*/ 1
< result
< column "1" "BIGINT"
< row
: int64 1
> query
| select 1
< result
< column "1" "BIGINT"
< row
: int64 1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select 1 /*!*/
< ok 0 0
> query
| select 1 /*!*/
< result
< column "1" "BIGINT"
< row
: int64 1
> query
| select 1
< result
< column "1" "BIGINT"
< row
: int64 1
