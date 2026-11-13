-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestPartitionMaintenanceModifiersMatchTiDB default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| ALTER TABLE pt OPTIMIZE PARTITION NO_WRITE_TO_BINLOG p0
< error
| 1146 42S02 Table 'test.pt' doesn't exist
> exec ^
: string
| ALTER TABLE pt OPTIMIZE PARTITION NO_WRITE_TO_BINLOG p0
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE pt OPTIMIZE PARTITION LOCAL p0, p1
< error
| 1146 42S02 Table 'test.pt' doesn't exist
> exec ^
: string
| ALTER TABLE pt OPTIMIZE PARTITION LOCAL p0, p1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE pt OPTIMIZE PARTITION LOCAL ALL
< error
| 1146 42S02 Table 'test.pt' doesn't exist
> exec ^
: string
| ALTER TABLE pt OPTIMIZE PARTITION LOCAL ALL
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE pt REPAIR PARTITION NO_WRITE_TO_BINLOG p0
< error
| 1146 42S02 Table 'test.pt' doesn't exist
> exec ^
: string
| ALTER TABLE pt REPAIR PARTITION NO_WRITE_TO_BINLOG p0
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE pt REPAIR PARTITION LOCAL ALL
< error
| 1146 42S02 Table 'test.pt' doesn't exist
> exec ^
: string
| ALTER TABLE pt REPAIR PARTITION LOCAL ALL
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE pt REBUILD PARTITION LOCAL p0
< error
| 1146 42S02 Table 'test.pt' doesn't exist
> exec ^
: string
| ALTER TABLE pt REBUILD PARTITION LOCAL p0
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE pt REBUILD PARTITION NO_WRITE_TO_BINLOG ALL
< error
| 1146 42S02 Table 'test.pt' doesn't exist
> exec ^
: string
| ALTER TABLE pt REBUILD PARTITION NO_WRITE_TO_BINLOG ALL
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE pt COALESCE PARTITION LOCAL 3
< error
| 1146 42S02 Table 'test.pt' doesn't exist
> exec ^
: string
| ALTER TABLE pt COALESCE PARTITION LOCAL 3
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE pt COALESCE PARTITION NO_WRITE_TO_BINLOG 3
< error
| 1146 42S02 Table 'test.pt' doesn't exist
> exec ^
: string
| ALTER TABLE pt COALESCE PARTITION NO_WRITE_TO_BINLOG 3
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE pt ADD PARTITION LOCAL
< error
| 1146 42S02 Table 'test.pt' doesn't exist
> exec ^
: string
| ALTER TABLE pt ADD PARTITION LOCAL
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE pt ADD PARTITION NO_WRITE_TO_BINLOG PARTITIONS 10
< error
| 1146 42S02 Table 'test.pt' doesn't exist
> exec ^
: string
| ALTER TABLE pt ADD PARTITION NO_WRITE_TO_BINLOG PARTITIONS 10
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE pt ADD PARTITION LOCAL (PARTITION p9 VALUES LESS THAN (100))
< error
| 1146 42S02 Table 'test.pt' doesn't exist
> exec ^
: string
| ALTER TABLE pt ADD PARTITION LOCAL (PARTITION p9 VALUES LESS THAN (100))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE pt REORGANIZE PARTITION LOCAL
< error
| 1146 42S02 Table 'test.pt' doesn't exist
> exec ^
: string
| ALTER TABLE pt REORGANIZE PARTITION LOCAL
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE pt REORGANIZE PARTITION NO_WRITE_TO_BINLOG p0 INTO (PARTITION s0 VALUES LESS THAN (5))
< error
| 1146 42S02 Table 'test.pt' doesn't exist
> exec ^
: string
| ALTER TABLE pt REORGANIZE PARTITION NO_WRITE_TO_BINLOG p0 INTO (PARTITION s0 VALUES LESS THAN (5))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE pt ADD PARTITION
< error
| 1146 42S02 Table 'test.pt' doesn't exist
> exec ^
: string
| ALTER TABLE pt ADD PARTITION
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE pt REORGANIZE PARTITION
< error
| 1146 42S02 Table 'test.pt' doesn't exist
> exec ^
: string
| ALTER TABLE pt REORGANIZE PARTITION
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE pt OPTIMIZE PARTITION LOCAL local
< error
| 1146 42S02 Table 'test.pt' doesn't exist
> exec ^
: string
| ALTER TABLE pt OPTIMIZE PARTITION LOCAL local
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE pt OPTIMIZE PARTITION LOCAL `local`
< error
| 1146 42S02 Table 'test.pt' doesn't exist
> exec ^
: string
| ALTER TABLE pt OPTIMIZE PARTITION LOCAL `local`
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE pt REBUILD PARTITION LOCAL local, local
< error
| 1146 42S02 Table 'test.pt' doesn't exist
> exec ^
: string
| ALTER TABLE pt REBUILD PARTITION LOCAL local, local
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE pt REORGANIZE PARTITION LOCAL remove INTO (PARTITION s0 VALUES LESS THAN (5))
< error
| 1146 42S02 Table 'test.pt' doesn't exist
> exec ^
: string
| ALTER TABLE pt REORGANIZE PARTITION LOCAL remove INTO (PARTITION s0 VALUES LESS THAN (5))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t REMOVE PARTITIONING
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t REMOVE PARTITIONING
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN a INT REMOVE PARTITIONING
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN a INT REMOVE PARTITIONING
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN a INT, ADD INDEX (n) REMOVE PARTITIONING
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN a INT, ADD INDEX (n) REMOVE PARTITIONING
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t LOCK = DEFAULT REMOVE PARTITIONING
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t LOCK = DEFAULT REMOVE PARTITIONING
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t MODIFY COLUMN n INT AFTER id REMOVE PARTITIONING
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t MODIFY COLUMN n INT AFTER id REMOVE PARTITIONING
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t CHANGE COLUMN n m INT AFTER id REMOVE PARTITIONING
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t CHANGE COLUMN n m INT AFTER id REMOVE PARTITIONING
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t PARTITION BY HASH(id)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t PARTITION BY HASH(id)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN a INT PARTITION BY HASH(id)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN a INT PARTITION BY HASH(id)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t COMMENT 'c' PARTITION BY HASH(id)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t COMMENT 'c' PARTITION BY HASH(id)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE pt OPTIMIZE LOCAL PARTITION p0
< error
| 1064 42000 …line 1 column 29 near "LOCAL PARTITION p0" 
> exec ^
: string
| ALTER TABLE pt REBUILD LOCAL PARTITION p0
< error
| 1064 42000 …line 1 column 28 near "LOCAL PARTITION p0" 
> exec ^
: string
| ALTER TABLE pt OPTIMIZE PARTITION p0 LOCAL
< error
| 1064 42000 …line 1 column 42 near "LOCAL" 
> exec ^
: string
| ALTER TABLE pt OPTIMIZE PARTITION ALL NO_WRITE_TO_BINLOG
< error
| 1064 42000 …line 1 column 56 near "NO_WRITE_TO_BINLOG" 
> exec ^
: string
| ALTER TABLE pt REPAIR PARTITION p0 LOCAL
< error
| 1064 42000 …line 1 column 40 near "LOCAL" 
> exec ^
: string
| ALTER TABLE pt REBUILD PARTITION ALL NO_WRITE_TO_BINLOG
< error
| 1064 42000 …line 1 column 55 near "NO_WRITE_TO_BINLOG" 
> exec ^
: string
| ALTER TABLE pt COALESCE PARTITION 3 LOCAL
< error
| 1064 42000 …line 1 column 41 near "LOCAL" 
> exec ^
: string
| ALTER TABLE pt ADD PARTITION PARTITIONS 10 LOCAL
< error
| 1064 42000 …line 1 column 48 near "LOCAL" 
> exec ^
: string
| ALTER TABLE pt ADD PARTITION (PARTITION p9 VALUES LESS THAN (100)) LOCAL
< error
| 1064 42000 …line 1 column 72 near "LOCAL" 
> exec ^
: string
| ALTER TABLE pt REORGANIZE PARTITION p0 LOCAL INTO (PARTITION s0 VALUES LESS THAN (5))
< error
| 1064 42000 …line 1 column 44 near "LOCAL INTO (PARTITION s0 VALUES LESS THAN (5))" 
> exec ^
: string
| ALTER TABLE pt OPTIMIZE PARTITION LOCAL NO_WRITE_TO_BINLOG p0
< error
| 1064 42000 …line 1 column 58 near "NO_WRITE_TO_BINLOG p0" 
> exec ^
: string
| ALTER TABLE pt OPTIMIZE PARTITION NO_WRITE_TO_BINLOG LOCAL p0
< error
| 1064 42000 …line 1 column 61 near "p0" 
> exec ^
: string
| ALTER TABLE pt OPTIMIZE PARTITION LOCAL LOCAL p0
< error
| 1064 42000 …line 1 column 48 near "p0" 
> exec ^
: string
| ALTER TABLE pt ADD PARTITION LOCAL NO_WRITE_TO_BINLOG PARTITIONS 10
< error
| 1064 42000 …line 1 column 53 near "NO_WRITE_TO_BINLOG PARTITIONS 10" 
> exec ^
: string
| ALTER TABLE pt COALESCE PARTITION LOCAL 3 LOCAL
< error
| 1064 42000 …line 1 column 47 near "LOCAL" 
> exec ^
: string
| ALTER TABLE pt OPTIMIZE PARTITION LOCAL
< error
| 1064 42000 …line 1 column 39 near "" 
> exec ^
: string
| ALTER TABLE pt OPTIMIZE PARTITION NO_WRITE_TO_BINLOG
< error
| 1064 42000 …line 1 column 52 near "" 
> exec ^
: string
| ALTER TABLE pt OPTIMIZE PARTITION local
< error
| 1064 42000 …line 1 column 39 near "" 
> exec ^
: string
| ALTER TABLE pt OPTIMIZE PARTITION no_write_to_binlog
< error
| 1064 42000 …line 1 column 52 near "" 
> exec ^
: string
| ALTER TABLE pt REPAIR PARTITION LOCAL
< error
| 1064 42000 …line 1 column 37 near "" 
> exec ^
: string
| ALTER TABLE pt REBUILD PARTITION LOCAL
< error
| 1064 42000 …line 1 column 38 near "" 
> exec ^
: string
| ALTER TABLE pt REORGANIZE PARTITION LOCAL INTO (PARTITION s0 VALUES LESS THAN (5))
< error
| 1064 42000 …line 1 column 46 near "INTO (PARTITION s0 VALUES LESS THAN (5))" 
> exec ^
: string
| ALTER TABLE pt COALESCE PARTITION
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| ALTER TABLE pt COALESCE PARTITION LOCAL
< error
| 1064 42000 …line 1 column 39 near "" 
> exec ^
: string
| ALTER TABLE pt ADD PARTITION LOCAL PARTITIONS
< error
| 1064 42000 …line 1 column 45 near "" 
> exec ^
: string
| ALTER TABLE pt ADD PARTITION LOCAL ALL
< error
| 1064 42000 …line 1 column 38 near "ALL" 
> exec ^
: string
| ALTER TABLE pt CHECK PARTITION LOCAL p0
< error
| 1064 42000 …line 1 column 39 near "p0" 
> exec ^
: string
| ALTER TABLE pt CHECK PARTITION NO_WRITE_TO_BINLOG p0
< error
| 1064 42000 …line 1 column 49 near "NO_WRITE_TO_BINLOG p0" 
> exec ^
: string
| ALTER TABLE pt TRUNCATE PARTITION LOCAL p0
< error
| 1064 42000 …line 1 column 42 near "p0" 
> exec ^
: string
| ALTER TABLE pt TRUNCATE PARTITION NO_WRITE_TO_BINLOG p0
< error
| 1064 42000 …line 1 column 52 near "NO_WRITE_TO_BINLOG p0" 
> exec ^
: string
| ALTER TABLE pt DROP PARTITION LOCAL p0
< error
| 1064 42000 …line 1 column 38 near "p0" 
> exec ^
: string
| ALTER TABLE pt ANALYZE PARTITION LOCAL p0
< error
| 1064 42000 …line 1 column 41 near "p0" 
> exec ^
: string
| ALTER TABLE pt ANALYZE PARTITION NO_WRITE_TO_BINLOG p0
< error
| 1064 42000 …line 1 column 51 near "NO_WRITE_TO_BINLOG p0" 
> exec ^
: string
| ALTER TABLE pt ANALYZE PARTITION p0 LOCAL
< error
| 1064 42000 …line 1 column 41 near "LOCAL" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN a INT, REMOVE PARTITIONING
< error
| 1064 42000 …line 1 column 38 near "REMOVE PARTITIONING" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN a INT, ADD INDEX (n), REMOVE PARTITIONING
< error
| 1064 42000 …line 1 column 53 near "REMOVE PARTITIONING" 
> exec ^
: string
| ALTER TABLE t REMOVE PARTITIONING, ADD COLUMN a INT
< error
| 1064 42000 …line 1 column 34 near ", ADD COLUMN a INT" 
> exec ^
: string
| ALTER TABLE t REMOVE PARTITIONING ADD COLUMN a INT
< error
| 1064 42000 …line 1 column 37 near "ADD COLUMN a INT" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN a INT REMOVE PARTITIONING REMOVE PARTITIONING
< error
| 1064 42000 …line 1 column 57 near "REMOVE PARTITIONING" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN a INT, PARTITION BY HASH(id)
< error
| 1064 42000 …line 1 column 41 near "PARTITION BY HASH(id)" 
> exec ^
: string
| ALTER TABLE t PARTITION BY HASH(id) ADD COLUMN a INT
< error
| 1064 42000 …line 1 column 39 near "ADD COLUMN a INT" 
> exec ^
: string
| ALTER TABLE t PARTITION BY HASH(id), ADD COLUMN a INT
< error
| 1064 42000 …line 1 column 36 near ", ADD COLUMN a INT" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN a INT REMOVE PARTITIONING PARTITION BY HASH(id)
< error
| 1064 42000 …line 1 column 60 near "PARTITION BY HASH(id)" 
> exec ^
: string
| ALTER TABLE t REMOVE PARTITIONING PARTITION BY HASH(id)
< error
| 1064 42000 …line 1 column 43 near "PARTITION BY HASH(id)" 
> exec ^
: string
| ALTER TABLE t PARTITION BY HASH(id) PARTITIONS 4, ALGORITHM=INPLACE
< error
| 1064 42000 …line 1 column 49 near ", ALGORITHM=INPLACE" 
> exec ^
: string
| ALTER TABLE t PARTITION BY HASH(id) PARTITIONS 4, LOCK=NONE
< error
| 1064 42000 …line 1 column 49 near ", LOCK=NONE" 
> exec ^
: string
| ALTER TABLE t PARTITION BY HASH(id) ALGORITHM=INPLACE
< error
| 1064 42000 …line 1 column 45 near "ALGORITHM=INPLACE" 
> exec ^
: string
| ALTER TABLE t PARTITION BY HASH(id), ENABLE KEYS
< error
| 1064 42000 …line 1 column 36 near ", ENABLE KEYS" 
> exec ^
: string
| ALTER TABLE t REMOVE PARTITIONING, ALGORITHM=INPLACE
< error
| 1064 42000 …line 1 column 34 near ", ALGORITHM=INPLACE" 
> exec ^
: string
| ALTER TABLE t REMOVE
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| ALTER TABLE t REMOVE PARTITION
< error
| 1064 42000 …line 1 column 30 near "PARTITION" 
