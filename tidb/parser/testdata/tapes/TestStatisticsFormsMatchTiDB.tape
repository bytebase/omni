-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestStatisticsFormsMatchTiDB default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| analyze table t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze local table t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze local table t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze no_write_to_binlog table t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze no_write_to_binlog table t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze local incremental table t index a
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze local incremental table t index a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze no_write_to_binlog incremental table t index a
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze no_write_to_binlog incremental table t index a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze incremental local table t index a
< error
| 1064 42000 …line 1 column 25 near "local table t index a" 
> exec ^
: string
| analyze t
< error
| 1064 42000 …line 1 column 9 near "t" 
> exec ^
: string
| analyze table
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| analyze `incremental` table t index a
< error
| 1064 42000 …line 1 column 21 near "`incremental` table t index a" 
> exec ^
: string
| analyze table t all columns
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t all columns
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t predicate columns
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t predicate columns
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t columns c1,c2
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t columns c1,c2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t columns `c1`, `c2`
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t columns `c1`, `c2`
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t1,t2 all columns
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| analyze table t1,t2 all columns
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t1,t2 predicate columns
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| analyze table t1,t2 predicate columns
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t1,t2 columns c1
< error
| 1064 42000 …line 1 column 27 near "columns c1" 
> exec ^
: string
| analyze table t columns
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| analyze table t columns c1,
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| analyze table t columns c1.c2
< error
| 1064 42000 …line 1 column 27 near ".c2" 
> exec ^
: string
| analyze table t columns t.c1,t.c2
< error
| 1064 42000 …line 1 column 26 near ".c1,t.c2" 
> exec ^
: string
| analyze table t `all` columns
< error
| 1064 42000 …line 1 column 21 near "`all` columns" 
> exec ^
: string
| analyze table t all `columns`
< error
| 1064 42000 …line 1 column 29 near "`columns`" 
> exec ^
: string
| analyze table t columns c1,c2 all columns
< error
| 1064 42000 …line 1 column 33 near "all columns" 
> exec ^
: string
| analyze table t all
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| analyze table t predicate
< error
| 1064 42000 …line 1 column 25 near "" 
> exec ^
: string
| analyze table t partition p0 partition p1
< error
| 1064 42000 …line 1 column 38 near "partition p1" 
> exec ^
: string
| analyze table t index
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t index
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t index a
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t index a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t index a,b
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t index a,b
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t index a, b
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t index a, b
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t index `a`
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t index `a`
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t index ,a
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t index ,a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t index ,a,b
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t index ,a,b
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t index primary
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t index primary
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t index ,primary
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t index ,primary
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t index a, primary
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t index a, primary
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t index primary, a
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t index primary, a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t index a,
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| analyze table t index ,,a
< error
| 1064 42000 …line 1 column 24 near ",a" 
> exec ^
: string
| analyze table t index a,,b
< error
| 1064 42000 …line 1 column 25 near ",b" 
> exec ^
: string
| analyze table t index 4
< error
| 1064 42000 …line 1 column 23 near "4" 
> exec ^
: string
| analyze table t1,t2 index a
< error
| 1064 42000 …line 1 column 25 near "index a" 
> exec ^
: string
| analyze table t index a all columns
< error
| 1064 42000 …line 1 column 27 near "all columns" 
> exec ^
: string
| analyze table t index a predicate columns
< error
| 1064 42000 …line 1 column 33 near "predicate columns" 
> exec ^
: string
| analyze table t index a columns c
< error
| 1064 42000 …line 1 column 31 near "columns c" 
> exec ^
: string
| analyze table t all columns index a
< error
| 1064 42000 …line 1 column 33 near "index a" 
> exec ^
: string
| analyze table t columns c1 index a
< error
| 1064 42000 …line 1 column 32 near "index a" 
> exec ^
: string
| analyze table t partition a
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t partition a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t partition p0,p1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t partition p0,p1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t partition `p0`
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t partition `p0`
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t partition a all columns
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t partition a all columns
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t partition a predicate columns
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t partition a predicate columns
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t partition a columns c1,c2
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t partition a columns c1,c2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t partition a index b
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t partition a index b
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t partition p0 index
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t partition p0 index
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t partition p0 index a,b with 4 buckets
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t partition p0 index a,b with 4 buckets
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table db.t partition p0 index b
< error
| 1146 42S02 Table 'db.t' doesn't exist
> exec ^
: string
| analyze table db.t partition p0 index b
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t1,t2 partition p0
< error
| 1064 42000 …line 1 column 29 near "partition p0" 
> exec ^
: string
| analyze table t partition
< error
| 1064 42000 …line 1 column 25 near "" 
> exec ^
: string
| analyze table t partition (p0)
< error
| 1064 42000 …line 1 column 27 near "(p0)" 
> exec ^
: string
| analyze table t partition p0 all columns index a
< error
| 1064 42000 …line 1 column 46 near "index a" 
> exec ^
: string
| analyze table t partition p0 index a all columns
< error
| 1064 42000 …line 1 column 40 near "all columns" 
> exec ^
: string
| analyze incremental table t index
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze incremental table t index
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze incremental table t index idx
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze incremental table t index idx
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze incremental table t partition p0 index b
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze incremental table t partition p0 index b
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze incremental table t
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| analyze incremental table t all columns
< error
| 1064 42000 …line 1 column 31 near "all columns" 
> exec ^
: string
| analyze incremental table t columns c1
< error
| 1064 42000 …line 1 column 35 near "columns c1" 
> exec ^
: string
| analyze incremental table t1,t2 index a
< error
| 1064 42000 …line 1 column 29 near ",t2 index a" 
> exec ^
: string
| analyze incremental table t update histogram on b
< error
| 1064 42000 …line 1 column 34 near "update histogram on b" 
> exec ^
: string
| analyze incremental table t partition p0
< error
| 1064 42000 …line 1 column 40 near "" 
> exec ^
: string
| analyze incremental table t with 4 buckets
< error
| 1064 42000 …line 1 column 32 near "with 4 buckets" 
> exec ^
: string
| analyze incremental table t drop histogram on b
< error
| 1064 42000 …line 1 column 32 near "drop histogram on b" 
> exec ^
: string
| analyze table t with 4 buckets
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t with 4 buckets
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t with 4 topn
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t with 4 topn
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t with 4 cmsketch width
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t with 4 cmsketch width
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t with 4 cmsketch depth
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t with 4 cmsketch depth
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t with 4 samples
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t with 4 samples
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t with 10 samplerate
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t with 10 samplerate
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t with 0.1 samplerate
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t with 0.1 samplerate
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t with 4 BUCKETS
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t with 4 BUCKETS
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t with 4 TopN
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t with 4 TopN
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t with 4 CMSketch Width
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t with 4 CMSketch Width
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t with 4 buckets, 4 topn, 4 cmsketch width, 4 cmsketch depth, 4 samples
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t with 4 buckets, 4 topn, 4 cmsketch width, 4 cmsketch depth, 4 samples
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t with 4 buckets,4 topn
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t with 4 buckets,4 topn
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t with 0.1 samplerate, 4 buckets
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t with 0.1 samplerate, 4 buckets
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t with 4 buckets, 0.5 samplerate
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t with 4 buckets, 0.5 samplerate
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t with 4 buckets, 4 buckets
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t with 4 buckets, 4 buckets
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t with 4 buckets, 8 buckets
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t with 4 buckets, 8 buckets
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t with 4 samplerate, 4 samplerate
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t with 4 samplerate, 4 samplerate
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t with 0 buckets
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t with 0 buckets
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t with 0 topn
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t with 0 topn
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t with 0 samples
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t with 0 samples
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t with 18446744073709551615 buckets
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t with 18446744073709551615 buckets
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t with 18446744073709551616 buckets
< error
| 1064 42000 …line 1 column 49 near "buckets" 
> exec ^
: string
| analyze table t with 9999999999999999999 buckets
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t with 9999999999999999999 buckets
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t with 4 penguins
< error
| 1064 42000 …line 1 column 31 near "penguins" 
> exec ^
: string
| analyze table t with 4 cmsketch
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| analyze table t with 4 width
< error
| 1064 42000 …line 1 column 28 near "width" 
> exec ^
: string
| analyze table t with 4 depth
< error
| 1064 42000 …line 1 column 28 near "depth" 
> exec ^
: string
| analyze table t with 4
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| analyze table t with buckets
< error
| 1064 42000 …line 1 column 28 near "buckets" 
> exec ^
: string
| analyze table t with
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| analyze table t with 1.5 buckets
< error
| 1064 42000 …line 1 column 32 near "buckets" 
> exec ^
: string
| analyze table t with 0.5 topn
< error
| 1064 42000 …line 1 column 29 near "topn" 
> exec ^
: string
| analyze table t with -1 buckets
< error
| 1064 42000 …line 1 column 22 near "-1 buckets" 
> exec ^
: string
| analyze table t with -0.1 samplerate
< error
| 1064 42000 …line 1 column 22 near "-0.1 samplerate" 
> exec ^
: string
| analyze table t with 1.0 samplerate
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t with 1.0 samplerate
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t with .5 samplerate
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t with .5 samplerate
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t with 1e3 samplerate
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t with 1e3 samplerate
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t with 1024 samplerate
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t with 1024 samplerate
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t with 0x10 buckets
< error
| 1064 42000 …line 1 column 25 near "0x10 buckets" 
> exec ^
: string
| analyze table t with 0x10 samplerate
< error
| 1064 42000 …line 1 column 25 near "0x10 samplerate" 
> exec ^
: string
| analyze table t with 4e0 buckets
< error
| 1064 42000 …line 1 column 32 near "buckets" 
> exec ^
: string
| analyze table t with '4' buckets
< error
| 1064 42000 …line 1 column 24 near "'4' buckets" 
> exec ^
: string
| analyze table t with 4 buckets 4 topn
< error
| 1064 42000 …line 1 column 32 near "4 topn" 
> exec ^
: string
| analyze table t index a with 4 buckets with 4 topn
< error
| 1064 42000 …line 1 column 43 near "with 4 topn" 
> exec ^
: string
| analyze table t with 4 buckets,
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| analyze table t with 4 `buckets`
< error
| 1064 42000 …line 1 column 32 near "`buckets`" 
> exec ^
: string
| analyze table t index a with 4 `topn`
< error
| 1064 42000 …line 1 column 37 near "`topn`" 
> exec ^
: string
| analyze table t index a with 4 buckets
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t index a with 4 buckets
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t index a with 4 buckets, 4 topn
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t index a with 4 buckets, 4 topn
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t index a with 0.1 samplerate
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t index a with 0.1 samplerate
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t partition a with 4 buckets
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t partition a with 4 buckets
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t partition a index b with 4 buckets
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t partition a index b with 4 buckets
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t partition p0 index a with 4 buckets
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t partition p0 index a with 4 buckets
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t partition p0 columns c1,c2
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t partition p0 columns c1,c2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t all columns with 1024 buckets
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t all columns with 1024 buckets
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze no_write_to_binlog incremental table t partition p0 index a
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze no_write_to_binlog incremental table t partition p0 index a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t partition a all columns with 1024 buckets
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t partition a all columns with 1024 buckets
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t partition a columns c1,c2 with 1024 buckets
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t partition a columns c1,c2 with 1024 buckets
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t partition a predicate columns with 1024 buckets
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t partition a predicate columns with 1024 buckets
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t all columns with 4 buckets
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t all columns with 4 buckets
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t predicate columns with 4 buckets
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t predicate columns with 4 buckets
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t columns c1,c2 with 4 topn
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t columns c1,c2 with 4 topn
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t columns c1 with 4 buckets, 4 topn
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t columns c1 with 4 buckets, 4 topn
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t1,t2 all columns with 4 topn
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| analyze table t1,t2 all columns with 4 topn
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t1,t2 predicate columns with 4 topn
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| analyze table t1,t2 predicate columns with 4 topn
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t1, t2 all columns with 4 buckets, 0.5 samplerate
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| analyze table t1, t2 all columns with 4 buckets, 0.5 samplerate
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze local table t with 4 buckets
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze local table t with 4 buckets
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t update histogram on b
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t update histogram on b
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t update histogram on c1, c2
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t update histogram on c1, c2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t drop histogram on b
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t drop histogram on b
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t drop histogram on c1, c2
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t drop histogram on c1, c2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze no_write_to_binlog table t update histogram on b
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze no_write_to_binlog table t update histogram on b
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t update histogram on b with 1024 buckets
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t update histogram on b with 1024 buckets
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t update histogram on b with 4 topn
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t update histogram on b with 4 topn
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t update histogram on b with 4 buckets, 4 topn
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t update histogram on b with 4 buckets, 4 topn
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t update histogram on b with 0.5 samplerate
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| analyze table t update histogram on b with 0.5 samplerate
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t drop histogram on b with 4 buckets
< error
| 1064 42000 …line 1 column 40 near "with 4 buckets" 
> exec ^
: string
| analyze table t1,t2 update histogram on b
< error
| 1064 42000 …line 1 column 26 near "update histogram on b" 
> exec ^
: string
| analyze table t1,t2 drop histogram on b
< error
| 1064 42000 …line 1 column 24 near "drop histogram on b" 
> exec ^
: string
| analyze table t update histogram on t.c1
< error
| 1064 42000 …line 1 column 38 near ".c1" 
> exec ^
: string
| analyze table t update histogram on t.c1, t.c2
< error
| 1064 42000 …line 1 column 38 near ".c1, t.c2" 
> exec ^
: string
| analyze table t drop histogram on t.c1, t.c2
< error
| 1064 42000 …line 1 column 36 near ".c1, t.c2" 
> exec ^
: string
| analyze table t partition p0 update histogram on b
< error
| 1064 42000 …line 1 column 35 near "update histogram on b" 
> exec ^
: string
| analyze table t index a update histogram on b
< error
| 1064 42000 …line 1 column 30 near "update histogram on b" 
> exec ^
: string
| analyze table t index a drop histogram on b
< error
| 1064 42000 …line 1 column 28 near "drop histogram on b" 
> exec ^
: string
| analyze table t all columns update histogram on b
< error
| 1064 42000 …line 1 column 34 near "update histogram on b" 
> exec ^
: string
| analyze table t drop table u
< error
| 1064 42000 …line 1 column 26 near "table u" 
> exec ^
: string
| analyze table t partition p0 update histogram on b with 4 buckets
< error
| 1064 42000 …line 1 column 35 near "update histogram on b with 4 buckets" 
> exec ^
: string
| analyze table t update histogram
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| analyze table t drop histogram
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| analyze table t update histogram on b using data '{}'
< error
| 1064 42000 …line 1 column 43 near "using data '{}'" 
> exec ^
: string
| analyze table t update histogram on b, c using data '{}'
< error
| 1064 42000 …line 1 column 46 near "using data '{}'" 
> exec ^
: string
| analyze table t1,t2
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| analyze table t1,t2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t1, t2, t3
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| analyze table t1, t2, t3
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table db1.t1
< error
| 1146 42S02 Table 'db1.t1' doesn't exist
> exec ^
: string
| analyze table db1.t1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| analyze table t1.*
< error
| 1064 42000 …line 1 column 18 near "*" 
> exec ^
: string
| drop stats t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| drop stats t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| drop stats db.t
< error
| 1146 42S02 Table 'db.t' doesn't exist
> exec ^
: string
| drop stats db.t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| drop stats `t`
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| drop stats `t`
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| drop stats t1, t2, t3
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| drop stats t1, t2, t3
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| drop stats t global
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| drop stats t global
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| drop stats t partition p0
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| drop stats t partition p0
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| drop stats t partition p0, p1, p2
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| drop stats t partition p0, p1, p2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| drop stats t1,t2 global
< error
| 1064 42000 …line 1 column 23 near "global" 
> exec ^
: string
| drop stats t1,t2 partition p0
< error
| 1064 42000 …line 1 column 26 near "partition p0" 
> exec ^
: string
| drop stats t global partition p0
< error
| 1064 42000 …line 1 column 29 near "partition p0" 
> exec ^
: string
| drop stats t global, t2
< error
| 1064 42000 …line 1 column 20 near ", t2" 
> exec ^
: string
| drop stats
< error
| 1064 42000 …line 1 column 10 near "" 
> exec ^
: string
| drop stats t partition
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| drop stats t partition p0,
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| drop stats t partition (p0)
< error
| 1064 42000 …line 1 column 24 near "(p0)" 
> exec ^
: string
| drop stats t partition (p0,p1)
< error
| 1064 42000 …line 1 column 24 near "(p0,p1)" 
> exec ^
: string
| drop `stats` t
< error
| 1064 42000 …line 1 column 12 near "`stats` t" 
> exec ^
: string
| load stats '/tmp/stats.json'
< ok 0 0
> exec ^
: string
| load stats "x"
< ok 0 0
> exec ^
: string
| load stats
< error
| 1064 42000 …line 1 column 10 near "" 
> exec ^
: string
| load stats 4
< error
| 1064 42000 …line 1 column 12 near "4" 
> exec ^
: string
| load stats 'a' 'b'
< error
| 1064 42000 …line 1 column 18 near "'b'" 
> exec ^
: string
| load `stats` 'x'
< error
| 1064 42000 …line 1 column 12 near "`stats` 'x'" 
> exec ^
: string
| create statistics stats1 (cardinality) on t(a,b,c)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| create statistics stats1 (cardinality) on t(a,b,c)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create statistics stats1(cardinality) on t(a,b,c)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| create statistics stats1(cardinality) on t(a,b,c)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create statistics stats2 (dependency) on t(a,b)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| create statistics stats2 (dependency) on t(a,b)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create statistics stats3 (correlation) on t(a,b)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| create statistics stats3 (correlation) on t(a,b)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create statistics s1 (CARDINALITY) on t(a)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| create statistics s1 (CARDINALITY) on t(a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create statistics `s1` (cardinality) on t(a)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| create statistics `s1` (cardinality) on t(a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create statistics if not exists stats1 (cardinality) on t(a,b,c)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| create statistics if not exists stats1 (cardinality) on t(a,b,c)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create statistics s1 (cardinality) on db.t(a)
< error
| 1146 42S02 Table 'db.t' doesn't exist
> exec ^
: string
| create statistics s1 (cardinality) on db.t(a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create statistics s1 (cardinality) on t(t.a)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| create statistics s1 (cardinality) on t(t.a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create statistics s1 (cardinality) on t(d.t.a)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| create statistics s1 (cardinality) on t(d.t.a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create statistics s1 (cardinality) on t(a.*)
< error
| 1064 42000 …line 1 column 43 near "*)" 
> exec ^
: string
| create statistics stats3 on t(a,b)
< error
| 1064 42000 …line 1 column 27 near "on t(a,b)" 
> exec ^
: string
| create statistics if not exists stats3 on t(a,b)
< error
| 1064 42000 …line 1 column 41 near "on t(a,b)" 
> exec ^
: string
| create statistics s1 (frequency) on t(a,b)
< error
| 1064 42000 …line 1 column 31 near "frequency) on t(a,b)" 
> exec ^
: string
| create statistics s1 (cardinality,dependency) on t(a)
< error
| 1064 42000 …line 1 column 34 near ",dependency) on t(a)" 
> exec ^
: string
| create statistics s1 (`cardinality`) on t(a)
< error
| 1064 42000 …line 1 column 35 near "`cardinality`) on t(a)" 
> exec ^
: string
| create statistics s1 (cardinality) on t
< error
| 1064 42000 …line 1 column 39 near "" 
> exec ^
: string
| create statistics s1 (cardinality) on t()
< error
| 1064 42000 …line 1 column 41 near ")" 
> exec ^
: string
| create statistics s1 (cardinality) on t(a,)
< error
| 1064 42000 …line 1 column 43 near ")" 
> exec ^
: string
| create `statistics` s1 (cardinality) on t(a)
< error
| 1064 42000 …line 1 column 19 near "`statistics` s1 (cardinality) on t(a)" 
> exec ^
: string
| drop statistics stats1
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| drop statistics stats1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| drop statistics `s1`
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| drop statistics `s1`
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| drop statistics
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| drop statistics s1, s2
< error
| 1064 42000 …line 1 column 19 near ", s2" 
> exec ^
: string
| drop statistics if exists s1
< error
| 1064 42000 …line 1 column 18 near "if exists s1" 
