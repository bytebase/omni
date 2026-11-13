-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestTiDBSpacedFuncWordsMatchEngine default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| select count(*)
< ok 0 0
> exec ^
: string
| select count ( * )
< error
| 1064 42000 …line 1 column 16 near "* )" 
> exec ^
: string
| select count(distinct a)
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| select count(distinct a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select count ( distinct a )
< error
| 1064 42000 …line 1 column 23 near "distinct a )" 
> exec ^
: string
| select count ( all a )
< error
| 1064 42000 …line 1 column 18 near "all a )" 
> exec ^
: string
| select count(a)
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| select count(a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select count ( a )
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select count ( a )
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select count(a,a)
< error
| 1064 42000 …line 1 column 15 near ",a)" 
> exec ^
: string
| select count ( a, a )
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select count ( a, a )
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select count ( 1 )
< error
| 1305 42000 FUNCTION test.count does not exist
> exec ^
: string
| select count ( 1 )
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select max(a)
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| select max(a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select max ( a )
< error
| 1064 42000 …line 1 column 12 near "( a )" 
> exec ^
: string
| select max('bbb')
< ok 0 0
> exec ^
: string
| select max ( 'bbb' )
< error
| 1064 42000 …line 1 column 12 near "( 'bbb' )" 
> exec ^
: string
| select max(distinct a)
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| select max(distinct a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select max ( distinct a )
< error
| 1064 42000 …line 1 column 12 near "( distinct a )" 
> exec ^
: string
| select sum(a)
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| select sum(a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select sum ( a )
< error
| 1064 42000 …line 1 column 12 near "( a )" 
> exec ^
: string
| select min(a)
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| select min(a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select min ( a )
< error
| 1064 42000 …line 1 column 12 near "( a )" 
> exec ^
: string
| select group_concat(a)
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| select group_concat(a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select group_concat ( a )
< error
| 1064 42000 …line 1 column 21 near "( a )" 
> exec ^
: string
| select group_concat(a separator ',')
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| select group_concat(a separator ',')
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select group_concat ( a separator ',' )
< error
| 1064 42000 …line 1 column 21 near "( a separator ',' )" 
> exec ^
: string
| select std(a)
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| select std(a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select std ( a )
< error
| 1064 42000 …line 1 column 12 near "( a )" 
> exec ^
: string
| select stddev ( a )
< error
| 1064 42000 …line 1 column 15 near "( a )" 
> exec ^
: string
| select stddev_pop ( a )
< error
| 1064 42000 …line 1 column 19 near "( a )" 
> exec ^
: string
| select stddev_samp(a)
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| select stddev_samp(a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select stddev_samp ( a )
< error
| 1064 42000 …line 1 column 20 near "( a )" 
> exec ^
: string
| select var_pop(a)
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| select var_pop(a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select var_pop ( a )
< error
| 1064 42000 …line 1 column 16 near "( a )" 
> exec ^
: string
| select var_samp ( a )
< error
| 1064 42000 …line 1 column 17 near "( a )" 
> exec ^
: string
| select variance ( a )
< error
| 1064 42000 …line 1 column 17 near "( a )" 
> exec ^
: string
| select bit_and(a)
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| select bit_and(a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select bit_and ( a )
< error
| 1064 42000 …line 1 column 16 near "( a )" 
> exec ^
: string
| select bit_or ( a )
< error
| 1064 42000 …line 1 column 15 near "( a )" 
> exec ^
: string
| select bit_xor ( a )
< error
| 1064 42000 …line 1 column 16 near "( a )" 
> exec ^
: string
| select approx_count_distinct(a)
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| select approx_count_distinct(a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select approx_count_distinct ( a )
< error
| 1064 42000 …line 1 column 30 near "( a )" 
> exec ^
: string
| select approx_percentile(a, 50)
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| select approx_percentile(a, 50)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select approx_percentile ( a, 50 )
< error
| 1064 42000 …line 1 column 26 near "( a, 50 )" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select substr ( a, 1 )
< error
| 1064 42000 …line 1 column 15 near "( a, 1 )" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select substring ( a from 1 )
< error
| 1064 42000 …line 1 column 18 near "( a from 1 )" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select trim ( a )
< error
| 1064 42000 …line 1 column 13 near "( a )" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select trim ( leading 'x' from a )
< error
| 1064 42000 …line 1 column 13 near "( leading 'x' from a )" 
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
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select position ( 'a' in b )
< error
| 1064 42000 …line 1 column 17 near "( 'a' in b )" 
> exec ^
: string
| select avg(a)
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| select avg(a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select avg ( a )
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| select avg ( a )
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select avg(distinct a)
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| select avg(distinct a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select avg ( distinct a )
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| select avg ( distinct a )
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select avg(*)
< error
| 1064 42000 …line 1 column 12 near "*)" 
> exec ^
: string
| select avg ( * )
< error
| 1064 42000 …line 1 column 14 near "* )" 
> exec ^
: string
| select json_arrayagg(a)
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| select json_arrayagg(a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select json_arrayagg (a)
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| select json_arrayagg (a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select json_objectagg(a, a)
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| select json_objectagg(a, a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select json_objectagg (a, a)
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| select json_objectagg (a, a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select sysdate(b)
< error
| 1064 42000 …line 1 column 16 near "b)" 
> exec ^
: string
| select sysdate ( b )
< error
| 1054 42S22 Unknown column 'b' in 'field list'
> exec ^
: string
| select sysdate ( b )
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select now()
< ok 0 0
> exec ^
: string
| select now ( )
< error
| 1064 42000 …line 1 column 12 near "( )" 
> exec ^
: string
| select curdate(1)
< error
| 1064 42000 …line 1 column 16 near "1)" 
> exec ^
: string
| select curdate ( 1 )
< error
| 1064 42000 …line 1 column 16 near "( 1 )" 
> exec ^
: string
| select abc(a)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select abc(a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select abc ( a )
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select abc ( a )
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select abc ( * )
< error
| 1064 42000 …line 1 column 14 near "* )" 
> exec ^
: string
| select foo (x)
< error
| 1054 42S22 Unknown column 'x' in 'field list'
> exec ^
: string
| select foo (x)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
