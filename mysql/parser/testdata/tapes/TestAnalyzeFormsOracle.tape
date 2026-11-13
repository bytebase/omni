-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestAnalyzeFormsOracle default#1
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| ANALYZE TABLE t
< ok 0 0
> exec ^
: string
| ANALYZE TABLE t1, t2
< ok 0 0
> exec ^
: string
| ANALYZE LOCAL TABLE t
< ok 0 0
> exec ^
: string
| ANALYZE NO_WRITE_TO_BINLOG TABLE t
< ok 0 0
> exec ^
: string
| ANALYZE t
< error
| 1064 42000 …near 't' at line 1
> exec ^
: string
| ANALYZE LOCAL t
< error
| 1064 42000 …near 't' at line 1
> exec ^
: string
| ANALYZE NO_WRITE_TO_BINLOG t
< error
| 1064 42000 …near 't' at line 1
> exec ^
: string
| ANALYZE TABLE
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| ANALYZE
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| ANALYZE TABLE t WITH 4 BUCKETS
< error
| 1064 42000 …near 'WITH 4 BUCKETS' at line 1
> exec ^
: string
| ANALYZE TABLE t WITH 4 TOPN
< error
| 1064 42000 …near 'WITH 4 TOPN' at line 1
> exec ^
: string
| ANALYZE TABLE t WITH 0.5 SAMPLERATE
< error
| 1064 42000 …near 'WITH 0.5 SAMPLERATE' at line 1
> exec ^
: string
| ANALYZE TABLE t INDEX a
< error
| 1064 42000 …near 'INDEX a' at line 1
> exec ^
: string
| ANALYZE TABLE t INDEX
< error
| 1064 42000 …near 'INDEX' at line 1
> exec ^
: string
| ANALYZE TABLE t ALL COLUMNS
< error
| 1064 42000 …near 'ALL COLUMNS' at line 1
> exec ^
: string
| ANALYZE TABLE t PREDICATE COLUMNS
< error
| 1064 42000 …near 'PREDICATE COLUMNS' at line 1
> exec ^
: string
| ANALYZE TABLE t COLUMNS c1, c2
< error
| 1064 42000 …near 'COLUMNS c1, c2' at line 1
> exec ^
: string
| ANALYZE TABLE t PARTITION p0
< error
| 1064 42000 …near 'PARTITION p0' at line 1
> exec ^
: string
| ANALYZE INCREMENTAL TABLE t INDEX a
< error
| 1064 42000 …near 'INCREMENTAL TABLE t INDEX a' at line 1
> exec ^
: string
| DROP STATS t
< error
| 1064 42000 …near 'STATS t' at line 1
> exec ^
: string
| LOAD STATS '/tmp/stats.json'
< error
| 1064 42000 …near 'STATS '/tmp/stats.json'' at line 1
> exec ^
: string
| CREATE STATISTICS s1 (CARDINALITY) ON t(a)
< error
| 1064 42000 …near 'STATISTICS s1 (CARDINALITY) ON t(a)' at line 1
> exec ^
: string
| DROP STATISTICS s1
< error
| 1064 42000 …near 'STATISTICS s1' at line 1
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON b
< ok 0 0
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON b WITH 4 BUCKETS
< ok 0 0
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON b USING DATA '{}'
< ok 0 0
> exec ^
: string
| ANALYZE TABLE t DROP HISTOGRAM ON b
< ok 0 0
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON b WITH 4 TOPN
< error
| 1064 42000 …near 'TOPN' at line 1
