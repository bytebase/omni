-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestMariaDBRefusesMySQLHistograms default#1
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON b
< error
| 1064 42000 …near 'UPDATE HISTOGRAM ON b' at line 1
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON b WITH 4 BUCKETS
< error
| 1064 42000 …near 'UPDATE HISTOGRAM ON b WITH 4 BUCKETS' at line 1
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON b USING DATA '{}'
< error
| 1064 42000 …near 'UPDATE HISTOGRAM ON b USING DATA '{}'' at line 1
> exec ^
: string
| ANALYZE TABLE t DROP HISTOGRAM ON b
< error
| 1064 42000 …near 'DROP HISTOGRAM ON b' at line 1
