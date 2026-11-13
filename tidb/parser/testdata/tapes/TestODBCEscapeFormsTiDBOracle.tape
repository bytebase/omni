-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestODBCEscapeFormsTiDBOracle default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| select {d '2020-01-01'}
< ok 0 0
> exec ^
: string
| select {D '2020-01-01'}
< ok 0 0
> exec ^
: string
| select {TS '2020-01-01 00:00:00'}
< ok 0 0
> exec ^
: string
| select {d 123}
< error
| 1292 22007 Incorrect date value: ''
> exec ^
: string
| select {d 123}
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select {d now()}
< error
| 1292 22007 Incorrect date value: ''
> exec ^
: string
| select {d now()}
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select {ts 1 + 1}
< error
| 1292 22007 Incorrect datetime value: '2'
> exec ^
: string
| select {ts 1 + 1}
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select {`d` '2020-01-01'}
< ok 0 0
> exec ^
: string
| select {abc 1}
< ok 0 0
> exec ^
: string
| select {select 1}
< error
| 1064 42000 …line 1 column 14 near "select 1}" 
> exec ^
: string
| select {d}
< error
| 1064 42000 …line 1 column 10 near "}" 
> exec ^
: string
| select {}
< error
| 1064 42000 …line 1 column 9 near "}" 
> exec ^
: string
| select {d '2020-01-01'} + 1
< ok 0 0
> exec ^
: string
| select { d '2020-01-01' }
< ok 0 0
> exec ^
: string
| select {"d" '2020-01-01'}
< error
| 1064 42000 …line 1 column 11 near ""d" '2020-01-01'}" 
> exec ^
: string
| select {d '2020-01-01'} as x
< ok 0 0
> exec ^
: string
| select {date '2020-01-01'}
< ok 0 0
> exec ^
: string
| select {timestamp '2020-01-01'}
< ok 0 0
> exec ^
: string
| select {from 1}
< error
| 1064 42000 …line 1 column 12 near "from 1}" 
> exec ^
: string
| select {d (select 1)}
< error
| 1292 22007 Incorrect date value: ''
> exec ^
: string
| select {d (select 1)}
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select {d 1, 2}
< error
| 1064 42000 …line 1 column 12 near ", 2}" 
> exec ^
: string
| insert into t values ({d '2020-01-01'})
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values ({d '2020-01-01'})
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select {ts {abc {d '1989-09-10'}}}
< error
| 1292 22007 Incorrect datetime value: '1989-09-10'
> exec ^
: string
| select {ts {abc {d '1989-09-10'}}}
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select {d '1' } from t where a in ({d '1'})
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select {d '1' } from t where a in ({d '1'})
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select {oj 1}
< ok 0 0
> exec ^
: string
| select {d '1'
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| select {abc *} from t
< error
| 1064 42000 …line 1 column 13 near "*} from t" 
> exec ^
: string
| select {abc t.*} from t
< error
| 1064 42000 …line 1 column 15 near "*} from t" 
> exec ^
: string
| select {d *}
< error
| 1064 42000 …line 1 column 11 near "*}" 
> exec ^
: string
| select {d default}
< error
| 1064 42000 …line 1 column 18 near "}" 
> exec ^
: string
| select {d default(a)} from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select {d default(a)} from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values ({d *})
< error
| 1064 42000 …line 1 column 26 near "*})" 
> exec ^
: string
| insert into t values ({abc default})
< error
| 1064 42000 …line 1 column 35 near "})" 
> exec ^
: string
| select * from {oj t1 left join t2 on t1.a = t2.a}
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| select * from {oj t1 left join t2 on t1.a = t2.a}
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from {oj t1}
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| select * from {oj t1}
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from {abc t1 left join t2 on t1.a = t2.a}
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| select * from {abc t1 left join t2 on t1.a = t2.a}
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from {oj t1 left join t2 on t1.a = t2.a}, t3
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| select * from {oj t1 left join t2 on t1.a = t2.a}, t3
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t3, {oj t1 left join t2 on t1.a = t2.a}
< error
| 1146 42S02 Table 'test.t3' doesn't exist
> exec ^
: string
| select * from t3, {oj t1 left join t2 on t1.a = t2.a}
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from {oj t1} join t2 on t1.a = t2.a
< error
| 1064 42000 …line 1 column 26 near "join t2 on t1.a = t2.a" 
> exec ^
: string
| select * from ({oj t1})
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| select * from ({oj t1})
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from {oj {oj t1}}
< error
| 1064 42000 …line 1 column 19 near "{oj t1}}" 
> exec ^
: string
| update {oj t1} set a = 1
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| update {oj t1} set a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| delete from {oj t1}
< error
| 1064 42000 …line 1 column 13 near "{oj t1}" 
> exec ^
: string
| delete t1 from {oj t1 left join t2 on t1.a = t2.a}
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| delete t1 from {oj t1 left join t2 on t1.a = t2.a}
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| delete from t1 using {oj t1 left join t2 on t1.a = t2.a}
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| delete from t1 using {oj t1 left join t2 on t1.a = t2.a}
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
