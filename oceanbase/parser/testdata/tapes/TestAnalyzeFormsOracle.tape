-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestAnalyzeFormsOracle default#1
> exec
| ANALYZE TABLE t
< error
| 1146 42S02 Table doesn't exist
> exec
| ANALYZE TABLE t1, t2
< error
| 1146 42S02 Table doesn't exist
> exec
| ANALYZE LOCAL TABLE t
< error
| 1064 42000 …near 'LOCAL TABLE t' at line 1
> exec
| ANALYZE NO_WRITE_TO_BINLOG TABLE t
< error
| 1064 42000 …near 'NO_WRITE_TO_BINLOG TABLE t' at line 1
> exec
| ANALYZE t
< error
| 1064 42000 …near 't' at line 1
> exec
| ANALYZE LOCAL t
< error
| 1064 42000 …near 'LOCAL t' at line 1
> exec
| ANALYZE NO_WRITE_TO_BINLOG t
< error
| 1064 42000 …near 'NO_WRITE_TO_BINLOG t' at line 1
> exec
| ANALYZE TABLE
< error
| 1064 42000 …near 'TABLE' at line 1
> exec
| ANALYZE
< error
| 1064 42000 …near 'ANALYZE' at line 1
> exec
| ANALYZE TABLE t WITH 4 BUCKETS
< error
| 1064 42000 …near 'WITH 4 BUCKETS' at line 1
> exec
| ANALYZE TABLE t WITH 4 TOPN
< error
| 1064 42000 …near 'WITH 4 TOPN' at line 1
> exec
| ANALYZE TABLE t WITH 0.5 SAMPLERATE
< error
| 1064 42000 …near 'WITH 0.5 SAMPLERATE' at line 1
> exec
| ANALYZE TABLE t INDEX a
< error
| 1064 42000 …near 'INDEX a' at line 1
> exec
| ANALYZE TABLE t INDEX
< error
| 1064 42000 …near 'INDEX' at line 1
> exec
| ANALYZE TABLE t ALL COLUMNS
< error
| 1064 42000 …near 'ALL COLUMNS' at line 1
> exec
| ANALYZE TABLE t PREDICATE COLUMNS
< error
| 1064 42000 …near 'PREDICATE COLUMNS' at line 1
> exec
| ANALYZE TABLE t COLUMNS c1, c2
< error
| 1064 42000 …near 'COLUMNS c1, c2' at line 1
> exec
| ANALYZE TABLE t PARTITION p0
< error
| 1064 42000 …near 'p0' at line 1
> exec
| ANALYZE INCREMENTAL TABLE t INDEX a
< error
| 1064 42000 …near 'INCREMENTAL TABLE t INDEX a' at line 1
> exec
| DROP STATS t
< error
| 1064 42000 …near 'STATS t' at line 1
> exec
| LOAD STATS '/tmp/stats.json'
< error
| 1064 42000 …near 'STATS '/tmp/stats.json'' at line 1
> exec
| CREATE STATISTICS s1 (CARDINALITY) ON t(a)
< error
| 1064 42000 …near 'STATISTICS s1 (CARDINALITY) ON t(a)' at line 1
> exec
| DROP STATISTICS s1
< error
| 1064 42000 …near 'STATISTICS s1' at line 1
> exec
| ANALYZE TABLE t UPDATE HISTOGRAM ON b
< error
| 1146 42S02 Table doesn't exist
> exec
| ANALYZE TABLE t UPDATE HISTOGRAM ON b WITH 4 BUCKETS
< error
| 1146 42S02 Table doesn't exist
> exec
| ANALYZE TABLE t UPDATE HISTOGRAM ON b USING DATA '{}'
< error
| 1064 42000 …near 'USING DATA '{}'' at line 1
> exec
| ANALYZE TABLE t DROP HISTOGRAM ON b
< error
| 1146 42S02 Table doesn't exist
> exec
| ANALYZE TABLE t UPDATE HISTOGRAM ON b WITH 4 TOPN
< error
| 1064 42000 …near 'TOPN' at line 1
