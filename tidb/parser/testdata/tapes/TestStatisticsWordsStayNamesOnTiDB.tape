-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestStatisticsWordsStayNamesOnTiDB default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| SELECT * FROM stats
< error
| 1146 42S02 Table 'test.stats' doesn't exist
> exec ^
: string
| SELECT * FROM stats
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM statistics
< error
| 1146 42S02 Table 'test.statistics' doesn't exist
> exec ^
: string
| SELECT * FROM statistics
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM topn
< error
| 1146 42S02 Table 'test.topn' doesn't exist
> exec ^
: string
| SELECT * FROM topn
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM cmsketch
< error
| 1146 42S02 Table 'test.cmsketch' doesn't exist
> exec ^
: string
| SELECT * FROM cmsketch
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM samples
< error
| 1146 42S02 Table 'test.samples' doesn't exist
> exec ^
: string
| SELECT * FROM samples
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM samplerate
< error
| 1146 42S02 Table 'test.samplerate' doesn't exist
> exec ^
: string
| SELECT * FROM samplerate
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM cardinality
< error
| 1146 42S02 Table 'test.cardinality' doesn't exist
> exec ^
: string
| SELECT * FROM cardinality
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM dependency
< error
| 1146 42S02 Table 'test.dependency' doesn't exist
> exec ^
: string
| SELECT * FROM dependency
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM correlation
< error
| 1146 42S02 Table 'test.correlation' doesn't exist
> exec ^
: string
| SELECT * FROM correlation
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM incremental
< error
| 1146 42S02 Table 'test.incremental' doesn't exist
> exec ^
: string
| SELECT * FROM incremental
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM predicate
< error
| 1146 42S02 Table 'test.predicate' doesn't exist
> exec ^
: string
| SELECT * FROM predicate
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM width
< error
| 1146 42S02 Table 'test.width' doesn't exist
> exec ^
: string
| SELECT * FROM width
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM depth
< error
| 1146 42S02 Table 'test.depth' doesn't exist
> exec ^
: string
| SELECT * FROM depth
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM buckets
< error
| 1146 42S02 Table 'test.buckets' doesn't exist
> exec ^
: string
| SELECT * FROM buckets
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t stats WHERE stats.id = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t stats WHERE stats.id = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t topn WHERE topn.id = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t topn WHERE topn.id = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t cmsketch WHERE cmsketch.id = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t cmsketch WHERE cmsketch.id = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t samplerate WHERE samplerate.id = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t samplerate WHERE samplerate.id = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t incremental WHERE incremental.id = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t incremental WHERE incremental.id = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t predicate WHERE predicate.id = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t predicate WHERE predicate.id = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT 1 AS statistics
< ok 0 0
> exec ^
: string
| SELECT 1 AS samples
< ok 0 0
> exec ^
: string
| SELECT 1 AS correlation
< ok 0 0
> exec ^
: string
| SELECT depth, width, topn FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT depth, width, topn FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT t.cardinality FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT t.cardinality FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE stats (samples int, width int, depth int, topn int)
< ok 0 0
> exec ^
: string
| CREATE TABLE statistics (cardinality int, dependency int, correlation int)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (samplerate double, incremental int, predicate int, INDEX cmsketch (incremental))
< ok 0 0
> exec ^
: string
| CREATE INDEX topn ON t (a)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX topn ON t (a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DROP TABLE stats
< ok 0 0
> exec ^
: string
| DROP TABLE statistics
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD COLUMN samplerate double
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN samplerate double
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ANALYZE TABLE stats
< error
| 1146 42S02 Table 'test.stats' doesn't exist
> exec ^
: string
| ANALYZE TABLE stats
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ANALYZE TABLE statistics COLUMNS samples, width
< error
| 1146 42S02 Table 'test.statistics' doesn't exist
> exec ^
: string
| ANALYZE TABLE statistics COLUMNS samples, width
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ANALYZE TABLE t INDEX topn WITH 4 BUCKETS
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ANALYZE TABLE t INDEX topn WITH 4 BUCKETS
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DROP STATS stats
< error
| 1146 42S02 Table 'test.stats' doesn't exist
> exec ^
: string
| DROP STATS stats
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DROP STATISTICS statistics
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| DROP STATISTICS statistics
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE STATISTICS statistics (CARDINALITY) ON stats(samples)
< error
| 1146 42S02 Table 'test.stats' doesn't exist
> exec ^
: string
| CREATE STATISTICS statistics (CARDINALITY) ON stats(samples)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
