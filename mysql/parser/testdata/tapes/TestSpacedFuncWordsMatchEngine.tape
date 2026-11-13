-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestSpacedFuncWordsMatchEngine default#1
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| select count(*)
< ok 0 0
> exec ^
: string
| select count ( * )
< error
| 1064 42000 …near '* )' at line 1
> exec ^
: string
| select count(distinct a)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select count(distinct a)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select count ( distinct a )
< error
| 1064 42000 …near 'distinct a )' at line 1
> exec ^
: string
| select count ( all a )
< error
| 1064 42000 …near 'all a )' at line 1
> exec ^
: string
| select count(a)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select count(a)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select count ( a )
< error
| 1630 42000 FUNCTION omni_t_mysql_parser.count does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| select count ( a )
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select count(a,a)
< error
| 1064 42000 …near ',a)' at line 1
> exec ^
: string
| select count ( a, a )
< error
| 1630 42000 FUNCTION omni_t_mysql_parser.count does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| select count ( a, a )
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select count ( )
< error
| 1630 42000 FUNCTION omni_t_mysql_parser.count does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| select count ( )
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select max(a)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select max(a)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select max ( a )
< error
| 1630 42000 FUNCTION omni_t_mysql_parser.max does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| select max ( a )
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select max(distinct a)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select max(distinct a)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select max ( distinct a )
< error
| 1064 42000 …near 'distinct a )' at line 1
> exec ^
: string
| select max(a,a)
< error
| 1064 42000 …near ',a)' at line 1
> exec ^
: string
| select max ( a, a )
< error
| 1630 42000 FUNCTION omni_t_mysql_parser.max does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| select max ( a, a )
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select max ( )
< error
| 1630 42000 FUNCTION omni_t_mysql_parser.max does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| select max ( )
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select sum(a)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select sum(a)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select sum ( a )
< error
| 1630 42000 FUNCTION omni_t_mysql_parser.sum does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| select sum ( a )
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select sum ( * )
< error
| 1064 42000 …near '* )' at line 1
> exec ^
: string
| select sum ( distinct a )
< error
| 1064 42000 …near 'distinct a )' at line 1
> exec ^
: string
| select min(a)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select min(a)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select min ( a )
< error
| 1630 42000 FUNCTION omni_t_mysql_parser.min does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| select min ( a )
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select group_concat(a)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select group_concat(a)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select group_concat ( a )
< error
| 1630 42000 FUNCTION omni_t_mysql_parser.group_concat does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| select group_concat ( a )
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select group_concat(distinct a)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select group_concat(distinct a)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select group_concat ( distinct a )
< error
| 1064 42000 …near 'distinct a )' at line 1
> exec ^
: string
| select group_concat(a separator ',')
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select group_concat(a separator ',')
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select group_concat ( a separator ',' )
< error
| 1064 42000 …near 'separator ',' )' at line 1
> exec ^
: string
| select group_concat ( a, a )
< error
| 1630 42000 FUNCTION omni_t_mysql_parser.group_concat does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| select group_concat ( a, a )
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select std(a)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select std(a)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select std ( a )
< error
| 1630 42000 FUNCTION omni_t_mysql_parser.std does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| select std ( a )
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select std ( * )
< error
| 1064 42000 …near '* )' at line 1
> exec ^
: string
| select stddev ( a )
< error
| 1630 42000 FUNCTION omni_t_mysql_parser.stddev does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| select stddev ( a )
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select stddev_pop ( a )
< error
| 1630 42000 FUNCTION omni_t_mysql_parser.stddev_pop does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| select stddev_pop ( a )
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select stddev_samp ( a )
< error
| 1630 42000 FUNCTION omni_t_mysql_parser.stddev_samp does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| select stddev_samp ( a )
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select var_pop ( a )
< error
| 1630 42000 FUNCTION omni_t_mysql_parser.var_pop does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| select var_pop ( a )
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select var_samp ( a )
< error
| 1630 42000 FUNCTION omni_t_mysql_parser.var_samp does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| select var_samp ( a )
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select variance ( a )
< error
| 1630 42000 FUNCTION omni_t_mysql_parser.variance does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| select variance ( a )
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select bit_and ( a )
< error
| 1630 42000 FUNCTION omni_t_mysql_parser.bit_and does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| select bit_and ( a )
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select bit_or ( a )
< error
| 1630 42000 FUNCTION omni_t_mysql_parser.bit_or does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| select bit_or ( a )
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select bit_xor ( a )
< error
| 1630 42000 FUNCTION omni_t_mysql_parser.bit_xor does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| select bit_xor ( a )
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select bit_and ( * )
< error
| 1064 42000 …near '* )' at line 1
> exec ^
: string
| select json_arrayagg(a)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select json_arrayagg(a)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select json_arrayagg ( a )
< error
| 1630 42000 FUNCTION omni_t_mysql_parser.json_arrayagg does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| select json_arrayagg ( a )
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select json_arrayagg(a,a)
< error
| 1064 42000 …near ',a)' at line 1
> exec ^
: string
| select json_arrayagg ( a, a )
< error
| 1630 42000 FUNCTION omni_t_mysql_parser.json_arrayagg does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| select json_arrayagg ( a, a )
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select json_objectagg(a, b)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select json_objectagg(a, b)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select json_objectagg ( a, b )
< error
| 1630 42000 FUNCTION omni_t_mysql_parser.json_objectagg does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| select json_objectagg ( a, b )
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select json_objectagg(a)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select json_objectagg ( a )
< error
| 1630 42000 FUNCTION omni_t_mysql_parser.json_objectagg does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| select json_objectagg ( a )
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select substr(a, 1)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select substr(a, 1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select substr ( a, 1 )
< error
| 1630 42000 FUNCTION omni_t_mysql_parser.substr does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| select substr ( a, 1 )
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select substr(a from 1)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select substr(a from 1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select substr ( a from 1 )
< error
| 1064 42000 …near 'from 1 )' at line 1
> exec ^
: string
| select substring(a from 1)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select substring(a from 1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select substring ( a from 1 )
< error
| 1064 42000 …near 'from 1 )' at line 1
> exec ^
: string
| select substring ( a )
< error
| 1630 42000 FUNCTION omni_t_mysql_parser.substring does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| select substring ( a )
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select trim(a)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select trim(a)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select trim ( a )
< error
| 1630 42000 FUNCTION omni_t_mysql_parser.trim does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| select trim ( a )
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select trim(leading 'x' from a)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select trim(leading 'x' from a)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select trim ( leading 'x' from a )
< error
| 1064 42000 …near 'leading 'x' from a )' at line 1
> exec ^
: string
| select position('a' in b)
< error
| 1054 42S22 Unknown column 'b' in 'field list'
> exec ^
: string
| select position('a' in b)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select position ( 'a' in b )
< error
| 1064 42000 …near 'b )' at line 1
> exec ^
: string
| select position(a)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select position ( a )
< error
| 1630 42000 FUNCTION omni_t_mysql_parser.position does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| select position ( a )
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select avg(a)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select avg(a)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select avg ( a )
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select avg ( a )
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select avg(distinct a)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select avg(distinct a)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select avg ( distinct a )
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select avg ( distinct a )
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select avg(*)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| select avg ( * )
< error
| 1064 42000 …near '* )' at line 1
> exec ^
: string
| select char(a using utf8)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select char(a using utf8)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select char ( a using utf8 )
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select char ( a using utf8 )
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select abc(a)
< error
| 1305 42000 FUNCTION omni_t_mysql_parser.abc does not exist
> exec ^
: string
| select abc(a)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select abc ( a )
< error
| 1305 42000 FUNCTION omni_t_mysql_parser.abc does not exist
> exec ^
: string
| select abc ( a )
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select abc ( * )
< error
| 1064 42000 …near '* )' at line 1
> exec ^
: string
| select foo (x)
< error
| 1305 42000 FUNCTION omni_t_mysql_parser.foo does not exist
> exec ^
: string
| select foo (x)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create table t (a int) partition by hash(sum(a))
< ok 0 0
> exec ^
: string
| create table t (a int) partition by hash(sum (a))
< error
| 1064 42000 Constant, random or timezone-dependent expressions in (sub)partitioning function are not allowed near 'sum (a))' at line 1
> exec ^
: string
| create table t (a int) partition by hash(max (a))
< error
| 1064 42000 Constant, random or timezone-dependent expressions in (sub)partitioning function are not allowed near 'max (a))' at line 1
> exec ^
: string
| create table t (a int) partition by hash(bit_and (a))
< error
| 1064 42000 Constant, random or timezone-dependent expressions in (sub)partitioning function are not allowed near 'bit_and (a))' at line 1
> exec ^
: string
| create table t (a int) partition by hash(`sum`(a))
< error
| 1064 42000 Constant, random or timezone-dependent expressions in (sub)partitioning function are not allowed near '`sum`(a))' at line 1
