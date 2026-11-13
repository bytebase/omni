-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestAnalyzeFormsOracle default#1
> exec
| PREPARE omni_parses_stmt FROM ?
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
| ANALYZE TABLE t UPDATE HISTOGRAM ON b WITH 4 TOPN
< error
| 1064 42000 …near 'UPDATE HISTOGRAM ON b WITH 4 TOPN' at line 1
