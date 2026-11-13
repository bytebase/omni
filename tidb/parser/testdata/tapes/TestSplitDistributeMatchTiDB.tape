-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestSplitDistributeMatchTiDB default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| SPLIT TABLE t1 BY ('a'),('b'),('c')
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SPLIT TABLE t1 BY ('a'),('b'),('c')
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SPLIT TABLE t1 BY (1)
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SPLIT TABLE t1 BY (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SPLIT TABLE t1 INDEX idx1 BY ('abc',123), ('xyz'), ('yz', 1000)
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SPLIT TABLE t1 INDEX idx1 BY ('abc',123), ('xyz'), ('yz', 1000)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SPLIT TABLE t1 BETWEEN ('a') AND ('z') REGIONS 10
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SPLIT TABLE t1 BETWEEN ('a') AND ('z') REGIONS 10
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SPLIT TABLE t1 INDEX idx1 BETWEEN ('a',1) AND ('z',2) REGIONS 10
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SPLIT TABLE t1 INDEX idx1 BETWEEN ('a',1) AND ('z',2) REGIONS 10
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SPLIT REGION FOR TABLE t1 BY ('a')
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SPLIT REGION FOR TABLE t1 BY ('a')
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SPLIT PARTITION TABLE t1 BY ('a')
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SPLIT PARTITION TABLE t1 BY ('a')
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SPLIT REGION FOR PARTITION TABLE t1 BY ('a')
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SPLIT REGION FOR PARTITION TABLE t1 BY ('a')
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SPLIT REGION FOR PARTITION TABLE t1 PARTITION (p0,p1) INDEX idx1 BETWEEN ('a') AND ('z') REGIONS 10
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SPLIT REGION FOR PARTITION TABLE t1 PARTITION (p0,p1) INDEX idx1 BETWEEN ('a') AND ('z') REGIONS 10
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SPLIT TABLE t1 BY ()
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SPLIT TABLE t1 BY ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SPLIT TABLE t1 BY ('a'), ()
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SPLIT TABLE t1 BY ('a'), ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SPLIT TABLE t1 BETWEEN () AND () REGIONS 10
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SPLIT TABLE t1 BETWEEN () AND () REGIONS 10
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SPLIT TABLE t1 INDEX idx1 BETWEEN () AND ('z') REGIONS 10
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SPLIT TABLE t1 INDEX idx1 BETWEEN () AND ('z') REGIONS 10
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SPLIT TABLE t1 BY ('a'),(1,2,3)
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SPLIT TABLE t1 BY ('a'),(1,2,3)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SPLIT TABLE t1 BETWEEN ('a','b') AND ('z') REGIONS 10
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SPLIT TABLE t1 BETWEEN ('a','b') AND ('z') REGIONS 10
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SPLIT TABLE t1 BETWEEN ('a') AND ('z') REGIONS 0
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SPLIT TABLE t1 BETWEEN ('a') AND ('z') REGIONS 0
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SPLIT TABLE t1 BY (1+1)
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SPLIT TABLE t1 BY (1+1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SPLIT TABLE t1 BY (now())
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SPLIT TABLE t1 BY (now())
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SPLIT TABLE t1 BY (-1)
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SPLIT TABLE t1 BY (-1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SPLIT TABLE t1 BY (null)
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SPLIT TABLE t1 BY (null)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SPLIT TABLE t1 BY (b)
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SPLIT TABLE t1 BY (b)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SPLIT TABLE probe.t1 BY ('a')
< error
| 1146 42S02 Table 'probe.t1' doesn't exist
> exec ^
: string
| SPLIT TABLE probe.t1 BY ('a')
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SPLIT TABLE `t1` INDEX `idx1` BY ('a')
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SPLIT TABLE `t1` INDEX `idx1` BY ('a')
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SPLIT TABLE t1 INDEX split BY ('a')
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SPLIT TABLE t1 INDEX split BY ('a')
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SPLIT TABLE t1 INDEX region BY ('a')
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SPLIT TABLE t1 INDEX region BY ('a')
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SPLIT TABLE t1 INDEX rule BY ('a')
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SPLIT TABLE t1 INDEX rule BY ('a')
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SPLIT TABLE t1 INDEX timeout BY ('a')
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SPLIT TABLE t1 INDEX timeout BY ('a')
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SPLIT TABLE t1 INDEX distribute BY ('a')
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SPLIT TABLE t1 INDEX distribute BY ('a')
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SPLIT TABLE split BY ('a')
< error
| 1146 42S02 Table 'test.split' doesn't exist
> exec ^
: string
| SPLIT TABLE split BY ('a')
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SPLIT TABLE t1 PARTITION (split, rule) BY ('a')
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SPLIT TABLE t1 PARTITION (split, rule) BY ('a')
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SPLIT TABLE t1 INDEX `primary` BY (1)
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SPLIT TABLE t1 INDEX `primary` BY (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DISTRIBUTE TABLE t1 RULE = 'leader-scatter' ENGINE = 'tikv'
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| DISTRIBUTE TABLE t1 RULE = 'leader-scatter' ENGINE = 'tikv'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DISTRIBUTE TABLE t1 RULE = "leader-scatter" ENGINE = "tikv"
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| DISTRIBUTE TABLE t1 RULE = "leader-scatter" ENGINE = "tikv"
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DISTRIBUTE TABLE t1 RULE 'leader-scatter' ENGINE 'tikv'
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| DISTRIBUTE TABLE t1 RULE 'leader-scatter' ENGINE 'tikv'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DISTRIBUTE TABLE t1 RULE = 'x' ENGINE = 'tikv' TIMEOUT = '30m'
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| DISTRIBUTE TABLE t1 RULE = 'x' ENGINE = 'tikv' TIMEOUT = '30m'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DISTRIBUTE TABLE t1 RULE 'x' ENGINE 'tikv' TIMEOUT '30m'
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| DISTRIBUTE TABLE t1 RULE 'x' ENGINE 'tikv' TIMEOUT '30m'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DISTRIBUTE TABLE t1 RULE = 'x' ENGINE = 'tikv' TIMEOUT = ''
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| DISTRIBUTE TABLE t1 RULE = 'x' ENGINE = 'tikv' TIMEOUT = ''
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DISTRIBUTE TABLE t1 PARTITION(p0,p1) RULE = 'learner-scatter' ENGINE = 'tikv'
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| DISTRIBUTE TABLE t1 PARTITION(p0,p1) RULE = 'learner-scatter' ENGINE = 'tikv'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DISTRIBUTE TABLE probe.t1 RULE = 'x' ENGINE = 'tikv'
< error
| 1146 42S02 Table 'probe.t1' doesn't exist
> exec ^
: string
| DISTRIBUTE TABLE probe.t1 RULE = 'x' ENGINE = 'tikv'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DISTRIBUTE TABLE split RULE = 'x' ENGINE = 'tikv'
< error
| 1146 42S02 Table 'test.split' doesn't exist
> exec ^
: string
| DISTRIBUTE TABLE split RULE = 'x' ENGINE = 'tikv'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT 1 AS split
< ok 0 0
> exec ^
: string
| SELECT 1 AS distribute
< ok 0 0
> exec ^
: string
| SELECT 1 AS region
< ok 0 0
> exec ^
: string
| SELECT 1 AS rule
< ok 0 0
> exec ^
: string
| SELECT 1 AS timeout
< ok 0 0
> exec ^
: string
| SELECT split, distribute, region, rule, timeout FROM t1
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SELECT split, distribute, region, rule, timeout FROM t1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t1 AS rule
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SELECT a FROM t1 AS rule
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT rule.a FROM t1 rule
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SELECT rule.a FROM t1 rule
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t1 split WHERE split.a = 'x'
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| DELETE FROM t1 split WHERE split.a = 'x'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t1 distribute WHERE distribute.a = 'x'
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| DELETE FROM t1 distribute WHERE distribute.a = 'x'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t1 region WHERE region.a = 'x'
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| DELETE FROM t1 region WHERE region.a = 'x'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t1 rule WHERE rule.a = 'x'
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| DELETE FROM t1 rule WHERE rule.a = 'x'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t1 timeout WHERE timeout.a = 'x'
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| DELETE FROM t1 timeout WHERE timeout.a = 'x'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE split (split int, distribute int, region int, rule int, timeout int)
< ok 0 0
> exec ^
: string
| CREATE INDEX split ON t1 (a)
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| CREATE INDEX split ON t1 (a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| UPDATE t1 SET rule = 1 WHERE timeout = 2
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| UPDATE t1 SET rule = 1 WHERE timeout = 2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| INSERT INTO t1 (rule, timeout) VALUES (1,2)
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| INSERT INTO t1 (rule, timeout) VALUES (1,2)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SAVEPOINT split
< ok 0 0
> exec ^
: string
| ANALYZE TABLE split
< error
| 1146 42S02 Table 'test.split' doesn't exist
> exec ^
: string
| ANALYZE TABLE split
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TRUNCATE TABLE region
< error
| 1146 42S02 Table 'test.region' doesn't exist
> exec ^
: string
| TRUNCATE TABLE region
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SPLIT TABLE t1
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| SPLIT TABLE
< error
| 1064 42000 …line 1 column 11 near "" 
> exec ^
: string
| SPLIT TABLE t1 INDEX idx1
< error
| 1064 42000 …line 1 column 25 near "" 
> exec ^
: string
| SPLIT TABLE t1 REGIONS 10
< error
| 1064 42000 …line 1 column 22 near "REGIONS 10" 
> exec ^
: string
| SPLIT TABLE t1 INDEX idx1 REGIONS 10
< error
| 1064 42000 …line 1 column 33 near "REGIONS 10" 
> exec ^
: string
| SPLIT TABLE t1 BY ('a') BETWEEN ('a') AND ('z') REGIONS 10
< error
| 1064 42000 …line 1 column 31 near "BETWEEN ('a') AND ('z') REGIONS 10" 
> exec ^
: string
| SPLIT TABLE t1 BETWEEN ('a') AND ('z') REGIONS 10 BY ('a')
< error
| 1064 42000 …line 1 column 52 near "BY ('a')" 
> exec ^
: string
| SPLIT TABLE t1 BY
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| SPLIT TABLE t1 INDEX idx1 BY
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| SPLIT TABLE t1 BY ('a'),
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| SPLIT TABLE t1 BY row('a')
< error
| 1064 42000 …line 1 column 21 near "row('a')" 
> exec ^
: string
| SPLIT TABLE t1 BETWEEN ('a') AND ('z')
< error
| 1064 42000 …line 1 column 38 near "" 
> exec ^
: string
| SPLIT TABLE t1 BETWEEN ('a') REGIONS 10
< error
| 1064 42000 …line 1 column 36 near "REGIONS 10" 
> exec ^
: string
| SPLIT TABLE t1 BETWEEN ('a') AND ('z') REGIONS
< error
| 1064 42000 …line 1 column 46 near "" 
> exec ^
: string
| SPLIT TABLE t1 BETWEEN (1) AND (2) REGIONS -1
< error
| 1064 42000 …line 1 column 44 near "-1" 
> exec ^
: string
| SPLIT TABLE t1 BETWEEN (1) AND (2) REGIONS +10
< error
| 1064 42000 …line 1 column 44 near "+10" 
> exec ^
: string
| SPLIT TABLE t1 BETWEEN (1) AND (2) REGIONS 10.5
< error
| 1064 42000 …line 1 column 47 near "10.5" 
> exec ^
: string
| SPLIT TABLE t1 BETWEEN (1) AND (2) REGIONS 1e3
< error
| 1064 42000 …line 1 column 46 near "1e3" 
> exec ^
: string
| SPLIT TABLE t1 BETWEEN (1) AND (2) REGIONS 0x10
< error
| 1064 42000 …line 1 column 47 near "0x10" 
> exec ^
: string
| SPLIT TABLE t1 BETWEEN (1) AND (2) REGIONS 'a'
< error
| 1064 42000 …line 1 column 46 near "'a'" 
> exec ^
: string
| SPLIT TABLE t1 BETWEEN (1) AND (2) REGIONS ?
< error
| 1064 42000 …line 1 column 44 near "?" 
> exec ^
: string
| SPLIT TABLE t1 BETWEEN (1) AND (2) REGIONS 99999999999999999999
< error
| 1064 42000 …line 1 column 63 near "99999999999999999999" 
> exec ^
: string
| SPLIT TABLE t1, t2 BY ('a')
< error
| 1064 42000 …line 1 column 15 near ", t2 BY ('a')" 
> exec ^
: string
| SPLIT TABLE t1 INDEX idx1, idx2 BY ('a')
< error
| 1064 42000 …line 1 column 26 near ", idx2 BY ('a')" 
> exec ^
: string
| SPLIT TABLE t1 INDEX BY (1)
< error
| 1064 42000 …line 1 column 23 near "BY (1)" 
> exec ^
: string
| SPLIT TABLE t1 INDEX primary BY (1)
< error
| 1064 42000 …line 1 column 28 near "primary BY (1)" 
> exec ^
: string
| SPLIT TABLE t1 PARTITION INDEX idx1 BY ('a')
< error
| 1064 42000 …line 1 column 30 near "INDEX idx1 BY ('a')" 
> exec ^
: string
| SPLIT TABLE t1 PARTITION () INDEX idx1 BY ('a')
< error
| 1064 42000 …line 1 column 27 near ") INDEX idx1 BY ('a')" 
> exec ^
: string
| SPLIT TABLE t1 INDEX idx1 PARTITION (p0) BY ('a')
< error
| 1064 42000 …line 1 column 35 near "PARTITION (p0) BY ('a')" 
> exec ^
: string
| SPLIT REGION FOR t1 BY ('a')
< error
| 1064 42000 …line 1 column 19 near "t1 BY ('a')" 
> exec ^
: string
| SPLIT REGION TABLE t1 BY ('a')
< error
| 1064 42000 …line 1 column 18 near "TABLE t1 BY ('a')" 
> exec ^
: string
| SPLIT REGION PARTITION TABLE t1 BY ('a')
< error
| 1064 42000 …line 1 column 22 near "PARTITION TABLE t1 BY ('a')" 
> exec ^
: string
| SPLIT PARTITION REGION FOR TABLE t1 BY ('a')
< error
| 1064 42000 …line 1 column 22 near "REGION FOR TABLE t1 BY ('a')" 
> exec ^
: string
| SPLIT PARTITION t1 BY ('a')
< error
| 1064 42000 …line 1 column 18 near "t1 BY ('a')" 
> exec ^
: string
| SPLIT REGION FOR PARTITION t1 BY ('a')
< error
| 1064 42000 …line 1 column 29 near "t1 BY ('a')" 
> exec ^
: string
| SPLIT t1 BY ('a')
< error
| 1064 42000 …line 1 column 8 near "t1 BY ('a')" 
> exec ^
: string
| SPLIT TABLE t1 BY (maxvalue)
< error
| 1064 42000 …line 1 column 27 near "maxvalue)" 
> exec ^
: string
| DISTRIBUTE TABLE t1
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| DISTRIBUTE TABLE t1 PARTITION(p0)
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| DISTRIBUTE TABLE t1 RULE = 'leader-scatter'
< error
| 1064 42000 …line 1 column 43 near "" 
> exec ^
: string
| DISTRIBUTE TABLE t1 ENGINE = 'tikv'
< error
| 1064 42000 …line 1 column 26 near "ENGINE = 'tikv'" 
> exec ^
: string
| DISTRIBUTE TABLE t1 TIMEOUT = '30m'
< error
| 1064 42000 …line 1 column 27 near "TIMEOUT = '30m'" 
> exec ^
: string
| DISTRIBUTE TABLE t1 ENGINE = 'tikv' RULE = 'leader-scatter'
< error
| 1064 42000 …line 1 column 26 near "ENGINE = 'tikv' RULE = 'leader-scatter'" 
> exec ^
: string
| DISTRIBUTE TABLE t1 TIMEOUT = '30m' RULE = 'x' ENGINE = 'tikv'
< error
| 1064 42000 …line 1 column 27 near "TIMEOUT = '30m' RULE = 'x' ENGINE = 'tikv'" 
> exec ^
: string
| DISTRIBUTE TABLE t1 RULE = 'x' ENGINE = 'tikv' TIMEOUT = '30m' RULE = 'y'
< error
| 1064 42000 …line 1 column 67 near "RULE = 'y'" 
> exec ^
: string
| DISTRIBUTE TABLE t1 RULE = leader ENGINE = 'tikv'
< error
| 1064 42000 …line 1 column 33 near "leader ENGINE = 'tikv'" 
> exec ^
: string
| DISTRIBUTE TABLE t1 RULE = 1 ENGINE = 'tikv'
< error
| 1064 42000 …line 1 column 28 near "1 ENGINE = 'tikv'" 
> exec ^
: string
| DISTRIBUTE TABLE t1 RULE = 'x' ENGINE = tikv
< error
| 1064 42000 …line 1 column 44 near "tikv" 
> exec ^
: string
| DISTRIBUTE TABLE t1 RULE = 'x' ENGINE = 'tikv' TIMEOUT = 30
< error
| 1064 42000 …line 1 column 59 near "30" 
> exec ^
: string
| DISTRIBUTE TABLE t1 RULE='x',ENGINE='tikv'
< error
| 1064 42000 …line 1 column 29 near ",ENGINE='tikv'" 
> exec ^
: string
| DISTRIBUTE TABLE t1 RULE = 'x' ENGINE = 'tikv' EXTRA = 'y'
< error
| 1064 42000 …line 1 column 52 near "EXTRA = 'y'" 
> exec ^
: string
| DISTRIBUTE TABLE t1 INDEX idx1 RULE = 'x' ENGINE = 'tikv'
< error
| 1064 42000 …line 1 column 25 near "INDEX idx1 RULE = 'x' ENGINE = 'tikv'" 
> exec ^
: string
| DISTRIBUTE TABLE t1, t2 RULE = 'x' ENGINE = 'tikv'
< error
| 1064 42000 …line 1 column 20 near ", t2 RULE = 'x' ENGINE = 'tikv'" 
> exec ^
: string
| DISTRIBUTE TABLE t1 PARTITION() RULE = 'x' ENGINE = 'tikv'
< error
| 1064 42000 …line 1 column 31 near ") RULE = 'x' ENGINE = 'tikv'" 
> exec ^
: string
| DISTRIBUTE TABLE t1 PARTITION RULE = 'x' ENGINE = 'tikv'
< error
| 1064 42000 …line 1 column 34 near "RULE = 'x' ENGINE = 'tikv'" 
> exec ^
: string
| DISTRIBUTE t1 RULE = 'x' ENGINE = 'tikv'
< error
| 1064 42000 …line 1 column 13 near "t1 RULE = 'x' ENGINE = 'tikv'" 
> exec ^
: string
| DISTRIBUTE TABLE t1 RULE = 'x' ENGINE = 'tikv' TIMEOUT = '1m' TIMEOUT = '2m'
< error
| 1064 42000 …line 1 column 69 near "TIMEOUT = '2m'" 
