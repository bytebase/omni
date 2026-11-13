-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestTiDBTemporalFormsMatchEngine default#1
> exec
| CREATE DATABASE IF NOT EXISTS omni_temporal
< ok 0 0
> exec
| USE omni_temporal
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS t (a DATETIME, b INT, c INT, `now` INT, `curdate` INT, `curtime` INT, `date_add` INT, `date_sub` INT, `sysdate` INT, `adddate` INT, `subdate` INT)
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp(0)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp(6)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp(7)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1426 42000 Too-big precision 7 specified for 'current_timestamp'. Maximum is 6.
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp(255)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1426 42000 Too-big precision 255 specified for 'current_timestamp'. Maximum is 6.
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp(256)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1426 42000 Too-big precision 256 specified for 'current_timestamp'. Maximum is 6.
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp(2147483647)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1426 42000 Too-big precision 2147483647 specified for 'current_timestamp'. Maximum is 6.
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp(2147483648)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp(4294967296)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp(18446744073709551615)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp(18446744073709551616)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 45 near "18446744073709551616)" 
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp(00000000000000000000001)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp(-1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 26 near "-1)" 
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp(+1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 26 near "+1)" 
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp(1.0)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 28 near "1.0)" 
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp(1e0)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 28 near "1e0)" 
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp(.5)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 27 near ".5)" 
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp('2')
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 28 near "'2')" 
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp(null)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 29 near "null)" 
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp(true)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 29 near "true)" 
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp(0x1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 28 near "0x1)" 
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp(x'01')
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 30 near "x'01')" 
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp(b'1')
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 29 near "b'1')" 
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp(b)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 26 near "b)" 
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp(@a)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 27 near "@a)" 
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp(?)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 26 near "?)" 
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp((1))
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 26 near "(1))" 
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp(1+1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 27 near "+1)" 
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp(1,2)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 27 near ",2)" 
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp(1,)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 27 near ",)" 
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp(
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 25 near "" 
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp (1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp /* c */ (1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp (-1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 27 near "-1)" 
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp.a from t current_timestamp
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 51 near "current_timestamp" 
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp a
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp as a
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select 1 as current_timestamp
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 29 near "current_timestamp" 
> exec
| SET @omni_temporal = ?
: string
| select 1 current_timestamp
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 26 near "current_timestamp" 
> exec
| SET @omni_temporal = ?
: string
| select `current_timestamp` from t
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'current_timestamp' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select `current_timestamp`(-1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use
> exec
| SET @omni_temporal = ?
: string
| select tp.current_timestamp(-1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1305 42000 FUNCTION tp.current_timestamp does not exist
> exec
| SET @omni_temporal = ?
: string
| select t.current_timestamp from t
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 't.current_timestamp' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp() over ()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 31 near "over ()" 
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp(1) over ()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 32 near "over ()" 
> exec
| SET @omni_temporal = ?
: string
| select current_time
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select current_time()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select current_time(6)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select current_time(7)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1426 42000 Too-big precision 7 specified for 'current_time'. Maximum is 6.
> exec
| SET @omni_temporal = ?
: string
| select current_time(-1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 21 near "-1)" 
> exec
| SET @omni_temporal = ?
: string
| select current_time(1.0)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 23 near "1.0)" 
> exec
| SET @omni_temporal = ?
: string
| select current_time('1')
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 23 near "'1')" 
> exec
| SET @omni_temporal = ?
: string
| select current_time(null)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 24 near "null)" 
> exec
| SET @omni_temporal = ?
: string
| select current_time(1,2)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 22 near ",2)" 
> exec
| SET @omni_temporal = ?
: string
| select current_time(0x1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 23 near "0x1)" 
> exec
| SET @omni_temporal = ?
: string
| select current_time (1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select current_time (-1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 22 near "-1)" 
> exec
| SET @omni_temporal = ?
: string
| select 1 as current_time
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 24 near "current_time" 
> exec
| SET @omni_temporal = ?
: string
| select `current_time`(-1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use
> exec
| SET @omni_temporal = ?
: string
| select tp.current_time(-1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1305 42000 FUNCTION tp.current_time does not exist
> exec
| SET @omni_temporal = ?
: string
| select curtime
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'curtime' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select curtime()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select curtime(6)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select curtime(7)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1426 42000 Too-big precision 7 specified for 'curtime'. Maximum is 6.
> exec
| SET @omni_temporal = ?
: string
| select curtime(-1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 16 near "-1)" 
> exec
| SET @omni_temporal = ?
: string
| select curtime(1.0)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 18 near "1.0)" 
> exec
| SET @omni_temporal = ?
: string
| select curtime('1')
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 18 near "'1')" 
> exec
| SET @omni_temporal = ?
: string
| select curtime(null)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 19 near "null)" 
> exec
| SET @omni_temporal = ?
: string
| select curtime(1,2)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 17 near ",2)" 
> exec
| SET @omni_temporal = ?
: string
| select curtime(0x1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 18 near "0x1)" 
> exec
| SET @omni_temporal = ?
: string
| select curtime(b)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 16 near "b)" 
> exec
| SET @omni_temporal = ?
: string
| select curtime (1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 16 near "(1)" 
> exec
| SET @omni_temporal = ?
: string
| select curtime (-1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 16 near "(-1)" 
> exec
| SET @omni_temporal = ?
: string
| select curtime (b)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 16 near "(b)" 
> exec
| SET @omni_temporal = ?
: string
| select curtime /**/ (-1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 21 near "(-1)" 
> exec
| SET @omni_temporal = ?
: string
| select curtime from t
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select curtime.a from curtime
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1146 42S02 Table 'omni_temporal.curtime' doesn't exist
> exec
| SET @omni_temporal = ?
: string
| select t.curtime from t
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select 1 as curtime
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select 1 curtime
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select `curtime`(-1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use
> exec
| SET @omni_temporal = ?
: string
| select tp.curtime(-1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1305 42000 FUNCTION tp.curtime does not exist
> exec
| SET @omni_temporal = ?
: string
| select curtime(1) over ()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 22 near "over ()" 
> exec
| SET @omni_temporal = ?
: string
| select curtime() over ()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 21 near "over ()" 
> exec
| SET @omni_temporal = ?
: string
| select curtime(255)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1426 42000 Too-big precision 255 specified for 'curtime'. Maximum is 6.
> exec
| SET @omni_temporal = ?
: string
| select curtime(256)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1426 42000 Too-big precision 256 specified for 'curtime'. Maximum is 6.
> exec
| SET @omni_temporal = ?
: string
| select curtime(2147483648)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use
> exec
| SET @omni_temporal = ?
: string
| select curtime(18446744073709551616)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 35 near "18446744073709551616)" 
> exec
| SET @omni_temporal = ?
: string
| select utc_timestamp
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select utc_timestamp()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select utc_timestamp(6)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select utc_timestamp(7)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1426 42000 Too-big precision 7 specified for 'utc_timestamp'. Maximum is 6.
> exec
| SET @omni_temporal = ?
: string
| select utc_timestamp(-1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 22 near "-1)" 
> exec
| SET @omni_temporal = ?
: string
| select utc_timestamp(1.0)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 24 near "1.0)" 
> exec
| SET @omni_temporal = ?
: string
| select utc_timestamp('1')
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 24 near "'1')" 
> exec
| SET @omni_temporal = ?
: string
| select utc_timestamp(null)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 25 near "null)" 
> exec
| SET @omni_temporal = ?
: string
| select utc_timestamp(1,2)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 23 near ",2)" 
> exec
| SET @omni_temporal = ?
: string
| select utc_timestamp (1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select 1 as utc_timestamp
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 25 near "utc_timestamp" 
> exec
| SET @omni_temporal = ?
: string
| select `utc_timestamp`(-1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use
> exec
| SET @omni_temporal = ?
: string
| select tp.utc_timestamp(-1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1305 42000 FUNCTION tp.utc_timestamp does not exist
> exec
| SET @omni_temporal = ?
: string
| select utc_time
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select utc_time()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select utc_time(6)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select utc_time(7)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1426 42000 Too-big precision 7 specified for 'utc_time'. Maximum is 6.
> exec
| SET @omni_temporal = ?
: string
| select utc_time(-1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 17 near "-1)" 
> exec
| SET @omni_temporal = ?
: string
| select utc_time(1.0)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 19 near "1.0)" 
> exec
| SET @omni_temporal = ?
: string
| select utc_time('1')
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 19 near "'1')" 
> exec
| SET @omni_temporal = ?
: string
| select utc_time(null)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 20 near "null)" 
> exec
| SET @omni_temporal = ?
: string
| select utc_time(1,2)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 18 near ",2)" 
> exec
| SET @omni_temporal = ?
: string
| select utc_time (1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select 1 as utc_time
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 20 near "utc_time" 
> exec
| SET @omni_temporal = ?
: string
| select localtime
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select localtime()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select localtime(6)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select localtime(7)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1426 42000 Too-big precision 7 specified for 'localtime'. Maximum is 6.
> exec
| SET @omni_temporal = ?
: string
| select localtime(-1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 18 near "-1)" 
> exec
| SET @omni_temporal = ?
: string
| select localtime(1.0)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 20 near "1.0)" 
> exec
| SET @omni_temporal = ?
: string
| select localtime('1')
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 20 near "'1')" 
> exec
| SET @omni_temporal = ?
: string
| select localtime(null)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 21 near "null)" 
> exec
| SET @omni_temporal = ?
: string
| select localtimestamp
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select localtimestamp()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select localtimestamp(6)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select localtimestamp(-1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 23 near "-1)" 
> exec
| SET @omni_temporal = ?
: string
| select localtimestamp(1.0)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 25 near "1.0)" 
> exec
| SET @omni_temporal = ?
: string
| select localtimestamp('1')
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 25 near "'1')" 
> exec
| SET @omni_temporal = ?
: string
| select localtimestamp(null)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 26 near "null)" 
> exec
| SET @omni_temporal = ?
: string
| select 1 as localtime
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 21 near "localtime" 
> exec
| SET @omni_temporal = ?
: string
| select now
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'now' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select now()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select now(0)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select now(6)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select now(7)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1426 42000 Too-big precision 7 specified for 'now'. Maximum is 6.
> exec
| SET @omni_temporal = ?
: string
| select now(255)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1426 42000 Too-big precision 255 specified for 'now'. Maximum is 6.
> exec
| SET @omni_temporal = ?
: string
| select now(256)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1426 42000 Too-big precision 256 specified for 'now'. Maximum is 6.
> exec
| SET @omni_temporal = ?
: string
| select now(-1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use
> exec
| SET @omni_temporal = ?
: string
| select now(1.0)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select now('1')
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select now(null)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select now(b)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'b' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select now(1,2)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1582 42000 Incorrect parameter count in the call to native function 'now'
> exec
| SET @omni_temporal = ?
: string
| select now(0x1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select now (1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 12 near "(1)" 
> exec
| SET @omni_temporal = ?
: string
| select now (-1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 12 near "(-1)" 
> exec
| SET @omni_temporal = ?
: string
| select now(*)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 12 near "*)" 
> exec
| SET @omni_temporal = ?
: string
| select now(distinct 1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 19 near "distinct 1)" 
> exec
| SET @omni_temporal = ?
: string
| select now from t
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select now.a from now
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1146 42S02 Table 'omni_temporal.now' doesn't exist
> exec
| SET @omni_temporal = ?
: string
| select t.now from t
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select 1 as now
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select 1 now
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select `now`(-1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use
> exec
| SET @omni_temporal = ?
: string
| select tp.now(-1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1305 42000 FUNCTION tp.now does not exist
> exec
| SET @omni_temporal = ?
: string
| select now() over ()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 17 near "over ()" 
> exec
| SET @omni_temporal = ?
: string
| select now(1) over ()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 18 near "over ()" 
> exec
| SET @omni_temporal = ?
: string
| select sysdate
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'sysdate' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select sysdate()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select sysdate(6)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select sysdate(7)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1426 42000 Too-big precision 7 specified for 'sysdate'. Maximum is 6.
> exec
| SET @omni_temporal = ?
: string
| select sysdate(-1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 16 near "-1)" 
> exec
| SET @omni_temporal = ?
: string
| select sysdate(1.0)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 18 near "1.0)" 
> exec
| SET @omni_temporal = ?
: string
| select sysdate('1')
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 18 near "'1')" 
> exec
| SET @omni_temporal = ?
: string
| select sysdate(null)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 19 near "null)" 
> exec
| SET @omni_temporal = ?
: string
| select sysdate(1,2)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 17 near ",2)" 
> exec
| SET @omni_temporal = ?
: string
| select sysdate(0x1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 18 near "0x1)" 
> exec
| SET @omni_temporal = ?
: string
| select sysdate (-1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use
> exec
| SET @omni_temporal = ?
: string
| select sysdate from t
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select 1 as sysdate
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select tp.sysdate(-1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1305 42000 FUNCTION tp.sysdate does not exist
> exec
| SET @omni_temporal = ?
: string
| select `sysdate`(-1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use
> exec
| SET @omni_temporal = ?
: string
| select sysdate() over ()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 21 near "over ()" 
> exec
| SET @omni_temporal = ?
: string
| select curdate
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'curdate' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select curdate()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select curdate(1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 16 near "1)" 
> exec
| SET @omni_temporal = ?
: string
| select curdate(0)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 16 near "0)" 
> exec
| SET @omni_temporal = ?
: string
| select curdate(null)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 19 near "null)" 
> exec
| SET @omni_temporal = ?
: string
| select curdate(b)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 16 near "b)" 
> exec
| SET @omni_temporal = ?
: string
| select curdate (1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 16 near "(1)" 
> exec
| SET @omni_temporal = ?
: string
| select curdate (b)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 16 near "(b)" 
> exec
| SET @omni_temporal = ?
: string
| select curdate from t
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select 1 as curdate
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select `curdate`(1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1582 42000 Incorrect parameter count in the call to native function 'curdate'
> exec
| SET @omni_temporal = ?
: string
| select tp.curdate(1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1305 42000 FUNCTION tp.curdate does not exist
> exec
| SET @omni_temporal = ?
: string
| select curdate() over ()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 21 near "over ()" 
> exec
| SET @omni_temporal = ?
: string
| select current_date
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select current_date()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select current_date(1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 21 near "1)" 
> exec
| SET @omni_temporal = ?
: string
| select current_date(null)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 24 near "null)" 
> exec
| SET @omni_temporal = ?
: string
| select current_date (1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 22 near "1)" 
> exec
| SET @omni_temporal = ?
: string
| select current_date ()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select 1 as current_date
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 24 near "current_date" 
> exec
| SET @omni_temporal = ?
: string
| select tp.current_date(1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1305 42000 FUNCTION tp.current_date does not exist
> exec
| SET @omni_temporal = ?
: string
| select `current_date`(1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1582 42000 Incorrect parameter count in the call to native function 'current_date'
> exec
| SET @omni_temporal = ?
: string
| select current_date() over ()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 26 near "over ()" 
> exec
| SET @omni_temporal = ?
: string
| select utc_date
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select utc_date()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select utc_date(1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 17 near "1)" 
> exec
| SET @omni_temporal = ?
: string
| select utc_date(null)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 20 near "null)" 
> exec
| SET @omni_temporal = ?
: string
| select utc_date (1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 18 near "1)" 
> exec
| SET @omni_temporal = ?
: string
| select 1 as utc_date
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 20 near "utc_date" 
> exec
| SET @omni_temporal = ?
: string
| select current_user
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select current_user()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select current_user(1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 21 near "1)" 
> exec
| SET @omni_temporal = ?
: string
| select current_user (1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 22 near "1)" 
> exec
| SET @omni_temporal = ?
: string
| select current_role
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select current_role()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| SELECT CURRENT_DATE, CURRENT_DATE(), CURDATE(1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 46 near "1)" 
> exec
| SET @omni_temporal = ?
: string
| SELECT CURRENT_DATE, CURRENT_DATE(), CURDATE()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select date_add
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'date_add' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select date_add from t
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select date_add()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 17 near ")" 
> exec
| SET @omni_temporal = ?
: string
| select date_add(a)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 18 near ")" 
> exec
| SET @omni_temporal = ?
: string
| select date_add(a, 10)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 21 near "10)" 
> exec
| SET @omni_temporal = ?
: string
| select date_add(a, 0.10)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 23 near "0.10)" 
> exec
| SET @omni_temporal = ?
: string
| select date_add(a, "11,11")
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 26 near ""11,11")" 
> exec
| SET @omni_temporal = ?
: string
| select date_add(a, b)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 20 near "b)" 
> exec
| SET @omni_temporal = ?
: string
| select date_add(a, null)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 23 near "null)" 
> exec
| SET @omni_temporal = ?
: string
| select date_add(a, interval 1 day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select date_add(a, INTERVAL 1 SECOND_MICROSECOND)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select date_add(a, interval b day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select date_add(a, interval interval(1,2) day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select date_add(a, interval 1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 30 near ")" 
> exec
| SET @omni_temporal = ?
: string
| select date_add(a, interval 1 foo)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 33 near "foo)" 
> exec
| SET @omni_temporal = ?
: string
| select date_add(a, interval 1 day + 1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 35 near "+ 1)" 
> exec
| SET @omni_temporal = ?
: string
| select date_add(a, interval 1 day + interval 1 day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 35 near "+ interval 1 day)" 
> exec
| SET @omni_temporal = ?
: string
| select date_add(a, 1 + interval 1 day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 20 near "1 + interval 1 day)" 
> exec
| SET @omni_temporal = ?
: string
| select date_add(a, (interval 1 day))
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 20 near "(interval 1 day))" 
> exec
| SET @omni_temporal = ?
: string
| select date_add(a, b + interval 1 day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 20 near "b + interval 1 day)" 
> exec
| SET @omni_temporal = ?
: string
| select date_add(a, interval 1 day, 1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 34 near ", 1)" 
> exec
| SET @omni_temporal = ?
: string
| select date_add(a, 1, 2)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 20 near "1, 2)" 
> exec
| SET @omni_temporal = ?
: string
| select date_add(a, interval 1 sql_tsi_day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select date_add(interval 1 day + a, interval 1 day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select date_add(a + interval 1 day, interval 1 day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select date_add(interval 1 day, interval 1 day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 31 near ", interval 1 day)" 
> exec
| SET @omni_temporal = ?
: string
| select date_add(*, interval 1 day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 17 near "*, interval 1 day)" 
> exec
| SET @omni_temporal = ?
: string
| select date_add(distinct a, interval 1 day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 24 near "distinct a, interval 1 day)" 
> exec
| SET @omni_temporal = ?
: string
| select date_add(all a, interval 1 day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 19 near "all a, interval 1 day)" 
> exec
| SET @omni_temporal = ?
: string
| select date_add(a, interval 1 day) over ()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 39 near "over ()" 
> exec
| SET @omni_temporal = ?
: string
| select date_add (a, 10)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 17 near "(a, 10)" 
> exec
| SET @omni_temporal = ?
: string
| select date_add (a, b)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 17 near "(a, b)" 
> exec
| SET @omni_temporal = ?
: string
| select date_add (a, interval 1 day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 17 near "(a, interval 1 day)" 
> exec
| SET @omni_temporal = ?
: string
| select date_add /**/ (a, 10)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 22 near "(a, 10)" 
> exec
| SET @omni_temporal = ?
: string
| select date_add.a from date_add
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1146 42S02 Table 'omni_temporal.date_add' doesn't exist
> exec
| SET @omni_temporal = ?
: string
| select t.date_add from t
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select 1 as date_add
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select 1 date_add
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select `date_add`(a, 10)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select `date_add`(a, interval 1 day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 36 near ")" 
> exec
| SET @omni_temporal = ?
: string
| select tp.date_add(a, 10)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1305 42000 FUNCTION tp.date_add does not exist
> exec
| SET @omni_temporal = ?
: string
| select tp.date_add(a, interval 1 day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 37 near ")" 
> exec
| SET @omni_temporal = ?
: string
| select date_add("2011-11-11 10:10:10.123456", 10)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 48 near "10)" 
> exec
| SET @omni_temporal = ?
: string
| select date_add("2011-11-11 10:10:10.123456", 0.10)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 50 near "0.10)" 
> exec
| SET @omni_temporal = ?
: string
| select date_add("2011-11-11 10:10:10.123456", "11,11")
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 53 near ""11,11")" 
> exec
| SET @omni_temporal = ?
: string
| select date_add("2011-11-11 10:10:10.123456", interval 10 microsecond)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select date_sub
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'date_sub' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select date_sub from t
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select date_sub()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 17 near ")" 
> exec
| SET @omni_temporal = ?
: string
| select date_sub(a)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 18 near ")" 
> exec
| SET @omni_temporal = ?
: string
| select date_sub(a, 10)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 21 near "10)" 
> exec
| SET @omni_temporal = ?
: string
| select date_sub(a, 0.10)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 23 near "0.10)" 
> exec
| SET @omni_temporal = ?
: string
| select date_sub(a, "11,11")
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 26 near ""11,11")" 
> exec
| SET @omni_temporal = ?
: string
| select date_sub(a, b)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 20 near "b)" 
> exec
| SET @omni_temporal = ?
: string
| select date_sub(a, null)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 23 near "null)" 
> exec
| SET @omni_temporal = ?
: string
| select date_sub(a, interval 1 day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select date_sub(a, interval 1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 30 near ")" 
> exec
| SET @omni_temporal = ?
: string
| select date_sub(a, interval 1 day + 1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 35 near "+ 1)" 
> exec
| SET @omni_temporal = ?
: string
| select date_sub(a, 1 + interval 1 day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 20 near "1 + interval 1 day)" 
> exec
| SET @omni_temporal = ?
: string
| select date_sub(a, (interval 1 day))
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 20 near "(interval 1 day))" 
> exec
| SET @omni_temporal = ?
: string
| select date_sub(a, interval 1 day, 1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 34 near ", 1)" 
> exec
| SET @omni_temporal = ?
: string
| select date_sub(a, interval 1 day) over ()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 39 near "over ()" 
> exec
| SET @omni_temporal = ?
: string
| select date_sub (a, 10)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 17 near "(a, 10)" 
> exec
| SET @omni_temporal = ?
: string
| select date_sub (a, interval 1 day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 17 near "(a, interval 1 day)" 
> exec
| SET @omni_temporal = ?
: string
| select 1 as date_sub
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select `date_sub`(a, 10)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select tp.date_sub(a, 10)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1305 42000 FUNCTION tp.date_sub does not exist
> exec
| SET @omni_temporal = ?
: string
| select adddate
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'adddate' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select adddate from t
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select adddate()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 16 near ")" 
> exec
| SET @omni_temporal = ?
: string
| select adddate(a)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 17 near ")" 
> exec
| SET @omni_temporal = ?
: string
| select adddate(a, 10)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select adddate(a, 0.10)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select adddate(a, "11,11")
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select adddate(a, b)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select adddate(a, interval 1 day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select adddate(a, interval 1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 29 near ")" 
> exec
| SET @omni_temporal = ?
: string
| select adddate(a, interval 1 day + 1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select adddate(a, b + interval 1 day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select adddate(a, (interval 1 day))
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 34 near "))" 
> exec
| SET @omni_temporal = ?
: string
| select adddate(a, 1, 2)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 20 near ", 2)" 
> exec
| SET @omni_temporal = ?
: string
| select adddate(a, interval 1 day) over ()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 38 near "over ()" 
> exec
| SET @omni_temporal = ?
: string
| select adddate (a, 10)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select adddate (a)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 18 near ")" 
> exec
| SET @omni_temporal = ?
: string
| select adddate (a, interval 1 day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select 1 as adddate
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select `adddate`(a, 10)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select subdate(a)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 17 near ")" 
> exec
| SET @omni_temporal = ?
: string
| select subdate(a, 10)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select subdate(a, interval 1 day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select subdate(a, 1, 2)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 20 near ", 2)" 
> exec
| SET @omni_temporal = ?
: string
| select subdate(a, interval 1 day + 1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select subdate (a, interval 1 day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select now ()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 12 near "()" 
> exec
| SET @omni_temporal = ?
: string
| select curdate ()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 16 near "()" 
> exec
| SET @omni_temporal = ?
: string
| select curtime ()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 16 near "()" 
> exec
| SET @omni_temporal = ?
: string
| select sysdate ()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select sysdate (1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select sysdate (b)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'b' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select sysdate (1,2)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1582 42000 Incorrect parameter count in the call to native function 'sysdate'
> exec
| SET @omni_temporal = ?
: string
| select sysdate (1) over ()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 23 near "over ()" 
> exec
| SET @omni_temporal = ?
: string
| select now (1) over ()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 12 near "(1) over ()" 
> exec
| SET @omni_temporal = ?
: string
| select curtime (1) over ()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 16 near "(1) over ()" 
> exec
| SET @omni_temporal = ?
: string
| select date_add ()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 17 near "()" 
> exec
| SET @omni_temporal = ?
: string
| select adddate ()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 17 near ")" 
> exec
| SET @omni_temporal = ?
: string
| select subdate ()
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 17 near ")" 
> exec
| SET @omni_temporal = ?
: string
| select adddate (a, b, c)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 21 near ", c)" 
> exec
| SET @omni_temporal = ?
: string
| select curdate/**/(1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 19 near "(1)" 
> exec
| SET @omni_temporal = ?
: string
| select now/**/(1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 15 near "(1)" 
> exec
| SET @omni_temporal = ?
: string
| select CURDATE (1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 16 near "(1)" 
> exec
| SET @omni_temporal = ?
: string
| select cUrTiMe(-1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 16 near "-1)" 
> exec
| SET @omni_temporal = ?
: string
| select t.current_time(-1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1305 42000 FUNCTION t.current_time does not exist
> exec
| SET @omni_temporal = ?
: string
| select t.now(-1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1305 42000 FUNCTION t.now does not exist
> exec
| SET @omni_temporal = ?
: string
| select t.curdate(1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1305 42000 FUNCTION t.curdate does not exist
> exec
| SET @omni_temporal = ?
: string
| select date_add(a, interval(1,2))
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 33 near ")" 
> exec
| SET @omni_temporal = ?
: string
| select date_sub(a, interval(1,2))
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 33 near ")" 
> exec
| SET @omni_temporal = ?
: string
| select adddate(a, interval(1,2))
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select subdate(a, interval(1,2))
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select date_add((interval 1 day), interval 1 day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 32 near "), interval 1 day)" 
> exec
| SET @omni_temporal = ?
: string
| select adddate(interval 1 day, 1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 30 near ", 1)" 
> exec
| SET @omni_temporal = ?
: string
| select adddate((interval 1 day), 1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 31 near "), 1)" 
> exec
| SET @omni_temporal = ?
: string
| select adddate(interval 1 day + a, interval 1 day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET @omni_temporal = ?
: string
| select adddate(*, 1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 16 near "*, 1)" 
> exec
| SET @omni_temporal = ?
: string
| select adddate(a, *)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 19 near "*)" 
> exec
| SET @omni_temporal = ?
: string
| select adddate(distinct a, 1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 23 near "distinct a, 1)" 
> exec
| SET @omni_temporal = ?
: string
| select adddate(a, distinct 1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 26 near "distinct 1)" 
> exec
| SET @omni_temporal = ?
: string
| select date_add(a, distinct interval 1 day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 27 near "distinct interval 1 day)" 
> exec
| SET @omni_temporal = ?
: string
| select date_add(a, interval 1 day,)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 34 near ",)" 
> exec
| SET @omni_temporal = ?
: string
| select abs(interval 1 day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 26 near ")" 
> exec
| SET @omni_temporal = ?
: string
| select concat(interval 1 day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 29 near ")" 
> exec
| SET @omni_temporal = ?
: string
| select coalesce(interval 1 day, 1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 31 near ", 1)" 
> exec
| SET @omni_temporal = ?
: string
| select interval 1 day + 1
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select 1 + interval 1 day
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select abs((interval 1 day))
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 27 near "))" 
> exec
| SET @omni_temporal = ?
: string
| select abs(interval 1 day + 1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select current_time(0xffffffffffffffffffff)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 42 near "0xffffffffffffffffffff)" 
> exec
| SET @omni_temporal = ?
: string
| select current_time(0X1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 23 near "0X1)"hex literal: invalid hexadecimal format 0X1 
> exec
| SET @omni_temporal = ?
: string
| select current_time(0b1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 23 near "0b1)" 
> exec
| SET @omni_temporal = ?
: string
| select now(0b1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select curtime(0xa)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 18 near "0xa)" 
> exec
| SET @omni_temporal = ?
: string
| select utc_time(0x1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 19 near "0x1)" 
> exec
| SET @omni_temporal = ?
: string
| select current_time(0002147483647)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1426 42000 Too-big precision 2147483647 specified for 'current_time'. Maximum is 6.
> exec
| SET @omni_temporal = ?
: string
| select current_time(2147483647)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1426 42000 Too-big precision 2147483647 specified for 'current_time'. Maximum is 6.
> exec
| SET @omni_temporal = ?
: string
| select current_time(18446744073709551615)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use
> exec
| SET @omni_temporal = ?
: string
| select current_time(18446744073709551616)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 40 near "18446744073709551616)" 
> exec
| SET @omni_temporal = ?
: string
| select now(2147483647)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1426 42000 Too-big precision 2147483647 specified for 'now'. Maximum is 6.
> exec
| SET @omni_temporal = ?
: string
| select now(2147483648)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use
> exec
| SET @omni_temporal = ?
: string
| select now(18446744073709551615)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use
> exec
| SET @omni_temporal = ?
: string
| select now(18446744073709551616)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use
> exec
| SET @omni_temporal = ?
: string
| select now(99999999999999999999999)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use
> exec
| SET @omni_temporal = ?
: string
| select curtime(2147483647)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1426 42000 Too-big precision 2147483647 specified for 'curtime'. Maximum is 6.
> exec
| SET @omni_temporal = ?
: string
| select curtime(18446744073709551615)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use
> exec
| SET @omni_temporal = ?
: string
| select sysdate(2147483647)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1426 42000 Too-big precision 2147483647 specified for 'sysdate'. Maximum is 6.
> exec
| SET @omni_temporal = ?
: string
| select sysdate(2147483648)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use
> exec
| SET @omni_temporal = ?
: string
| select sysdate(18446744073709551615)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use
> exec
| SET @omni_temporal = ?
: string
| select sysdate(18446744073709551616)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 35 near "18446744073709551616)" 
> exec
| SET @omni_temporal = ?
: string
| select utc_time(2147483648)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use
> exec
| SET @omni_temporal = ?
: string
| select utc_timestamp(2147483648)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use
> exec
| SET @omni_temporal = ?
: string
| select localtime(2147483648)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use
> exec
| SET @omni_temporal = ?
: string
| select localtimestamp(2147483648)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use
> exec
| SET @omni_temporal = ?
: string
| create table t1 (a timestamp default current_timestamp(6) on update current_timestamp(6))
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| create table t2 (a timestamp default current_timestamp(-1))
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 56 near "-1))" 
> exec
| SET @omni_temporal = ?
: string
| create table t3 (a timestamp default now(6) on update now(6))
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| create table t5 (a timestamp on update current_timestamp(1.0))
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 60 near "1.0))" 
> exec
| SET @omni_temporal = ?
: string
| create table t6 (a date default curdate(1))
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 41 near "1))" 
> exec
| SET @omni_temporal = ?
: string
| create table t7 (a date default (curdate(1)))
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 42 near "1)))" 
> exec
| SET @omni_temporal = ?
: string
| create table t8 (a datetime default (date_add(now(), 1)))
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 45 near "date_add(now(), 1)))" 
> exec
| SET @omni_temporal = ?
: string
| create table t10 (a datetime default current_timestamp(7))
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| create table t11 (a datetime default now(7))
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp(00000000000000000000000018446744073709551615)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp(00000000000000000000000018446744073709551616)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 69 near "00000000000000000000000018446744073709551616)" 
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp(0000000000002147483647)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1426 42000 Too-big precision 2147483647 specified for 'current_timestamp'. Maximum is 6.
> exec
| SET @omni_temporal = ?
: string
| select current_timestamp(0000000000002147483648)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use
> exec
| SET @omni_temporal = ?
: string
| select now(0000000000002147483648)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use
> exec
| SET @omni_temporal = ?
: string
| select current_time(x'0xaa')
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 23 near "x'0xaa')" 
> exec
| SET @omni_temporal = ?
: string
| select current_date (
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 21 near "" 
> exec
| SET @omni_temporal = ?
: string
| select now(1)(2)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 14 near "(2)" 
> exec
| SET @omni_temporal = ?
: string
| select curtime(1)(2)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 18 near "(2)" 
> exec
| SET @omni_temporal = ?
: string
| select current_time(1)(2)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 23 near "(2)" 
> exec
| SET @omni_temporal = ?
: string
| select now(1) x
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select curtime(1) as x from t
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select now() + interval 1 day
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select curdate() - interval 1 day
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select date_add(a, interval 1 day) x from t
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select date_sub(a, interval 1 day) as x from t
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select date_add(a, interval 1 day) + 1 from t
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select date_add(a,
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 18 near "" 
> exec
| SET @omni_temporal = ?
: string
| select date_add(a
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 17 near "" 
> exec
| SET @omni_temporal = ?
: string
| select date_add(a, interval 1 day
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 33 near "" 
> exec
| SET @omni_temporal = ?
: string
| select date_add(, interval 1 day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 17 near ", interval 1 day)" 
> exec
| SET @omni_temporal = ?
: string
| select adddate(a,)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 18 near ")" 
> exec
| SET @omni_temporal = ?
: string
| select adddate(,a)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 16 near ",a)" 
> exec
| SET @omni_temporal = ?
: string
| select adddate(a, 1
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 19 near "" 
> exec
| SET @omni_temporal = ?
: string
| select interval(interval 1 day, 1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 31 near ", 1)" 
> exec
| SET @omni_temporal = ?
: string
| select if(interval 1 day, 1, 2)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 25 near ", 1, 2)" 
> exec
| SET @omni_temporal = ?
: string
| select left(interval 1 day, 1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 27 near ", 1)" 
> exec
| SET @omni_temporal = ?
: string
| select count(interval 1 day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 28 near ")" 
> exec
| SET @omni_temporal = ?
: string
| select sum(interval 1 day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 26 near ")" 
> exec
| SET @omni_temporal = ?
: string
| select max(a + interval 1 day) from t
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select date_format(now(), '%Y')
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select date_format(a + interval 1 day, '%Y') from t
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select extract(day from a + interval 1 day) from t
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select timestampadd(day, 1, now(6))
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select now(now())
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use
> exec
| SET @omni_temporal = ?
: string
| select curtime(now())
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 18 near "now())" 
> exec
| SET @omni_temporal = ?
: string
| select date_add(now(6), interval 1 day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select date_add(date_add(a, interval 1 day), interval 1 day) from t
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select date_add(a, interval date_add(b, interval 1 day) day) from t
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select adddate(adddate(a, 1), 1) from t
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select date_sub(curdate(), interval 1 month)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select sysdate(3) + interval 1 day
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select utc_timestamp(3) + 1
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| update t set a = now(6)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| update t set a = now(-1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use
> exec
| SET @omni_temporal = ?
: string
| update t set a = date_add(a, interval 1 day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| update t set a = date_add(a, 1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 30 near "1)" 
> exec
| SET @omni_temporal = ?
: string
| delete from t where a < date_sub(now(), interval 1 day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| delete from t where a < date_sub(now(), 1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 41 near "1)" 
> exec
| SET @omni_temporal = ?
: string
| insert into t (a) values (current_timestamp(3))
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| insert into t (a) values (current_timestamp(-1))
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 45 near "-1))" 
> exec
| SET @omni_temporal = ?
: string
| select * from t where a > curdate() order by now(6)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_temporal_stmt
< ok 0 0
> exec
| SET @omni_temporal = ?
: string
| select * from t where a > curdate(1)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 35 near "1)" 
> exec
| SET @omni_temporal = ?
: string
| select bit_and(interval 1 day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 30 near ")" 
> exec
| SET @omni_temporal = ?
: string
| select json_arrayagg(interval 1 day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 36 near ")" 
> exec
| SET @omni_temporal = ?
: string
| select char(interval 1 day)
< ok 0 0
> exec
| PREPARE omni_temporal_stmt FROM @omni_temporal
< error
| 1064 42000 …line 1 column 27 near ")" 
> exec
| DROP DATABASE IF EXISTS omni_temporal
< ok 0 0
