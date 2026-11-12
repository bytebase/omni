-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOceanBaseTemporalFormsMatchEngine default#1
> exec
| CREATE DATABASE IF NOT EXISTS omni_temporal
< ok 1 0
> exec
| USE omni_temporal
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a DATETIME, b INT, c INT, `now` INT, `curdate` INT, `curtime` INT, `date_add` INT, `date_sub` INT, `sysdate` INT, `adddate` INT, `subdate` INT)
< ok 0 0
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select current_timestamp
< ok 0 0
> exec ^
: string
| select current_timestamp()
< ok 0 0
> exec ^
: string
| select current_timestamp(0)
< ok 0 0
> exec ^
: string
| select current_timestamp(6)
< ok 0 0
> exec ^
: string
| select current_timestamp(7)
< error
| 1426 42000 Too big precision 7 specified for column 'now'. Maximum is 6.
> exec ^
: string
| select current_timestamp(255)
< error
| 1426 42000 Too big precision 255 specified for column 'now'. Maximum is 6.
> exec ^
: string
| select current_timestamp(256)
< error
| 1426 42000 Too big precision 256 specified for column 'now'. Maximum is 6.
> exec ^
: string
| select current_timestamp(2147483647)
< ok 0 0
> exec ^
: string
| select current_timestamp(2147483648)
< ok 0 0
> exec ^
: string
| select current_timestamp(4294967296)
< ok 0 0
> exec ^
: string
| select current_timestamp(18446744073709551615)
< ok 0 0
> exec ^
: string
| select current_timestamp(18446744073709551616)
< error
| 1064 42000 …near '18446744073709551616)' at line 1
> exec ^
: string
| select current_timestamp(00000000000000000000001)
< ok 0 0
> exec ^
: string
| select current_timestamp(-1)
< error
| 1064 42000 …near '-1)' at line 1
> exec ^
: string
| select current_timestamp(+1)
< error
| 1064 42000 …near '+1)' at line 1
> exec ^
: string
| select current_timestamp(1.0)
< error
| 1064 42000 …near '1.0)' at line 1
> exec ^
: string
| select current_timestamp(1e0)
< error
| 1064 42000 …near '1e0)' at line 1
> exec ^
: string
| select current_timestamp(.5)
< error
| 1064 42000 …near '.5)' at line 1
> exec ^
: string
| select current_timestamp('2')
< error
| 1064 42000 …near ''2')' at line 1
> exec ^
: string
| select current_timestamp(null)
< error
| 1064 42000 …near 'null)' at line 1
> exec ^
: string
| select current_timestamp(true)
< error
| 1064 42000 …near 'true)' at line 1
> exec ^
: string
| select current_timestamp(0x1)
< error
| 1064 42000 …near '0x1)' at line 1
> exec ^
: string
| select current_timestamp(x'01')
< error
| 1064 42000 …near 'x'01')' at line 1
> exec ^
: string
| select current_timestamp(b'1')
< error
| 1064 42000 …near 'b'1')' at line 1
> exec ^
: string
| select current_timestamp(b)
< error
| 1064 42000 …near 'b)' at line 1
> exec ^
: string
| select current_timestamp(@a)
< error
| 1064 42000 …near '@a)' at line 1
> exec ^
: string
| select current_timestamp(?)
< error
| 1064 42000 …near '?)' at line 1
> exec ^
: string
| select current_timestamp((1))
< error
| 1064 42000 …near '(1))' at line 1
> exec ^
: string
| select current_timestamp(1+1)
< error
| 1064 42000 …near '+1)' at line 1
> exec ^
: string
| select current_timestamp(1,2)
< error
| 1064 42000 …near ',2)' at line 1
> exec ^
: string
| select current_timestamp(1,)
< error
| 1064 42000 …near ',)' at line 1
> exec ^
: string
| select current_timestamp(
< error
| 1064 42000 …near '(' at line 1
> exec ^
: string
| select current_timestamp (1)
< ok 0 0
> exec ^
: string
| select current_timestamp /* c */ (1)
< ok 0 0
> exec ^
: string
| select current_timestamp (-1)
< error
| 1064 42000 …near '-1)' at line 1
> exec ^
: string
| select current_timestamp.a from t current_timestamp
< error
| 1064 42000 …near 'current_timestamp' at line 1
> exec ^
: string
| select current_timestamp a
< ok 0 0
> exec ^
: string
| select current_timestamp as a
< ok 0 0
> exec ^
: string
| select 1 as current_timestamp
< error
| 1064 42000 …near 'current_timestamp' at line 1
> exec ^
: string
| select 1 current_timestamp
< error
| 1064 42000 …near 'current_timestamp' at line 1
> exec ^
: string
| select `current_timestamp` from t
< error
| 1054 42S22 Unknown column 'current_timestamp' in 'field list'
> exec ^
: string
| select `current_timestamp`(-1)
< error
| 1582 42000 Incorrect parameter count in the call to native function 'current_timestamp'
> exec ^
: string
| select tp.current_timestamp(-1)
< error
| 1054 42S22 Unknown column '' in 'field list'
> exec ^
: string
| select t.current_timestamp from t
< error
| 1054 42S22 Unknown column 't.current_timestamp' in 'field list'
> exec ^
: string
| select current_timestamp() over ()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| select current_timestamp(1) over ()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| select current_time
< ok 0 0
> exec ^
: string
| select current_time()
< ok 0 0
> exec ^
: string
| select current_time(6)
< ok 0 0
> exec ^
: string
| select current_time(7)
< error
| 1426 42000 Too big precision 7 specified for column 'current_time'. Maximum is 6.
> exec ^
: string
| select current_time(-1)
< error
| 1064 42000 …near '-1)' at line 1
> exec ^
: string
| select current_time(1.0)
< error
| 1064 42000 …near '1.0)' at line 1
> exec ^
: string
| select current_time('1')
< error
| 1064 42000 …near ''1')' at line 1
> exec ^
: string
| select current_time(null)
< error
| 1064 42000 …near 'null)' at line 1
> exec ^
: string
| select current_time(1,2)
< error
| 1064 42000 …near ',2)' at line 1
> exec ^
: string
| select current_time(0x1)
< error
| 1064 42000 …near '0x1)' at line 1
> exec ^
: string
| select current_time (1)
< ok 0 0
> exec ^
: string
| select current_time (-1)
< error
| 1064 42000 …near '-1)' at line 1
> exec ^
: string
| select 1 as current_time
< error
| 1064 42000 …near 'current_time' at line 1
> exec ^
: string
| select `current_time`(-1)
< error
| 1305 42000 FUNCTION current_time does not exist
> exec ^
: string
| select tp.current_time(-1)
< error
| 1054 42S22 Unknown column '' in 'field list'
> exec ^
: string
| select curtime
< error
| 1054 42S22 Unknown column 'curtime' in 'field list'
> exec ^
: string
| select curtime()
< ok 0 0
> exec ^
: string
| select curtime(6)
< ok 0 0
> exec ^
: string
| select curtime(7)
< error
| 1426 42000 Too big precision 7 specified for column 'current_time'. Maximum is 6.
> exec ^
: string
| select curtime(-1)
< error
| 1064 42000 …near '-1)' at line 1
> exec ^
: string
| select curtime(1.0)
< error
| 1064 42000 …near '1.0)' at line 1
> exec ^
: string
| select curtime('1')
< error
| 1064 42000 …near ''1')' at line 1
> exec ^
: string
| select curtime(null)
< error
| 1064 42000 …near 'null)' at line 1
> exec ^
: string
| select curtime(1,2)
< error
| 1064 42000 …near ',2)' at line 1
> exec ^
: string
| select curtime(0x1)
< error
| 1064 42000 …near '0x1)' at line 1
> exec ^
: string
| select curtime(b)
< error
| 1064 42000 …near 'b)' at line 1
> exec ^
: string
| select curtime (1)
< ok 0 0
> exec ^
: string
| select curtime (-1)
< error
| 1064 42000 …near '-1)' at line 1
> exec ^
: string
| select curtime (b)
< error
| 1064 42000 …near 'b)' at line 1
> exec ^
: string
| select curtime /**/ (-1)
< error
| 1064 42000 …near '-1)' at line 1
> exec ^
: string
| select curtime from t
< ok 0 0
> exec ^
: string
| select curtime.a from curtime
< error
| 1146 42S02 Table 'omni_temporal.curtime' doesn't exist
> exec ^
: string
| select t.curtime from t
< ok 0 0
> exec ^
: string
| select 1 as curtime
< ok 0 0
> exec ^
: string
| select 1 curtime
< ok 0 0
> exec ^
: string
| select `curtime`(-1)
< error
| 1582 42000 Incorrect parameter count in the call to native function 'curtime'
> exec ^
: string
| select tp.curtime(-1)
< error
| 1054 42S22 Unknown column '' in 'field list'
> exec ^
: string
| select curtime(1) over ()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| select curtime() over ()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| select curtime(255)
< error
| 1426 42000 Too big precision 255 specified for column 'current_time'. Maximum is 6.
> exec ^
: string
| select curtime(256)
< error
| 1426 42000 Too big precision 256 specified for column 'current_time'. Maximum is 6.
> exec ^
: string
| select curtime(2147483648)
< ok 0 0
> exec ^
: string
| select curtime(18446744073709551616)
< error
| 1064 42000 …near '18446744073709551616)' at line 1
> exec ^
: string
| select utc_timestamp
< ok 0 0
> exec ^
: string
| select utc_timestamp()
< ok 0 0
> exec ^
: string
| select utc_timestamp(6)
< ok 0 0
> exec ^
: string
| select utc_timestamp(7)
< error
| 1426 42000 Too big precision 7 specified for column 'utc_timestamp'. Maximum is 6.
> exec ^
: string
| select utc_timestamp(-1)
< error
| 1064 42000 …near '-1)' at line 1
> exec ^
: string
| select utc_timestamp(1.0)
< error
| 1064 42000 …near '1.0)' at line 1
> exec ^
: string
| select utc_timestamp('1')
< error
| 1064 42000 …near ''1')' at line 1
> exec ^
: string
| select utc_timestamp(null)
< error
| 1064 42000 …near 'null)' at line 1
> exec ^
: string
| select utc_timestamp(1,2)
< error
| 1064 42000 …near ',2)' at line 1
> exec ^
: string
| select utc_timestamp (1)
< ok 0 0
> exec ^
: string
| select 1 as utc_timestamp
< error
| 1064 42000 …near 'utc_timestamp' at line 1
> exec ^
: string
| select `utc_timestamp`(-1)
< error
| 1582 42000 Incorrect parameter count in the call to native function 'utc_timestamp'
> exec ^
: string
| select tp.utc_timestamp(-1)
< error
| 1054 42S22 Unknown column '' in 'field list'
> exec ^
: string
| select utc_time
< ok 0 0
> exec ^
: string
| select utc_time()
< ok 0 0
> exec ^
: string
| select utc_time(6)
< ok 0 0
> exec ^
: string
| select utc_time(7)
< error
| 1426 42000 Too big precision 7 specified for column 'utc_time'. Maximum is 6.
> exec ^
: string
| select utc_time(-1)
< error
| 1064 42000 …near '-1)' at line 1
> exec ^
: string
| select utc_time(1.0)
< error
| 1064 42000 …near '1.0)' at line 1
> exec ^
: string
| select utc_time('1')
< error
| 1064 42000 …near ''1')' at line 1
> exec ^
: string
| select utc_time(null)
< error
| 1064 42000 …near 'null)' at line 1
> exec ^
: string
| select utc_time(1,2)
< error
| 1064 42000 …near ',2)' at line 1
> exec ^
: string
| select utc_time (1)
< ok 0 0
> exec ^
: string
| select 1 as utc_time
< error
| 1064 42000 …near 'utc_time' at line 1
> exec ^
: string
| select localtime
< ok 0 0
> exec ^
: string
| select localtime()
< ok 0 0
> exec ^
: string
| select localtime(6)
< ok 0 0
> exec ^
: string
| select localtime(7)
< error
| 1426 42000 Too big precision 7 specified for column 'now'. Maximum is 6.
> exec ^
: string
| select localtime(-1)
< error
| 1064 42000 …near '-1)' at line 1
> exec ^
: string
| select localtime(1.0)
< error
| 1064 42000 …near '1.0)' at line 1
> exec ^
: string
| select localtime('1')
< error
| 1064 42000 …near ''1')' at line 1
> exec ^
: string
| select localtime(null)
< error
| 1064 42000 …near 'null)' at line 1
> exec ^
: string
| select localtimestamp
< ok 0 0
> exec ^
: string
| select localtimestamp()
< ok 0 0
> exec ^
: string
| select localtimestamp(6)
< ok 0 0
> exec ^
: string
| select localtimestamp(-1)
< error
| 1064 42000 …near '-1)' at line 1
> exec ^
: string
| select localtimestamp(1.0)
< error
| 1064 42000 …near '1.0)' at line 1
> exec ^
: string
| select localtimestamp('1')
< error
| 1064 42000 …near ''1')' at line 1
> exec ^
: string
| select localtimestamp(null)
< error
| 1064 42000 …near 'null)' at line 1
> exec ^
: string
| select 1 as localtime
< error
| 1064 42000 …near 'localtime' at line 1
> exec ^
: string
| select now
< error
| 1054 42S22 Unknown column 'now' in 'field list'
> exec ^
: string
| select now()
< ok 0 0
> exec ^
: string
| select now(0)
< ok 0 0
> exec ^
: string
| select now(6)
< ok 0 0
> exec ^
: string
| select now(7)
< error
| 1426 42000 Too big precision 7 specified for column 'now'. Maximum is 6.
> exec ^
: string
| select now(255)
< error
| 1426 42000 Too big precision 255 specified for column 'now'. Maximum is 6.
> exec ^
: string
| select now(256)
< error
| 1426 42000 Too big precision 256 specified for column 'now'. Maximum is 6.
> exec ^
: string
| select now(-1)
< error
| 1064 42000 …near '-1)' at line 1
> exec ^
: string
| select now(1.0)
< error
| 1064 42000 …near '1.0)' at line 1
> exec ^
: string
| select now('1')
< error
| 1064 42000 …near ''1')' at line 1
> exec ^
: string
| select now(null)
< error
| 1064 42000 …near 'null)' at line 1
> exec ^
: string
| select now(b)
< error
| 1064 42000 …near 'b)' at line 1
> exec ^
: string
| select now(1,2)
< error
| 1064 42000 …near ',2)' at line 1
> exec ^
: string
| select now(0x1)
< error
| 1064 42000 …near '0x1)' at line 1
> exec ^
: string
| select now (1)
< ok 0 0
> exec ^
: string
| select now (-1)
< error
| 1064 42000 …near '-1)' at line 1
> exec ^
: string
| select now(*)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| select now(distinct 1)
< error
| 1064 42000 …near 'distinct 1)' at line 1
> exec ^
: string
| select now from t
< ok 0 0
> exec ^
: string
| select now.a from now
< error
| 1146 42S02 Table 'omni_temporal.now' doesn't exist
> exec ^
: string
| select t.now from t
< ok 0 0
> exec ^
: string
| select 1 as now
< ok 0 0
> exec ^
: string
| select 1 now
< ok 0 0
> exec ^
: string
| select `now`(-1)
< error
| 1305 42000 FUNCTION now does not exist
> exec ^
: string
| select tp.now(-1)
< error
| 1054 42S22 Unknown column '' in 'field list'
> exec ^
: string
| select now() over ()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| select now(1) over ()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| select sysdate
< error
| 1054 42S22 Unknown column 'sysdate' in 'field list'
> exec ^
: string
| select sysdate()
< ok 0 0
> exec ^
: string
| select sysdate(6)
< ok 0 0
> exec ^
: string
| select sysdate(7)
< error
| 1426 42000 Too big precision 7 specified for column 'sysdate'. Maximum is 6.
> exec ^
: string
| select sysdate(-1)
< error
| 1064 42000 …near '-1)' at line 1
> exec ^
: string
| select sysdate(1.0)
< error
| 1064 42000 …near '1.0)' at line 1
> exec ^
: string
| select sysdate('1')
< error
| 1064 42000 …near ''1')' at line 1
> exec ^
: string
| select sysdate(null)
< error
| 1064 42000 …near 'null)' at line 1
> exec ^
: string
| select sysdate(1,2)
< error
| 1064 42000 …near ',2)' at line 1
> exec ^
: string
| select sysdate(0x1)
< error
| 1064 42000 …near '0x1)' at line 1
> exec ^
: string
| select sysdate (-1)
< error
| 1064 42000 …near '-1)' at line 1
> exec ^
: string
| select sysdate from t
< ok 0 0
> exec ^
: string
| select 1 as sysdate
< ok 0 0
> exec ^
: string
| select tp.sysdate(-1)
< error
| 1054 42S22 Unknown column '' in 'field list'
> exec ^
: string
| select `sysdate`(-1)
< error
| 1582 42000 Incorrect parameter count in the call to native function 'sysdate'
> exec ^
: string
| select sysdate() over ()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| select curdate
< error
| 1054 42S22 Unknown column 'curdate' in 'field list'
> exec ^
: string
| select curdate()
< ok 0 0
> exec ^
: string
| select curdate(1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| select curdate(0)
< error
| 1064 42000 …near '0)' at line 1
> exec ^
: string
| select curdate(null)
< error
| 1064 42000 …near 'null)' at line 1
> exec ^
: string
| select curdate(b)
< error
| 1064 42000 …near 'b)' at line 1
> exec ^
: string
| select curdate (1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| select curdate (b)
< error
| 1064 42000 …near 'b)' at line 1
> exec ^
: string
| select curdate from t
< ok 0 0
> exec ^
: string
| select 1 as curdate
< ok 0 0
> exec ^
: string
| select `curdate`(1)
< error
| 1305 42000 FUNCTION curdate does not exist
> exec ^
: string
| select tp.curdate(1)
< error
| 1054 42S22 Unknown column '' in 'field list'
> exec ^
: string
| select curdate() over ()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| select current_date
< ok 0 0
> exec ^
: string
| select current_date()
< ok 0 0
> exec ^
: string
| select current_date(1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| select current_date(null)
< error
| 1064 42000 …near 'null)' at line 1
> exec ^
: string
| select current_date (1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| select current_date ()
< ok 0 0
> exec ^
: string
| select 1 as current_date
< error
| 1064 42000 …near 'current_date' at line 1
> exec ^
: string
| select tp.current_date(1)
< error
| 1054 42S22 Unknown column '' in 'field list'
> exec ^
: string
| select `current_date`(1)
< error
| 1305 42000 FUNCTION current_date does not exist
> exec ^
: string
| select current_date() over ()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| select utc_date
< ok 0 0
> exec ^
: string
| select utc_date()
< ok 0 0
> exec ^
: string
| select utc_date(1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| select utc_date(null)
< error
| 1064 42000 …near 'null)' at line 1
> exec ^
: string
| select utc_date (1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| select 1 as utc_date
< error
| 1064 42000 …near 'utc_date' at line 1
> exec ^
: string
| select current_user
< ok 0 0
> exec ^
: string
| select current_user()
< ok 0 0
> exec ^
: string
| select current_user(1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| select current_user (1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| select current_role
< error
| 1054 42S22 Unknown column 'current_role' in 'field list'
> exec ^
: string
| select current_role()
< ok 0 0
> exec ^
: string
| SELECT CURRENT_DATE, CURRENT_DATE(), CURDATE(1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| SELECT CURRENT_DATE, CURRENT_DATE(), CURDATE()
< ok 0 0
> exec ^
: string
| select date_add
< error
| 1054 42S22 Unknown column 'date_add' in 'field list'
> exec ^
: string
| select date_add from t
< ok 0 0
> exec ^
: string
| select date_add()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select date_add(a)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select date_add(a, 10)
< error
| 1064 42000 …near '10)' at line 1
> exec ^
: string
| select date_add(a, 0.10)
< error
| 1064 42000 …near '0.10)' at line 1
> exec ^
: string
| select date_add(a, "11,11")
< error
| 1064 42000 …near '"11,11")' at line 1
> exec ^
: string
| select date_add(a, b)
< error
| 1064 42000 …near 'b)' at line 1
> exec ^
: string
| select date_add(a, null)
< error
| 1064 42000 …near 'null)' at line 1
> exec ^
: string
| select date_add(a, interval 1 day)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select date_add(a, INTERVAL 1 SECOND_MICROSECOND)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select date_add(a, interval b day)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select date_add(a, interval 1+1 day)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select date_add(a, interval (1) day)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select date_add(a, interval interval(1,2) day)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select date_add(a, interval 1)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select date_add(a, interval 1 foo)
< error
| 1064 42000 …near 'foo)' at line 1
> exec ^
: string
| select date_add(a, interval 1 day + 1)
< error
| 1064 42000 …near '+ 1)' at line 1
> exec ^
: string
| select date_add(a, interval 1 day + interval 1 day)
< error
| 1064 42000 …near '+ interval 1 day)' at line 1
> exec ^
: string
| select date_add(a, 1 + interval 1 day)
< error
| 1064 42000 …near '1 + interval 1 day)' at line 1
> exec ^
: string
| select date_add(a, (interval 1 day))
< error
| 1064 42000 …near '(interval 1 day))' at line 1
> exec ^
: string
| select date_add(a, b + interval 1 day)
< error
| 1064 42000 …near 'b + interval 1 day)' at line 1
> exec ^
: string
| select date_add(a, interval 1 day, 1)
< error
| 1064 42000 …near ', 1)' at line 1
> exec ^
: string
| select date_add(a, 1, 2)
< error
| 1064 42000 …near '1, 2)' at line 1
> exec ^
: string
| select date_add(a, interval 1 sql_tsi_day)
< error
| 1064 42000 …near 'sql_tsi_day)' at line 1
> exec ^
: string
| select date_add(interval 1 day + a, interval 1 day)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select date_add(a + interval 1 day, interval 1 day)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select date_add(interval 1 day, interval 1 day)
< error
| 1064 42000 …near ', interval 1 day)' at line 1
> exec ^
: string
| select date_add(*, interval 1 day)
< error
| 1064 42000 …near '*, interval 1 day)' at line 1
> exec ^
: string
| select date_add(distinct a, interval 1 day)
< error
| 1064 42000 …near 'distinct a, interval 1 day)' at line 1
> exec ^
: string
| select date_add(all a, interval 1 day)
< error
| 1064 42000 …near 'all a, interval 1 day)' at line 1
> exec ^
: string
| select date_add(a, interval 1 day) over ()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| select date_add (a, 10)
< error
| 1064 42000 …near '10)' at line 1
> exec ^
: string
| select date_add (a, b)
< error
| 1064 42000 …near 'b)' at line 1
> exec ^
: string
| select date_add (a, interval 1 day)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select date_add /**/ (a, 10)
< error
| 1064 42000 …near '10)' at line 1
> exec ^
: string
| select date_add.a from date_add
< error
| 1146 42S02 Table 'omni_temporal.date_add' doesn't exist
> exec ^
: string
| select t.date_add from t
< ok 0 0
> exec ^
: string
| select 1 as date_add
< ok 0 0
> exec ^
: string
| select 1 date_add
< ok 0 0
> exec ^
: string
| select `date_add`(a, 10)
< error
| 1582 42000 Incorrect parameter count in the call to native function 'date_add'
> exec ^
: string
| select `date_add`(a, interval 1 day)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select tp.date_add(a, 10)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select tp.date_add(a, interval 1 day)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select date_add("2011-11-11 10:10:10.123456", 10)
< error
| 1064 42000 …near '10)' at line 1
> exec ^
: string
| select date_add("2011-11-11 10:10:10.123456", 0.10)
< error
| 1064 42000 …near '0.10)' at line 1
> exec ^
: string
| select date_add("2011-11-11 10:10:10.123456", "11,11")
< error
| 1064 42000 …near '"11,11")' at line 1
> exec ^
: string
| select date_add("2011-11-11 10:10:10.123456", interval 10 microsecond)
< ok 0 0
> exec ^
: string
| select date_sub
< error
| 1054 42S22 Unknown column 'date_sub' in 'field list'
> exec ^
: string
| select date_sub from t
< ok 0 0
> exec ^
: string
| select date_sub()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select date_sub(a)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select date_sub(a, 10)
< error
| 1064 42000 …near '10)' at line 1
> exec ^
: string
| select date_sub(a, 0.10)
< error
| 1064 42000 …near '0.10)' at line 1
> exec ^
: string
| select date_sub(a, "11,11")
< error
| 1064 42000 …near '"11,11")' at line 1
> exec ^
: string
| select date_sub(a, b)
< error
| 1064 42000 …near 'b)' at line 1
> exec ^
: string
| select date_sub(a, null)
< error
| 1064 42000 …near 'null)' at line 1
> exec ^
: string
| select date_sub(a, interval 1 day)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select date_sub(a, interval 1)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select date_sub(a, interval 1 day + 1)
< error
| 1064 42000 …near '+ 1)' at line 1
> exec ^
: string
| select date_sub(a, 1 + interval 1 day)
< error
| 1064 42000 …near '1 + interval 1 day)' at line 1
> exec ^
: string
| select date_sub(a, (interval 1 day))
< error
| 1064 42000 …near '(interval 1 day))' at line 1
> exec ^
: string
| select date_sub(a, interval 1 day, 1)
< error
| 1064 42000 …near ', 1)' at line 1
> exec ^
: string
| select date_sub(a, interval 1 day) over ()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| select date_sub (a, 10)
< error
| 1064 42000 …near '10)' at line 1
> exec ^
: string
| select date_sub (a, interval 1 day)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select 1 as date_sub
< ok 0 0
> exec ^
: string
| select `date_sub`(a, 10)
< error
| 1582 42000 Incorrect parameter count in the call to native function 'date_sub'
> exec ^
: string
| select tp.date_sub(a, 10)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select adddate
< error
| 1054 42S22 Unknown column 'adddate' in 'field list'
> exec ^
: string
| select adddate from t
< ok 0 0
> exec ^
: string
| select adddate()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select adddate(a)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select adddate(a, 10)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select adddate(a, 0.10)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select adddate(a, "11,11")
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select adddate(a, b)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select adddate(a, interval 1 day)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select adddate(a, interval 1)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select adddate(a, interval 1 day + 1)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select adddate(a, b + interval 1 day)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select adddate(a, (interval 1 day))
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| select adddate(a, 1, 2)
< error
| 1064 42000 …near ', 2)' at line 1
> exec ^
: string
| select adddate(a, interval 1 day) over ()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| select adddate (a, 10)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select adddate (a)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select adddate (a, interval 1 day)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select 1 as adddate
< ok 0 0
> exec ^
: string
| select `adddate`(a, 10)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select subdate(a)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select subdate(a, 10)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select subdate(a, interval 1 day)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select subdate(a, 1, 2)
< error
| 1064 42000 …near ', 2)' at line 1
> exec ^
: string
| select subdate(a, interval 1 day + 1)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select subdate (a, interval 1 day)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select now ()
< ok 0 0
> exec ^
: string
| select curdate ()
< ok 0 0
> exec ^
: string
| select curtime ()
< ok 0 0
> exec ^
: string
| select sysdate ()
< ok 0 0
> exec ^
: string
| select sysdate (1)
< ok 0 0
> exec ^
: string
| select sysdate (b)
< error
| 1064 42000 …near 'b)' at line 1
> exec ^
: string
| select sysdate (1,2)
< error
| 1064 42000 …near ',2)' at line 1
> exec ^
: string
| select sysdate (1) over ()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| select now (1) over ()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| select curtime (1) over ()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| select date_add ()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select adddate ()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select subdate ()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select adddate (a, b, c)
< error
| 1064 42000 …near ', c)' at line 1
> exec ^
: string
| select curdate/**/(1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| select now/**/(1)
< ok 0 0
> exec ^
: string
| select CURDATE (1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| select cUrTiMe(-1)
< error
| 1064 42000 …near '-1)' at line 1
> exec ^
: string
| select t.current_time(-1)
< error
| 1305 42000 FUNCTION current_time does not exist
> exec ^
: string
| select t.now(-1)
< error
| 1305 42000 FUNCTION now does not exist
> exec ^
: string
| select t.curdate(1)
< error
| 1305 42000 FUNCTION curdate does not exist
> exec ^
: string
| select date_add(a, interval(1,2))
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select date_sub(a, interval(1,2))
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select adddate(a, interval(1,2))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select subdate(a, interval(1,2))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select date_add((interval 1 day), interval 1 day)
< error
| 1064 42000 …near '), interval 1 day)' at line 1
> exec ^
: string
| select adddate(interval 1 day, 1)
< error
| 1064 42000 …near ', 1)' at line 1
> exec ^
: string
| select adddate((interval 1 day), 1)
< error
| 1064 42000 …near '), 1)' at line 1
> exec ^
: string
| select adddate(interval 1 day + a, interval 1 day)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select adddate(*, 1)
< error
| 1064 42000 …near '*, 1)' at line 1
> exec ^
: string
| select adddate(a, *)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| select adddate(distinct a, 1)
< error
| 1064 42000 …near 'distinct a, 1)' at line 1
> exec ^
: string
| select adddate(a, distinct 1)
< error
| 1064 42000 …near 'distinct 1)' at line 1
> exec ^
: string
| select date_add(a, distinct interval 1 day)
< error
| 1064 42000 …near 'distinct interval 1 day)' at line 1
> exec ^
: string
| select date_add(a, interval 1 day,)
< error
| 1064 42000 …near ',)' at line 1
> exec ^
: string
| select abs(interval 1 day)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select concat(interval 1 day)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select coalesce(interval 1 day, 1)
< error
| 1064 42000 …near ', 1)' at line 1
> exec ^
: string
| select interval 1 day + 1
< ok 0 0
> exec ^
: string
| select 1 + interval 1 day
< ok 0 0
> exec ^
: string
| select abs((interval 1 day))
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| select abs(interval 1 day + 1)
< ok 0 0
> exec ^
: string
| select current_time(0xffffffffffffffffffff)
< error
| 1064 42000 …near '0xffffffffffffffffffff)' at line 1
> exec ^
: string
| select current_time(0X1)
< error
| 1064 42000 …near '0X1)' at line 1
> exec ^
: string
| select current_time(0b1)
< error
| 1064 42000 …near '0b1)' at line 1
> exec ^
: string
| select now(0b1)
< error
| 1064 42000 …near '0b1)' at line 1
> exec ^
: string
| select curtime(0xa)
< error
| 1064 42000 …near '0xa)' at line 1
> exec ^
: string
| select utc_time(0x1)
< error
| 1064 42000 …near '0x1)' at line 1
> exec ^
: string
| select current_time(0002147483647)
< ok 0 0
> exec ^
: string
| select current_time(2147483647)
< ok 0 0
> exec ^
: string
| select current_time(18446744073709551615)
< ok 0 0
> exec ^
: string
| select current_time(18446744073709551616)
< error
| 1064 42000 …near '18446744073709551616)' at line 1
> exec ^
: string
| select now(2147483647)
< ok 0 0
> exec ^
: string
| select now(2147483648)
< ok 0 0
> exec ^
: string
| select now(18446744073709551615)
< ok 0 0
> exec ^
: string
| select now(18446744073709551616)
< error
| 1064 42000 …near '18446744073709551616)' at line 1
> exec ^
: string
| select now(99999999999999999999999)
< error
| 1064 42000 …near '99999999999999999999999)' at line 1
> exec ^
: string
| select curtime(2147483647)
< ok 0 0
> exec ^
: string
| select curtime(18446744073709551615)
< ok 0 0
> exec ^
: string
| select sysdate(2147483647)
< ok 0 0
> exec ^
: string
| select sysdate(2147483648)
< ok 0 0
> exec ^
: string
| select sysdate(18446744073709551615)
< ok 0 0
> exec ^
: string
| select sysdate(18446744073709551616)
< error
| 1064 42000 …near '18446744073709551616)' at line 1
> exec ^
: string
| select utc_time(2147483648)
< ok 0 0
> exec ^
: string
| select utc_timestamp(2147483648)
< ok 0 0
> exec ^
: string
| select localtime(2147483648)
< ok 0 0
> exec ^
: string
| select localtimestamp(2147483648)
< ok 0 0
> exec ^
: string
| create table t1 (a timestamp default current_timestamp(6) on update current_timestamp(6))
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| create table t2 (a timestamp default current_timestamp(-1))
< error
| 1064 42000 …near '-1))' at line 1
> exec ^
: string
| create table t3 (a timestamp default now(6) on update now(6))
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| create table t4 (a timestamp default now(-1))
< error
| 1064 42000 …near '-1))' at line 1
> exec ^
: string
| create table t5 (a timestamp on update current_timestamp(1.0))
< error
| 1064 42000 …near '1.0))' at line 1
> exec ^
: string
| create table t6 (a date default curdate(1))
< error
| 1064 42000 …near 'curdate(1))' at line 1
> exec ^
: string
| create table t7 (a date default (curdate(1)))
< error
| 1064 42000 …near '1)))' at line 1
> exec ^
: string
| create table t8 (a datetime default (date_add(now(), 1)))
< error
| 1064 42000 …near '1)))' at line 1
> exec ^
: string
| create table t9 (a datetime default (date_add(now(), interval 1 day)))
< ok 0 0
> exec ^
: string
| create table t10 (a datetime default current_timestamp(7))
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| create table t11 (a datetime default now(7))
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| select current_timestamp(00000000000000000000000018446744073709551615)
< ok 0 0
> exec ^
: string
| select current_timestamp(00000000000000000000000018446744073709551616)
< error
| 1064 42000 …near '00000000000000000000000018446744073709551616)' at line 1
> exec ^
: string
| select current_timestamp(0000000000002147483647)
< ok 0 0
> exec ^
: string
| select current_timestamp(0000000000002147483648)
< ok 0 0
> exec ^
: string
| select now(0000000000002147483648)
< ok 0 0
> exec ^
: string
| select current_time(x'0xaa')
< error
| 1064 42000 …near 'x'0xaa')' at line 1
> exec ^
: string
| select now(x'0xaa')
< error
| 1064 42000 …near 'x'0xaa')' at line 1
> exec ^
: string
| select current_date (
< error
| 1064 42000 …near '(' at line 1
> exec ^
: string
| select now(1)(2)
< error
| 1064 42000 …near '(2)' at line 1
> exec ^
: string
| select curtime(1)(2)
< error
| 1064 42000 …near '(2)' at line 1
> exec ^
: string
| select current_time(1)(2)
< error
| 1064 42000 …near '(2)' at line 1
> exec ^
: string
| select now(1) x
< ok 0 0
> exec ^
: string
| select curtime(1) as x from t
< ok 0 0
> exec ^
: string
| select now() + interval 1 day
< ok 0 0
> exec ^
: string
| select curdate() - interval 1 day
< ok 0 0
> exec ^
: string
| select date_add(a, interval 1 day) x from t
< ok 0 0
> exec ^
: string
| select date_sub(a, interval 1 day) as x from t
< ok 0 0
> exec ^
: string
| select date_add(a, interval 1 day) + 1 from t
< ok 0 0
> exec ^
: string
| select date_add(a,
< error
| 1064 42000 …near ',' at line 1
> exec ^
: string
| select date_add(a
< error
| 1064 42000 …near 'a' at line 1
> exec ^
: string
| select date_add(a, interval 1 day
< error
| 1064 42000 …near 'day' at line 1
> exec ^
: string
| select date_add(, interval 1 day)
< error
| 1064 42000 …near ', interval 1 day)' at line 1
> exec ^
: string
| select adddate(a,)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select adddate(,a)
< error
| 1064 42000 …near ',a)' at line 1
> exec ^
: string
| select adddate(a, 1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| select interval(interval 1 day, 1)
< error
| 1064 42000 …near ', 1)' at line 1
> exec ^
: string
| select if(interval 1 day, 1, 2)
< error
| 1064 42000 …near ', 1, 2)' at line 1
> exec ^
: string
| select left(interval 1 day, 1)
< error
| 1064 42000 …near ', 1)' at line 1
> exec ^
: string
| select count(interval 1 day)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select sum(interval 1 day)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select max(a + interval 1 day) from t
< ok 0 0
> exec ^
: string
| select date_format(now(), '%Y')
< ok 0 0
> exec ^
: string
| select date_format(a + interval 1 day, '%Y') from t
< ok 0 0
> exec ^
: string
| select extract(day from a + interval 1 day) from t
< ok 0 0
> exec ^
: string
| select timestampadd(day, 1, now(6))
< ok 0 0
> exec ^
: string
| select now(now())
< error
| 1064 42000 …near 'now())' at line 1
> exec ^
: string
| select curtime(now())
< error
| 1064 42000 …near 'now())' at line 1
> exec ^
: string
| select date_add(now(6), interval 1 day)
< ok 0 0
> exec ^
: string
| select date_add(date_add(a, interval 1 day), interval 1 day) from t
< ok 0 0
> exec ^
: string
| select date_add(a, interval date_add(b, interval 1 day) day) from t
< ok 0 0
> exec ^
: string
| select adddate(adddate(a, 1), 1) from t
< ok 0 0
> exec ^
: string
| select date_sub(curdate(), interval 1 month)
< ok 0 0
> exec ^
: string
| select sysdate(3) + interval 1 day
< ok 0 0
> exec ^
: string
| select utc_timestamp(3) + 1
< ok 0 0
> exec ^
: string
| update t set a = now(6)
< ok 0 0
> exec ^
: string
| update t set a = now(-1)
< error
| 1064 42000 …near '-1)' at line 1
> exec ^
: string
| update t set a = date_add(a, interval 1 day)
< ok 0 0
> exec ^
: string
| update t set a = date_add(a, 1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| delete from t where a < date_sub(now(), interval 1 day)
< ok 0 0
> exec ^
: string
| delete from t where a < date_sub(now(), 1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| insert into t (a) values (current_timestamp(3))
< ok 0 0
> exec ^
: string
| insert into t (a) values (current_timestamp(-1))
< error
| 1064 42000 …near '-1))' at line 1
> exec ^
: string
| select * from t where a > curdate() order by now(6)
< ok 0 0
> exec ^
: string
| select * from t where a > curdate(1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| select bit_and(interval 1 day)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select json_arrayagg(interval 1 day)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select char(interval 1 day)
< error
| 1064 42000 …near ')' at line 1
> exec
| DROP DATABASE IF EXISTS omni_temporal
< ok 1 0
