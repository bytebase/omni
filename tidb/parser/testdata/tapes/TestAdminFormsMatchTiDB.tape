-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestAdminFormsMatchTiDB default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| ADMIN
< error
| 1064 42000 …line 1 column 5 near "" 
> exec ^
: string
| ADMIN foo
< error
| 1064 42000 …line 1 column 9 near "foo" 
> exec ^
: string
| ADMIN SHOW
< error
| 1064 42000 …line 1 column 10 near "" 
> exec ^
: string
| ADMIN SHOW TELEMETRY
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| ADMIN CREATE WORKLOAD SNAPSHOT
< error
| 1064 42000 …line 1 column 12 near "CREATE WORKLOAD SNAPSHOT" 
> exec ^
: string
| admin show ddl
< ok 0 0
> exec ^
: string
| Admin Show Ddl Jobs
< ok 0 0
> exec ^
: string
| ADMIN SHOW DDL
< ok 0 0
> exec ^
: string
| ADMIN SHOW DDL JOBS
< ok 0 0
> exec ^
: string
| ADMIN SHOW DDL JOBS 20
< ok 0 0
> exec ^
: string
| ADMIN SHOW DDL JOBS 0
< ok 0 0
> exec ^
: string
| ADMIN SHOW DDL JOBS 9223372036854775807
< ok 0 0
> exec ^
: string
| ADMIN SHOW DDL JOBS 9223372036854775808
< error
| 1064 42000 …line 1 column 39 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| ADMIN SHOW DDL JOBS -1
< error
| 1064 42000 …line 1 column 21 near "-1" 
> exec ^
: string
| ADMIN SHOW DDL JOBS 1.0
< error
| 1064 42000 …line 1 column 23 near "1.0" 
> exec ^
: string
| ADMIN SHOW DDL JOBS 1e3
< error
| 1064 42000 …line 1 column 23 near "1e3" 
> exec ^
: string
| ADMIN SHOW DDL JOBS WHERE id > 0
< error
| 1054 42S22 Unknown column 'id' in 'where clause'
> exec ^
: string
| ADMIN SHOW DDL JOBS WHERE id > 0
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN SHOW DDL JOBS 20 WHERE id = 0
< error
| 1054 42S22 Unknown column 'id' in 'where clause'
> exec ^
: string
| ADMIN SHOW DDL JOBS 20 WHERE id = 0
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN SHOW DDL JOBS 20 WHERE id = 0 AND state = 'x'
< error
| 1054 42S22 Unknown column 'id' in 'where clause'
> exec ^
: string
| ADMIN SHOW DDL JOBS 20 WHERE id = 0 AND state = 'x'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN SHOW DDL JOBS WHERE id > 0 AND state IN ('a', 'b')
< error
| 1054 42S22 Unknown column 'id' in 'where clause'
> exec ^
: string
| ADMIN SHOW DDL JOBS WHERE id > 0 AND state IN ('a', 'b')
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN SHOW DDL JOBS 20 WHERE id IN (SELECT 1)
< error
| 1054 42S22 Unknown column 'id' in 'where clause'
> exec ^
: string
| ADMIN SHOW DDL JOBS 20 WHERE id IN (SELECT 1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN SHOW DDL JOBS 20 WHERE 1
< ok 0 0
> exec ^
: string
| ADMIN SHOW DDL JOBS WHERE
< error
| 1064 42000 …line 1 column 25 near "" 
> exec ^
: string
| ADMIN SHOW DDL extra
< error
| 1064 42000 …line 1 column 20 near "extra" 
> exec ^
: string
| ADMIN SHOW DDL `jobs`
< error
| 1064 42000 …line 1 column 21 near "`jobs`" 
> exec ^
: string
| ADMIN SHOW DDL JOBS 20 LIMIT 5
< error
| 1064 42000 …line 1 column 28 near "LIMIT 5" 
> exec ^
: string
| ADMIN SHOW DDL JOBS 20 OFFSET 2
< error
| 1064 42000 …line 1 column 29 near "OFFSET 2" 
> exec ^
: string
| ADMIN SHOW DDL JOBS 20 WHERE id = 0 LIMIT 1
< error
| 1064 42000 …line 1 column 41 near "LIMIT 1" 
> exec ^
: string
| ADMIN SHOW DDL JOB QUERIES 1
< ok 0 0
> exec ^
: string
| ADMIN SHOW DDL JOB QUERIES 1, 2, 3, 4
< ok 0 0
> exec ^
: string
| ADMIN SHOW DDL JOB QUERIES
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| ADMIN SHOW DDL JOB QUERIES 1,
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| ADMIN SHOW DDL JOB QUERIES 9223372036854775808
< error
| 1064 42000 …line 1 column 46 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| ADMIN SHOW DDL JOB QUERIES LIMIT 5
< ok 0 0
> exec ^
: string
| ADMIN SHOW DDL JOB QUERIES LIMIT 0
< ok 0 0
> exec ^
: string
| ADMIN SHOW DDL JOB QUERIES LIMIT 5, 10
< ok 0 0
> exec ^
: string
| ADMIN SHOW DDL JOB QUERIES LIMIT 3 OFFSET 2
< ok 0 0
> exec ^
: string
| ADMIN SHOW DDL JOB QUERIES LIMIT 22 OFFSET 0
< ok 0 0
> exec ^
: string
| ADMIN SHOW DDL JOB QUERIES LIMIT 18446744073709551615
< ok 0 0
> exec ^
: string
| ADMIN SHOW DDL JOB QUERIES LIMIT 18446744073709551616
< error
| 1064 42000 …line 1 column 53 near "18446744073709551616" 
> exec ^
: string
| ADMIN SHOW DDL JOB QUERIES LIMIT
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| ADMIN SHOW DDL JOB QUERIES LIMIT -1
< error
| 1064 42000 …line 1 column 34 near "-1" 
> exec ^
: string
| ADMIN SHOW DDL JOB QUERIES LIMIT 5 OFFSET
< error
| 1064 42000 …line 1 column 41 near "" 
> exec ^
: string
| ADMIN SHOW DDL JOB QUERIES LIMIT 5 OFFSET 1 extra
< error
| 1064 42000 …line 1 column 49 near "extra" 
> exec ^
: string
| ADMIN SHOW DDL JOB QUERIES 1 LIMIT 2
< error
| 1064 42000 …line 1 column 34 near "LIMIT 2" 
> exec ^
: string
| ADMIN SHOW DDL JOB QUERIES LIMIT 2 WHERE id = 1
< error
| 1064 42000 …line 1 column 40 near "WHERE id = 1" 
> exec ^
: string
| ADMIN SHOW DDL JOB 1
< error
| 1064 42000 …line 1 column 20 near "1" 
> exec ^
: string
| ADMIN SHOW t1 NEXT_ROW_ID
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| ADMIN SHOW t1 NEXT_ROW_ID
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN SHOW test.t1 NEXT_ROW_ID
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| ADMIN SHOW test.t1 NEXT_ROW_ID
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN SHOW ddl NEXT_ROW_ID
< error
| 1146 42S02 Table 'test.ddl' doesn't exist
> exec ^
: string
| ADMIN SHOW ddl NEXT_ROW_ID
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN SHOW slow NEXT_ROW_ID
< error
| 1146 42S02 Table 'test.slow' doesn't exist
> exec ^
: string
| ADMIN SHOW slow NEXT_ROW_ID
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN SHOW bdr NEXT_ROW_ID
< error
| 1146 42S02 Table 'test.bdr' doesn't exist
> exec ^
: string
| ADMIN SHOW bdr NEXT_ROW_ID
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN SHOW jobs NEXT_ROW_ID
< error
| 1146 42S02 Table 'test.jobs' doesn't exist
> exec ^
: string
| ADMIN SHOW jobs NEXT_ROW_ID
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN SHOW ddl.t NEXT_ROW_ID
< error
| 1146 42S02 Table 'ddl.t' doesn't exist
> exec ^
: string
| ADMIN SHOW ddl.t NEXT_ROW_ID
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN SHOW slow.t NEXT_ROW_ID
< error
| 1146 42S02 Table 'slow.t' doesn't exist
> exec ^
: string
| ADMIN SHOW slow.t NEXT_ROW_ID
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN SHOW bdr.t NEXT_ROW_ID
< error
| 1146 42S02 Table 'bdr.t' doesn't exist
> exec ^
: string
| ADMIN SHOW bdr.t NEXT_ROW_ID
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN SHOW `ddl` NEXT_ROW_ID
< error
| 1146 42S02 Table 'test.ddl' doesn't exist
> exec ^
: string
| ADMIN SHOW `ddl` NEXT_ROW_ID
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN SHOW `next_row_id` NEXT_ROW_ID
< error
| 1146 42S02 Table 'test.next_row_id' doesn't exist
> exec ^
: string
| ADMIN SHOW `next_row_id` NEXT_ROW_ID
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN SHOW next_row_id NEXT_ROW_ID
< error
| 1146 42S02 Table 'test.next_row_id' doesn't exist
> exec ^
: string
| ADMIN SHOW next_row_id NEXT_ROW_ID
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN SHOW t1
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| ADMIN SHOW NEXT_ROW_ID
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| ADMIN SHOW t1 NEXT_ROW_ID extra
< error
| 1064 42000 …line 1 column 31 near "extra" 
> exec ^
: string
| ADMIN SHOW `ddl`
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| ADMIN SHOW ddl `next_row_id`
< error
| 1064 42000 …line 1 column 28 near "`next_row_id`" 
> exec ^
: string
| ADMIN SHOW DDL JOBS NEXT_ROW_ID
< error
| 1064 42000 …line 1 column 31 near "NEXT_ROW_ID" 
> exec ^
: string
| ADMIN SHOW t.next_row_id
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| ADMIN SHOW SLOW TOP ALL NEXT_ROW_ID
< error
| 1064 42000 …line 1 column 35 near "NEXT_ROW_ID" 
> exec ^
: string
| ADMIN SHOW SLOW RECENT 11
< ok 0 0
> exec ^
: string
| ADMIN SHOW SLOW TOP 3
< ok 0 0
> exec ^
: string
| ADMIN SHOW SLOW TOP INTERNAL 7
< ok 0 0
> exec ^
: string
| ADMIN SHOW SLOW TOP ALL 9
< ok 0 0
> exec ^
: string
| ADMIN SHOW SLOW
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| ADMIN SHOW SLOW TOP
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| ADMIN SHOW SLOW RECENT
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| ADMIN SHOW SLOW RECENT -1
< error
| 1064 42000 …line 1 column 24 near "-1" 
> exec ^
: string
| ADMIN SHOW SLOW TOP foo 3
< error
| 1064 42000 …line 1 column 23 near "foo 3" 
> exec ^
: string
| ADMIN SHOW SLOW TOP 3 extra
< error
| 1064 42000 …line 1 column 27 near "extra" 
> exec ^
: string
| ADMIN SHOW `slow` RECENT 11
< error
| 1064 42000 …line 1 column 24 near "RECENT 11" 
> exec ^
: string
| ADMIN SHOW SLOW TOP ALL 9 WHERE id = 1
< error
| 1064 42000 …line 1 column 31 near "WHERE id = 1" 
> exec ^
: string
| ADMIN SHOW BDR ROLE
< ok 0 0
> exec ^
: string
| ADMIN SET BDR ROLE PRIMARY
< ok 0 0
> exec ^
: string
| ADMIN SET BDR ROLE primary
< ok 0 0
> exec ^
: string
| ADMIN SET BDR ROLE SECONDARY
< ok 0 0
> exec ^
: string
| ADMIN UNSET BDR ROLE
< ok 0 0
> exec ^
: string
| ADMIN SHOW BDR ROLE extra
< error
| 1064 42000 …line 1 column 25 near "extra" 
> exec ^
: string
| ADMIN SET BDR ROLE
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| ADMIN SET BDR ROLE foo
< error
| 1064 42000 …line 1 column 22 near "foo" 
> exec ^
: string
| ADMIN SET BDR ROLE `primary`
< error
| 1064 42000 …line 1 column 28 near "`primary`" 
> exec ^
: string
| ADMIN SET BDR PRIMARY
< error
| 1064 42000 …line 1 column 21 near "PRIMARY" 
> exec ^
: string
| ADMIN UNSET BDR
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| ADMIN CHECK TABLE t1
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| ADMIN CHECK TABLE t1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN CHECK TABLE t1, t2
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| ADMIN CHECK TABLE t1, t2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN CHECK TABLE test.t1
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| ADMIN CHECK TABLE test.t1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN CHECK TABLE `table`
< error
| 1146 42S02 Table 'test.table' doesn't exist
> exec ^
: string
| ADMIN CHECK TABLE `table`
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN CHECK TABLE
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| ADMIN CHECK TABLES t1
< error
| 1064 42000 …line 1 column 18 near "TABLES t1" 
> exec ^
: string
| ADMIN CHECK TABLE t1 QUICK
< error
| 1064 42000 …line 1 column 26 near "QUICK" 
> exec ^
: string
| ADMIN CHECK TABLE t1 FOR UPGRADE
< error
| 1064 42000 …line 1 column 24 near "FOR UPGRADE" 
> exec ^
: string
| ADMIN CHECKSUM TABLE t1
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| ADMIN CHECKSUM TABLE t1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN CHECKSUM TABLE t1, t2
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| ADMIN CHECKSUM TABLE t1, t2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN CHECKSUM TABLE
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| ADMIN CHECKSUM t1
< error
| 1064 42000 …line 1 column 17 near "t1" 
> exec ^
: string
| ADMIN CHECKSUM TABLE t1 QUICK
< error
| 1064 42000 …line 1 column 29 near "QUICK" 
> exec ^
: string
| CHECK TABLE t1
< error
| 1064 42000 …line 1 column 5 near "CHECK TABLE t1" 
> exec ^
: string
| CHECKSUM TABLE t1
< error
| 1064 42000 …line 1 column 8 near "CHECKSUM TABLE t1" 
> exec ^
: string
| ADMIN CHECK INDEX t idx
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ADMIN CHECK INDEX t idx
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN CHECK INDEX t.a idx
< error
| 1146 42S02 Table 't.a' doesn't exist
> exec ^
: string
| ADMIN CHECK INDEX t.a idx
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN CHECK INDEX `t` `idx`
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ADMIN CHECK INDEX `t` `idx`
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN CHECK INDEX t `select`
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ADMIN CHECK INDEX t `select`
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN CHECK INDEX next_row_id idx
< error
| 1146 42S02 Table 'test.next_row_id' doesn't exist
> exec ^
: string
| ADMIN CHECK INDEX next_row_id idx
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN CHECK INDEX t select
< error
| 1064 42000 …line 1 column 26 near "select" 
> exec ^
: string
| ADMIN CHECK INDEX t
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| ADMIN CHECK INDEX t idx (1, 2)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ADMIN CHECK INDEX t idx (1, 2)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN CHECK INDEX t idx (1, 2), (4, 5)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ADMIN CHECK INDEX t idx (1, 2), (4, 5)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN CHECK INDEX t idx (1)
< error
| 1064 42000 …line 1 column 27 near ")" 
> exec ^
: string
| ADMIN CHECK INDEX t idx ()
< error
| 1064 42000 …line 1 column 26 near ")" 
> exec ^
: string
| ADMIN CHECK INDEX t idx 1
< error
| 1064 42000 …line 1 column 25 near "1" 
> exec ^
: string
| ADMIN CHECK INDEX t idx (-1, 2)
< error
| 1064 42000 …line 1 column 26 near "-1, 2)" 
> exec ^
: string
| ADMIN CHECK INDEX t idx (1, 2),
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| ADMIN RECOVER INDEX t1 idx_a
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| ADMIN RECOVER INDEX t1 idx_a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN RECOVER INDEX t1
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| ADMIN RECOVER TABLE t1
< error
| 1064 42000 …line 1 column 19 near "TABLE t1" 
> exec ^
: string
| ADMIN RECOVER INDEX t1 idx_a (1, 2)
< error
| 1064 42000 …line 1 column 30 near "(1, 2)" 
> exec ^
: string
| ADMIN CLEANUP INDEX t1 idx_a
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| ADMIN CLEANUP INDEX t1 idx_a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN CLEANUP INDEX t1
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| ADMIN CLEANUP INDEX t1 idx_a (1, 2)
< error
| 1064 42000 …line 1 column 30 near "(1, 2)" 
> exec ^
: string
| ADMIN CLEANUP TABLE LOCK t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ADMIN CLEANUP TABLE LOCK t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN CLEANUP TABLE LOCK t1,t2
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| ADMIN CLEANUP TABLE LOCK t1,t2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN CLEANUP TABLE LOCK db.t
< error
| 1146 42S02 Table 'db.t' doesn't exist
> exec ^
: string
| ADMIN CLEANUP TABLE LOCK db.t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN CLEANUP TABLE LOCK
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| ADMIN CLEANUP TABLE t
< error
| 1064 42000 …line 1 column 21 near "t" 
> exec ^
: string
| ADMIN CANCEL DDL JOBS 1
< ok 0 0
> exec ^
: string
| ADMIN CANCEL DDL JOBS 1, 2
< ok 0 0
> exec ^
: string
| ADMIN CANCEL DDL JOBS 0
< ok 0 0
> exec ^
: string
| ADMIN PAUSE DDL JOBS 5
< ok 0 0
> exec ^
: string
| ADMIN PAUSE DDL JOBS 1, 3
< ok 0 0
> exec ^
: string
| ADMIN RESUME DDL JOBS 3
< ok 0 0
> exec ^
: string
| ADMIN RESUME DDL JOBS 1, 2
< ok 0 0
> exec ^
: string
| ADMIN CANCEL DDL JOBS
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| ADMIN CANCEL DDL JOBS 1,
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| ADMIN CANCEL DDL JOBS str_not_num
< error
| 1064 42000 …line 1 column 33 near "str_not_num" 
> exec ^
: string
| ADMIN CANCEL DDL JOBS 0x1
< error
| 1064 42000 …line 1 column 25 near "0x1" 
> exec ^
: string
| ADMIN PAUSE DDL JOBS
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| ADMIN PAUSE DDL JOBS str_not_num
< error
| 1064 42000 …line 1 column 32 near "str_not_num" 
> exec ^
: string
| ADMIN RESUME DDL JOBS
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| ADMIN RESUME DDL JOBS str_not_num
< error
| 1064 42000 …line 1 column 33 near "str_not_num" 
> exec ^
: string
| ADMIN CANCEL DDL JOB 1
< error
| 1064 42000 …line 1 column 20 near "JOB 1" 
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = 2
< ok 0 0
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 thread = 2
< ok 0 0
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 BATCH_SIZE = 3
< error
| 1105 HY000 the value 3 for batch_size is out of range [32, 10240]
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 BATCH_SIZE = 3
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 MAX_WRITE_SPEED = 4
< ok 0 0
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 MAX_WRITE_SPEED = _UTF8MB4'4MiB'
< ok 0 0
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = 2, BATCH_SIZE = 3
< error
| 1105 HY000 the value 3 for batch_size is out of range [32, 10240]
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = 2, BATCH_SIZE = 3
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 `thread` = 2
< ok 0 0
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 foo = 1
< error
| 1105 HY000 unsupported admin alter ddl jobs config: foo
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 foo = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = -2
< error
| 1105 HY000 the value -2 for thread is out of range [1, 256]
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = -2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = +2
< ok 0 0
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = 1.5
< error
| 1105 HY000 the value for thread is invalid, only integer is allowed
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = 1.5
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = 'x'
< error
| 1105 HY000 the value for thread is invalid, only integer is allowed
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = 'x'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = 'a' 'b'
< error
| 1105 HY000 the value for thread is invalid, only integer is allowed
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = 'a' 'b'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = N'1'
< error
| 1105 HY000 the value for thread is invalid, only integer is allowed
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = N'1'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = 0xFF
< error
| 1105 HY000 the value for thread is invalid, only integer is allowed
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = 0xFF
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = 0b101
< error
| 1105 HY000 the value for thread is invalid, only integer is allowed
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = 0b101
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = x'41'
< error
| 1105 HY000 the value for thread is invalid, only integer is allowed
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = x'41'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = b'1'
< error
| 1105 HY000 the value for thread is invalid, only integer is allowed
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = b'1'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = _utf8mb4 x'41'
< error
| 1105 HY000 the value for thread is invalid, only integer is allowed
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = _utf8mb4 x'41'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = TRUE
< ok 0 0
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = FALSE
< error
| 1105 HY000 the value 0 for thread is out of range [1, 256]
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = FALSE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = NULL
< error
| 1105 HY000 the value for thread is invalid, only integer is allowed
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = NULL
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN ALTER DDL JOBS 1
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD =
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = 2 BATCH_SIZE = 3
< error
| 1064 42000 …line 1 column 44 near "BATCH_SIZE = 3" 
> exec ^
: string
| ADMIN ALTER DDL JOBS 1, 2 THREAD = 2
< error
| 1064 42000 …line 1 column 23 near ", 2 THREAD = 2" 
> exec ^
: string
| ADMIN ALTER DDL JOB 1 THREAD = 2
< error
| 1064 42000 …line 1 column 19 near "JOB 1 THREAD = 2" 
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = @a
< error
| 1064 42000 …line 1 column 34 near "@a" 
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = (1)
< error
| 1064 42000 …line 1 column 33 near "(1)" 
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = 1+1
< error
| 1064 42000 …line 1 column 34 near "+1" 
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = DATE '2020-01-01'
< error
| 1064 42000 …line 1 column 36 near "DATE '2020-01-01'" 
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = -'a'
< error
| 1064 42000 …line 1 column 36 near "'a'" 
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = -0xFF
< error
| 1064 42000 …line 1 column 37 near "0xFF" 
> exec ^
: string
| ADMIN ALTER DDL JOBS 1 THREAD = _utf8mb4'x' COLLATE utf8mb4_bin
< error
| 1064 42000 …line 1 column 51 near "COLLATE utf8mb4_bin" 
> exec ^
: string
| ADMIN RELOAD EXPR_PUSHDOWN_BLACKLIST
< ok 0 0
> exec ^
: string
| ADMIN RELOAD OPT_RULE_BLACKLIST
< ok 0 0
> exec ^
: string
| ADMIN RELOAD STATISTICS
< ok 0 0
> exec ^
: string
| ADMIN RELOAD STATS_EXTENDED
< ok 0 0
> exec ^
: string
| ADMIN RELOAD BINDINGS
< ok 0 0
> exec ^
: string
| ADMIN RELOAD
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| ADMIN RELOAD foo
< error
| 1064 42000 …line 1 column 16 near "foo" 
> exec ^
: string
| ADMIN RELOAD EXPR_PUSHDOWN_BLACKLIST extra
< error
| 1064 42000 …line 1 column 42 near "extra" 
> exec ^
: string
| ADMIN FLUSH BINDINGS
< ok 0 0
> exec ^
: string
| ADMIN CAPTURE BINDINGS
< ok 0 0
> exec ^
: string
| ADMIN EVOLVE BINDINGS
< error
| 1105 HY000 Cannot enable baseline evolution feature, it is not generally available now
> exec ^
: string
| ADMIN EVOLVE BINDINGS
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN FLUSH BINDING
< error
| 1064 42000 …line 1 column 19 near "BINDING" 
> exec ^
: string
| ADMIN CAPTURE BINDING
< error
| 1064 42000 …line 1 column 21 near "BINDING" 
> exec ^
: string
| ADMIN EVOLVE BINDING
< error
| 1064 42000 …line 1 column 20 near "BINDING" 
> exec ^
: string
| ADMIN FLUSH BINDINGS extra
< error
| 1064 42000 …line 1 column 26 near "extra" 
> exec ^
: string
| ADMIN FLUSH PLAN_CACHE
< ok 0 0
> exec ^
: string
| ADMIN FLUSH SESSION PLAN_CACHE
< ok 0 0
> exec ^
: string
| ADMIN FLUSH GLOBAL PLAN_CACHE
< ok 0 0
> exec ^
: string
| ADMIN FLUSH INSTANCE PLAN_CACHE
< ok 0 0
> exec ^
: string
| ADMIN FLUSH LOCAL PLAN_CACHE
< error
| 1064 42000 …line 1 column 17 near "LOCAL PLAN_CACHE" 
> exec ^
: string
| ADMIN FLUSH SESSION
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| ADMIN FLUSH INSTANCE PLAN_CACHE extra
< error
| 1064 42000 …line 1 column 37 near "extra" 
> exec ^
: string
| ADMIN PLUGINS ENABLE audit
< ok 0 0
> exec ^
: string
| ADMIN PLUGINS ENABLE audit, whitelist
< ok 0 0
> exec ^
: string
| ADMIN PLUGINS DISABLE audit, whitelist
< ok 0 0
> exec ^
: string
| ADMIN PLUGINS ENABLE `audit`
< ok 0 0
> exec ^
: string
| ADMIN PLUGINS
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| ADMIN PLUGINS ENABLE
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| ADMIN PLUGINS TOGGLE audit
< error
| 1064 42000 …line 1 column 20 near "TOGGLE audit" 
> exec ^
: string
| ADMIN PLUGINS ENABLE select
< error
| 1064 42000 …line 1 column 27 near "select" 
> exec ^
: string
| ADMIN PLUGINS ENABLE audit whitelist
< error
| 1064 42000 …line 1 column 36 near "whitelist" 
> exec ^
: string
| ADMIN PLUGINS ENABLE audit,
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| ADMIN REPAIR TABLE t CREATE TABLE t (a int)
< error
| 8215 HY000 Failed to repair table: TiDB is not in REPAIR MODE
> exec ^
: string
| ADMIN REPAIR TABLE t CREATE TABLE t (a int)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN REPAIR TABLE t CREATE TABLE t (a char(1), b int)
< error
| 8215 HY000 Failed to repair table: TiDB is not in REPAIR MODE
> exec ^
: string
| ADMIN REPAIR TABLE t CREATE TABLE t (a char(1), b int)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN REPAIR TABLE t CREATE TABLE IF NOT EXISTS t (a int)
< error
| 8215 HY000 Failed to repair table: TiDB is not in REPAIR MODE
> exec ^
: string
| ADMIN REPAIR TABLE t CREATE TABLE IF NOT EXISTS t (a int)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN REPAIR TABLE t CREATE TEMPORARY TABLE t (a int)
< error
| 8215 HY000 Failed to repair table: TiDB is not in REPAIR MODE
> exec ^
: string
| ADMIN REPAIR TABLE t CREATE TEMPORARY TABLE t (a int)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN REPAIR TABLE t CREATE TABLE t LIKE t2
< error
| 8215 HY000 Failed to repair table: TiDB is not in REPAIR MODE
> exec ^
: string
| ADMIN REPAIR TABLE t CREATE TABLE t LIKE t2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN REPAIR TABLE t CREATE TABLE t AS SELECT 1
< error
| 8215 HY000 Failed to repair table: TiDB is not in REPAIR MODE
> exec ^
: string
| ADMIN REPAIR TABLE t CREATE TABLE t AS SELECT 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN REPAIR TABLE t CREATE TABLE t AS
< error
| 8215 HY000 Failed to repair table: TiDB is not in REPAIR MODE
> exec ^
: string
| ADMIN REPAIR TABLE t CREATE TABLE t AS
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN REPAIR TABLE t CREATE TABLE t (a int) AS
< error
| 8215 HY000 Failed to repair table: TiDB is not in REPAIR MODE
> exec ^
: string
| ADMIN REPAIR TABLE t CREATE TABLE t (a int) AS
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN REPAIR TABLE t CREATE TABLE t (a int) ENGINE=InnoDB
< error
| 8215 HY000 Failed to repair table: TiDB is not in REPAIR MODE
> exec ^
: string
| ADMIN REPAIR TABLE t CREATE TABLE t (a int) ENGINE=InnoDB
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN REPAIR TABLE t CREATE TABLE t (a int) PARTITION BY HASH(a) PARTITIONS 2
< error
| 8215 HY000 Failed to repair table: TiDB is not in REPAIR MODE
> exec ^
: string
| ADMIN REPAIR TABLE t CREATE TABLE t (a int) PARTITION BY HASH(a) PARTITIONS 2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN REPAIR TABLE t.a CREATE TABLE t (a int)
< error
| 8215 HY000 Failed to repair table: TiDB is not in REPAIR MODE
> exec ^
: string
| ADMIN REPAIR TABLE t.a CREATE TABLE t (a int)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN REPAIR TABLE `t` CREATE TABLE t (a int)
< error
| 8215 HY000 Failed to repair table: TiDB is not in REPAIR MODE
> exec ^
: string
| ADMIN REPAIR TABLE `t` CREATE TABLE t (a int)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ADMIN REPAIR TABLE t
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| ADMIN REPAIR TABLE t CREATE VIEW v AS SELECT 1
< error
| 1064 42000 …line 1 column 32 near "VIEW v AS SELECT 1" 
> exec ^
: string
| ADMIN REPAIR TABLE t CREATE OR REPLACE VIEW v AS SELECT 1
< error
| 1064 42000 …line 1 column 30 near "OR REPLACE VIEW v AS SELECT 1" 
> exec ^
: string
| ADMIN REPAIR TABLE t CREATE INDEX i ON t (a)
< error
| 1064 42000 …line 1 column 33 near "INDEX i ON t (a)" 
> exec ^
: string
| ADMIN REPAIR TABLE t CREATE TABLE t (a int) /*T! AUTO_RANDOM */
< error
| 1064 42000 …line 1 column 60 near "AUTO_RANDOM */" 
