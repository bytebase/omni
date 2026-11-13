-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestStatsLockFormsOracle default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| LOCK STATS t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| LOCK STATS t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| UNLOCK STATS t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| UNLOCK STATS t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| lock stats test.t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| lock stats test.t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| unlock stats test.t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| unlock stats test.t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| LOCK STATS t, t2
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| LOCK STATS t, t2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| UNLOCK STATS t, t2
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| UNLOCK STATS t, t2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| LOCK STATS test.t, test.t2
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| LOCK STATS test.t, test.t2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| LOCK STATS `test`.`t`, t2
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| LOCK STATS `test`.`t`, t2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| LOCK STATS t PARTITION (p0)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| LOCK STATS t PARTITION (p0)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| LOCK STATS t PARTITION (p0, p1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| LOCK STATS t PARTITION (p0, p1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| LOCK STATS t PARTITION p0
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| LOCK STATS t PARTITION p0
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| LOCK STATS t PARTITION p0, p1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| LOCK STATS t PARTITION p0, p1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| UNLOCK STATS t PARTITION (p0, p1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| UNLOCK STATS t PARTITION (p0, p1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| UNLOCK STATS t PARTITION p0
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| UNLOCK STATS t PARTITION p0
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| LOCK STATS test.t PARTITION p0
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| LOCK STATS test.t PARTITION p0
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| LOCK STATS t PARTITION `p0`
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| LOCK STATS t PARTITION `p0`
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| LOCK STATS t, t2 PARTITION (p0)
< error
| 1064 42000 …line 1 column 26 near "PARTITION (p0)" 
> exec ^
: string
| LOCK STATS t, t2 PARTITION p0
< error
| 1064 42000 …line 1 column 26 near "PARTITION p0" 
> exec ^
: string
| LOCK STATS t PARTITION (p0), t2
< error
| 1064 42000 …line 1 column 28 near ", t2" 
> exec ^
: string
| LOCK STATS t PARTITION (p0), t2 PARTITION (p1)
< error
| 1064 42000 …line 1 column 28 near ", t2 PARTITION (p1)" 
> exec ^
: string
| LOCK STATS
< error
| 1064 42000 …line 1 column 10 near "" 
> exec ^
: string
| UNLOCK STATS
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| LOCK STATS PARTITION (p0)
< error
| 1064 42000 …line 1 column 20 near "PARTITION (p0)" 
> exec ^
: string
| LOCK STATS t,
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| LOCK STATS t PARTITION
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| LOCK STATS t PARTITION ()
< error
| 1064 42000 …line 1 column 25 near ")" 
> exec ^
: string
| LOCK STATS t PARTITION (p0,)
< error
| 1064 42000 …line 1 column 28 near ")" 
> exec ^
: string
| LOCK STATS t PARTITION ((p0))
< error
| 1064 42000 …line 1 column 25 near "(p0))" 
> exec ^
: string
| LOCK STATS t PARTITION (p0) (p1)
< error
| 1064 42000 …line 1 column 29 near "(p1)" 
> exec ^
: string
| LOCK STATS t PARTITION (p0) PARTITION (p1)
< error
| 1064 42000 …line 1 column 37 near "PARTITION (p1)" 
> exec ^
: string
| LOCK STATS t PARTITION p0 p1
< error
| 1064 42000 …line 1 column 28 near "p1" 
> exec ^
: string
| LOCK STATS = t
< error
| 1064 42000 …line 1 column 12 near "= t" 
> exec ^
: string
| UNLOCK STATS = t
< error
| 1064 42000 …line 1 column 14 near "= t" 
> exec ^
: string
| LOCK STATS 1
< error
| 1064 42000 …line 1 column 12 near "1" 
> exec ^
: string
| LOCK STATS 't'
< error
| 1064 42000 …line 1 column 14 near "'t'" 
> exec ^
: string
| LOCK STATS select
< error
| 1064 42000 …line 1 column 17 near "select" 
> exec ^
: string
| LOCK STATS t PARTITION partition
< error
| 1064 42000 …line 1 column 32 near "partition" 
> exec ^
: string
| LOCK STATS t READ
< error
| 1064 42000 …line 1 column 17 near "READ" 
> exec ^
: string
| LOCK STAT t
< error
| 1064 42000 …line 1 column 9 near "STAT t" 
> exec ^
: string
| LOCK STATISTICS t
< error
| 1064 42000 …line 1 column 15 near "STATISTICS t" 
> exec ^
: string
| UNLOCK STAT t
< error
| 1064 42000 …line 1 column 11 near "STAT t" 
> exec ^
: string
| LOCK STATS `t`
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| LOCK STATS `t`
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| LOCK `STATS` t
< error
| 1064 42000 …line 1 column 12 near "`STATS` t" 
> exec ^
: string
| UNLOCK `STATS` t
< error
| 1064 42000 …line 1 column 14 near "`STATS` t" 
