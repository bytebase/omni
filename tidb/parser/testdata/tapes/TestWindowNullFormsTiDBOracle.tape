-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestWindowNullFormsTiDBOracle default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| SELECT FIRST_VALUE(val) RESPECT NULLS OVER w FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT FIRST_VALUE(val) RESPECT NULLS OVER w FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT FIRST_VALUE(val) IGNORE NULLS OVER w FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT FIRST_VALUE(val) IGNORE NULLS OVER w FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT LAST_VALUE(val) IGNORE NULLS OVER w FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT LAST_VALUE(val) IGNORE NULLS OVER w FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT LEAD(val) IGNORE NULLS OVER w FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT LEAD(val) IGNORE NULLS OVER w FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT LEAD(val, 2) RESPECT NULLS OVER w FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT LEAD(val, 2) RESPECT NULLS OVER w FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT LAG(val, 2, 0) IGNORE NULLS OVER w FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT LAG(val, 2, 0) IGNORE NULLS OVER w FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT NTH_VALUE(val, 233) FROM FIRST OVER w FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT NTH_VALUE(val, 233) FROM FIRST OVER w FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT NTH_VALUE(val, 233) FROM LAST OVER w FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT NTH_VALUE(val, 233) FROM LAST OVER w FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT NTH_VALUE(val, 233) FROM LAST IGNORE NULLS OVER w FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT NTH_VALUE(val, 233) FROM LAST IGNORE NULLS OVER w FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT NTH_VALUE(val, 233) FROM FIRST RESPECT NULLS OVER w FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT NTH_VALUE(val, 233) FROM FIRST RESPECT NULLS OVER w FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT NTH_VALUE(val, 1) IGNORE NULLS OVER w FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT NTH_VALUE(val, 1) IGNORE NULLS OVER w FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT FIRST_VALUE(val) IGNORE NULLS OVER () FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT FIRST_VALUE(val) IGNORE NULLS OVER () FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT FIRST_VALUE(val) IGNORE NULLS FROM t
< error
| 1064 42000 …line 1 column 41 near "FROM t" 
> exec ^
: string
| SELECT NTH_VALUE(val, 1) FROM LAST FROM t
< error
| 1064 42000 …line 1 column 39 near "FROM t" 
> exec ^
: string
| SELECT NTH_VALUE(val, 1) IGNORE NULLS FROM LAST OVER w FROM t
< error
| 1064 42000 …line 1 column 42 near "FROM LAST OVER w FROM t" 
> exec ^
: string
| SELECT FIRST_VALUE(val) FROM LAST OVER w FROM t
< error
| 1064 42000 …line 1 column 28 near "FROM LAST OVER w FROM t" 
> exec ^
: string
| SELECT LEAD(val) FROM FIRST OVER w FROM t
< error
| 1064 42000 …line 1 column 21 near "FROM FIRST OVER w FROM t" 
> exec ^
: string
| SELECT COUNT(val) IGNORE NULLS OVER w FROM t
< error
| 1064 42000 …line 1 column 24 near "IGNORE NULLS OVER w FROM t" 
> exec ^
: string
| SELECT SUM(val) RESPECT NULLS OVER w FROM t
< error
| 1064 42000 …line 1 column 29 near "NULLS OVER w FROM t" 
> exec ^
: string
| SELECT ROW_NUMBER() IGNORE NULLS OVER w FROM t
< error
| 1064 42000 …line 1 column 26 near "IGNORE NULLS OVER w FROM t" 
> exec ^
: string
| SELECT FIRST_VALUE(val) IGNORE OVER w FROM t
< error
| 1064 42000 …line 1 column 35 near "OVER w FROM t" 
> exec ^
: string
| SELECT FIRST_VALUE(val) RESPECT OVER w FROM t
< error
| 1064 42000 …line 1 column 36 near "OVER w FROM t" 
> exec ^
: string
| SELECT FIRST_VALUE(val) IGNORE NULL OVER w FROM t
< error
| 1064 42000 …line 1 column 35 near "NULL OVER w FROM t" 
> exec ^
: string
| SELECT `first_value`(val) IGNORE NULLS OVER w FROM t
< error
| 1064 42000 …line 1 column 32 near "IGNORE NULLS OVER w FROM t" 
