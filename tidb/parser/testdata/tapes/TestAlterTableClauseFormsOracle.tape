-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestAlterTableClauseFormsOracle default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| ALTER TABLE t CACHE
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t CACHE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t NOCACHE
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t NOCACHE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t CACHE, NOCACHE
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t CACHE, NOCACHE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t CACHE NOCACHE
< error
| 1064 42000 …line 1 column 27 near "NOCACHE" 
> exec ^
: string
| ALTER TABLE t CACHE = 1
< error
| 1064 42000 …line 1 column 21 near "= 1" 
> exec ^
: string
| ALTER TABLE t READ ONLY
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t READ ONLY
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t READ WRITE
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t READ WRITE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t READ ONLY, READ WRITE
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t READ ONLY, READ WRITE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t READ
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| ALTER TABLE t READ ONLY READ WRITE
< error
| 1064 42000 …line 1 column 28 near "READ WRITE" 
> exec ^
: string
| ALTER TABLE t READ ONLY ENGINE = InnoDB
< error
| 1064 42000 …line 1 column 30 near "ENGINE = InnoDB" 
> exec ^
: string
| ALTER TABLE t ENGINE = InnoDB READ ONLY
< error
| 1064 42000 …line 1 column 34 near "READ ONLY" 
> exec ^
: string
| alter table t compact
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t compact
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t compact tiflash replica
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t compact tiflash replica
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t compact partition p1,p2
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t compact partition p1,p2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t compact partition p1,p2 tiflash replica
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t compact partition p1,p2 tiflash replica
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t COMPACT PARTITION
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| ALTER TABLE t COMPACT TIFLASH
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| ALTER TABLE t COMPACT REPLICA
< error
| 1064 42000 …line 1 column 29 near "REPLICA" 
> exec ^
: string
| ALTER TABLE t COMPACT, CACHE
< error
| 1064 42000 …line 1 column 22 near ", CACHE" 
> exec ^
: string
| ALTER TABLE t CACHE, COMPACT
< error
| 1064 42000 …line 1 column 28 near "COMPACT" 
> exec ^
: string
| ALTER TABLE t ATTRIBUTES='str'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ATTRIBUTES='str'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ATTRIBUTES 'str1,str2'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ATTRIBUTES 'str1,str2'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ATTRIBUTES=DeFaUlT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ATTRIBUTES=DeFaUlT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t STATS_OPTIONS="str1,str2"
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t STATS_OPTIONS="str1,str2"
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t STATS_OPTIONS=DEFAULT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t STATS_OPTIONS=DEFAULT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t STATS_OPTIONS = 'a', ATTRIBUTES = 'b'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t STATS_OPTIONS = 'a', ATTRIBUTES = 'b'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ATTRIBUTES
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| ALTER TABLE t ATTRIBUTES = 1
< error
| 1064 42000 …line 1 column 28 near "1" 
> exec ^
: string
| ALTER TABLE t ATTRIBUTES 1
< error
| 1064 42000 …line 1 column 26 near "1" 
> exec ^
: string
| ALTER TABLE t STATS_OPTIONS
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| ALTER TABLE t STATS_OPTIONS = 1
< error
| 1064 42000 …line 1 column 31 near "1" 
> exec ^
: string
| ALTER TABLE t ATTRIBUTES = 'a' ATTRIBUTES = 'b'
< error
| 1064 42000 …line 1 column 41 near "ATTRIBUTES = 'b'" 
> exec ^
: string
| ALTER TABLE t ATTRIBUTES = 'a' ENGINE = x
< error
| 1064 42000 …line 1 column 37 near "ENGINE = x" 
> exec ^
: string
| ALTER TABLE t ENGINE = x ATTRIBUTES = 'a'
< error
| 1064 42000 …line 1 column 35 near "ATTRIBUTES = 'a'" 
> exec ^
: string
| CREATE TABLE t (a INT) STATS_OPTIONS = 'str'
< error
| 1064 42000 …line 1 column 36 near "STATS_OPTIONS = 'str'" 
> exec ^
: string
| ALTER TABLE t PARTITION p
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t PARTITION p
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t PARTITION p ATTRIBUTES='str'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t PARTITION p ATTRIBUTES='str'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t PARTITION p ATTRIBUTES=DEFAULT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t PARTITION p ATTRIBUTES=DEFAULT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t PARTITION p ENGINE = x MAX_ROWS 5
< error
| 1286 42000 Unknown storage engine 'x'
> exec ^
: string
| ALTER TABLE t PARTITION p ENGINE = x MAX_ROWS 5
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t PARTITION p PLACEMENT POLICY = x
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t PARTITION p PLACEMENT POLICY = x
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t PARTITION p STORAGE ENGINE = x
< error
| 1286 42000 Unknown storage engine 'x'
> exec ^
: string
| ALTER TABLE t PARTITION p STORAGE ENGINE = x
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t CACHE PARTITION p ATTRIBUTES = 'a'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t CACHE PARTITION p ATTRIBUTES = 'a'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN a INT PARTITION p ENGINE = x
< error
| 1286 42000 Unknown storage engine 'x'
> exec ^
: string
| ALTER TABLE t ADD COLUMN a INT PARTITION p ENGINE = x
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t PARTITION
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| ALTER TABLE t PARTITION p ENGINE = x, MAX_ROWS 5
< error
| 1064 42000 …line 1 column 37 near ", MAX_ROWS 5" 
> exec ^
: string
| ALTER TABLE t PARTITION p STORAGE DISK
< error
| 1064 42000 …line 1 column 38 near "DISK" 
> exec ^
: string
| ALTER TABLE t PARTITION p ATTRIBUTES = 'a', CACHE
< error
| 1064 42000 …line 1 column 43 near ", CACHE" 
> exec ^
: string
| ALTER TABLE t CACHE, PARTITION p ATTRIBUTES = 'a'
< error
| 1064 42000 …line 1 column 30 near "PARTITION p ATTRIBUTES = 'a'" 
> exec ^
: string
| ALTER TABLE t PARTITION p ATTRIBUTES = 'a', ATTRIBUTES = 'b'
< error
| 1064 42000 …line 1 column 43 near ", ATTRIBUTES = 'b'" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN a INT, PARTITION p ENGINE = x
< error
| 1064 42000 …line 1 column 41 near "PARTITION p ENGINE = x" 
> exec ^
: string
| ALTER TABLE t_n UNION ( ), KEY_BLOCK_SIZE = 1
< error
| 8232 HY000 CREATE/ALTER table with union option is not supported
> exec ^
: string
| ALTER TABLE t_n UNION ( ), KEY_BLOCK_SIZE = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t UNION = (d.t1, d.t2)
< error
| 8232 HY000 CREATE/ALTER table with union option is not supported
> exec ^
: string
| ALTER TABLE t UNION = (d.t1, d.t2)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE d_n.t_n LOCK DEFAULT , UNION = ( t_n , d_n.t_n ) REMOVE PARTITIONING
< error
| 8232 HY000 CREATE/ALTER table with union option is not supported
> exec ^
: string
| ALTER TABLE d_n.t_n LOCK DEFAULT , UNION = ( t_n , d_n.t_n ) REMOVE PARTITIONING
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE d_n.t_n ALGORITHM = DEFAULT , MAX_ROWS 10, UNION ( d_n.t_n ) , ROW_FORMAT REDUNDANT, STATS_PERSISTENT = DEFAULT
< error
| 8232 HY000 CREATE/ALTER table with union option is not supported
> exec ^
: string
| ALTER TABLE d_n.t_n ALGORITHM = DEFAULT , MAX_ROWS 10, UNION ( d_n.t_n ) , ROW_FORMAT REDUNDANT, STATS_PERSISTENT = DEFAULT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ROW_FORMAT REDUNDANT STATS_PERSISTENT = DEFAULT KEY_BLOCK_SIZE 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ROW_FORMAT REDUNDANT STATS_PERSISTENT = DEFAULT KEY_BLOCK_SIZE 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t UNION
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| ALTER TABLE t UNION = (t1,)
< error
| 1064 42000 …line 1 column 27 near ")" 
> exec ^
: string
| ALTER TABLE t UNION = (,)
< error
| 1064 42000 …line 1 column 24 near ",)" 
> exec ^
: string
| ALTER TABLE t LOCK DEFAULT, UNION = (t1), REMOVE PARTITIONING
< error
| 1064 42000 …line 1 column 48 near "REMOVE PARTITIONING" 
> exec ^
: string
| ALTER TABLE t ENGINE = InnoDB,
< error
| 1064 42000 …line 1 column 30 near "" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN a INT,
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| ALTER TABLE t ALGORITHM = DEFAULT,
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| ALTER TABLE t ALGORITHM = DEFAULT ALGORITHM = DEFAULT
< error
| 1064 42000 …line 1 column 43 near "ALGORITHM = DEFAULT" 
> exec ^
: string
| ALTER TABLE t ALGORITHM = DEFAULT LOCK = NONE
< error
| 1064 42000 …line 1 column 38 near "LOCK = NONE" 
> exec ^
: string
| ALTER TABLE t ALGORITHM = DEFAULT, LOCK = NONE
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALGORITHM = DEFAULT, LOCK = NONE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t TTL_ENABLE = 'ON'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t TTL_ENABLE = 'ON'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t TTL_ENABLE = 'off'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t TTL_ENABLE = 'off'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t TTL_ENABLE = 'test_case'
< error
| 1064 42000 …line 1 column 38 near ""The TTL_ENABLE option has to be set 'ON' or 'OFF' 
> exec ^
: string
| alter table t /*T![ttl] TTL_ENABLE = 'test_case' */
< error
| 1064 42000 …line 1 column 51 near ""The TTL_ENABLE option has to be set 'ON' or 'OFF' 
> exec ^
: string
| ALTER TABLE t FORCE
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t FORCE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t FORCE AUTO_INCREMENT 3
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t FORCE AUTO_INCREMENT 3
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t FORCE AUTO_INCREMENT = 3
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t FORCE AUTO_INCREMENT = 3
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t force auto_random_base = 50
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t force auto_random_base = 50
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t auto_increment 30, force auto_random_base 40
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t auto_increment 30, force auto_random_base 40
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t FORCE, AUTO_INCREMENT 3
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t FORCE, AUTO_INCREMENT 3
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t FORCE ENGINE = InnoDB
< error
| 1064 42000 …line 1 column 26 near "ENGINE = InnoDB" 
> exec ^
: string
| ALTER TABLE t FORCE AUTO_ID_CACHE = 3
< error
| 1064 42000 …line 1 column 33 near "AUTO_ID_CACHE = 3" 
> exec ^
: string
| ALTER TABLE t FORCE AUTO_INCREMENT
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| alter table t partition by range FIELDS(a) (partition x values less than maxvalue)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t partition by range FIELDS(a) (partition x values less than maxvalue)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t partition by list FIELDS(a) (PARTITION p0 VALUES IN (5, 10, 15))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t partition by list FIELDS(a) (PARTITION p0 VALUES IN (5, 10, 15))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t partition by range FIELDS(a,b,c) (partition p1 values less than (1,1,1))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t partition by range FIELDS(a,b,c) (partition p1 values less than (1,1,1))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t PARTITION BY KEY FIELDS(a)
< error
| 1064 42000 …line 1 column 37 near "FIELDS(a)" 
> exec ^
: string
| alter table t partition by range(a) update indexes (a local)
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec ^
: string
| alter table t partition by range(a) update indexes (a local)
| ; SELECT FROM
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec ^
: string
| alter table t add column a int partition by hash(a) update indexes (idx_a global)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add column a int partition by hash(a) update indexes (idx_a global)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add column a int partition by hash(a) update indexes (idx_a global, idx_b local)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add column a int partition by hash(a) update indexes (idx_a global, idx_b local)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t PARTITION BY HASH(a) UPDATE INDEXES ()
< error
| 1064 42000 …line 1 column 52 near ")" 
> exec ^
: string
| ALTER TABLE t PARTITION BY HASH(a) UPDATE INDEXES (i)
< error
| 1064 42000 …line 1 column 53 near ")" 
> exec ^
: string
| ALTER TABLE t PARTITION BY HASH(a) UPDATE INDEXES (i GLOBAL,)
< error
| 1064 42000 …line 1 column 61 near ")" 
> exec ^
: string
| ALTER TABLE t LAST PARTITION LESS THAN (1000)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t LAST PARTITION LESS THAN (1000)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t LAST PARTITION LESS THAN (1000) NO_WRITE_TO_BINLOG
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t LAST PARTITION LESS THAN (1000) NO_WRITE_TO_BINLOG
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t FIRST PARTITION LESS THAN (1000)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t FIRST PARTITION LESS THAN (1000)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t FIRST PARTITION LESS THAN (1000) IF EXISTS
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t FIRST PARTITION LESS THAN (1000) IF EXISTS
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t split MAXVALUE PARTITION LESS THAN (1000)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t split MAXVALUE PARTITION LESS THAN (1000)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t merge first PARTITION LESS THAN (1000)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t merge first PARTITION LESS THAN (1000)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN a INT, LAST PARTITION LESS THAN (1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN a INT, LAST PARTITION LESS THAN (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t LAST PARTITION LESS THAN (1), ADD COLUMN a INT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t LAST PARTITION LESS THAN (1), ADD COLUMN a INT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t LAST PARTITION LESS THAN 1000
< error
| 1064 42000 …line 1 column 43 near "1000" 
> exec ^
: string
| ALTER TABLE t LAST PARTITION LESS THAN (MAXVALUE)
< error
| 1064 42000 …line 1 column 48 near "MAXVALUE)" 
> exec ^
: string
| ALTER TABLE t SPLIT PARTITION LESS THAN (1000)
< error
| 1064 42000 …line 1 column 29 near "PARTITION LESS THAN (1000)" 
> exec ^
: string
| ALTER TABLE t MERGE PARTITION LESS THAN (1000)
< error
| 1064 42000 …line 1 column 29 near "PARTITION LESS THAN (1000)" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN a INT, SPLIT MAXVALUE PARTITION LESS THAN (1)
< error
| 1064 42000 …line 1 column 37 near "SPLIT MAXVALUE PARTITION LESS THAN (1)" 
> exec ^
: string
| ALTER TABLE t SPLIT MAXVALUE PARTITION LESS THAN (1), ADD COLUMN a INT
< error
| 1064 42000 …line 1 column 53 near ", ADD COLUMN a INT" 
> exec ^
: string
| ALTER TABLE t ANALYZE PARTITION p
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ANALYZE PARTITION p
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ANALYZE PARTITION p, q
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ANALYZE PARTITION p, q
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t analyze partition a with 4 buckets
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t analyze partition a with 4 buckets
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t analyze partition a index b
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t analyze partition a index b
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t analyze partition a index b with 4 buckets
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t analyze partition a index b with 4 buckets
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ANALYZE PARTITION p INDEX
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ANALYZE PARTITION p INDEX
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ANALYZE PARTITION p INDEX i, j WITH 4 BUCKETS, 3 TOPN
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ANALYZE PARTITION p INDEX i, j WITH 4 BUCKETS, 3 TOPN
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ANALYZE PARTITION p WITH 2 CMSKETCH DEPTH
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ANALYZE PARTITION p WITH 2 CMSKETCH DEPTH
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ANALYZE PARTITION p WITH 2 CMSKETCH WIDTH
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ANALYZE PARTITION p WITH 2 CMSKETCH WIDTH
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ANALYZE PARTITION p WITH 2 SAMPLES
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ANALYZE PARTITION p WITH 2 SAMPLES
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ANALYZE PARTITION p WITH 0.5 SAMPLERATE
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ANALYZE PARTITION p WITH 0.5 SAMPLERATE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ANALYZE PARTITION p WITH 2 SAMPLERATE
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ANALYZE PARTITION p WITH 2 SAMPLERATE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ANALYZE PARTITION ALL
< error
| 1064 42000 …line 1 column 35 near "ALL" 
> exec ^
: string
| ALTER TABLE t ANALYZE PARTITION p, ADD COLUMN a INT
< error
| 1064 42000 …line 1 column 38 near "ADD COLUMN a INT" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN a INT, ANALYZE PARTITION p
< error
| 1064 42000 …line 1 column 39 near "ANALYZE PARTITION p" 
> exec ^
: string
| ALTER TABLE t ANALYZE PARTITION p WITH
< error
| 1064 42000 …line 1 column 38 near "" 
> exec ^
: string
| ALTER TABLE t ANALYZE PARTITION p WITH 4
< error
| 1064 42000 …line 1 column 40 near "" 
> exec ^
: string
| ALTER TABLE t ANALYZE PARTITION p WITH 4 BUCKETS,
< error
| 1064 42000 …line 1 column 49 near "" 
> exec ^
: string
| ALTER TABLE t ANALYZE PARTITION p WITH 2 CMSKETCH
< error
| 1064 42000 …line 1 column 49 near "" 
> exec ^
: string
| ALTER TABLE t ANALYZE PARTITION p WITH 0.5 BUCKETS
< error
| 1064 42000 …line 1 column 50 near "BUCKETS" 
> exec ^
: string
| ALTER TABLE t ANALYZE PARTITION p WITH 4 BUCKETS 3 TOPN
< error
| 1064 42000 …line 1 column 50 near "3 TOPN" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN IF NOT EXISTS (a SMALLINT UNSIGNED)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN IF NOT EXISTS (a SMALLINT UNSIGNED)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN IF NOT EXISTS a SMALLINT UNSIGNED AFTER b
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN IF NOT EXISTS a SMALLINT UNSIGNED AFTER b
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD IF NOT EXISTS a INT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD IF NOT EXISTS a INT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE employees ADD PARTITION IF NOT EXISTS ( PARTITION P2 VALUES LESS THAN MAXVALUE)
< error
| 1146 42S02 Table 'test.employees' doesn't exist
> exec ^
: string
| ALTER TABLE employees ADD PARTITION IF NOT EXISTS ( PARTITION P2 VALUES LESS THAN MAXVALUE)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE employees ADD PARTITION IF NOT EXISTS PARTITIONS 5
< error
| 1146 42S02 Table 'test.employees' doesn't exist
> exec ^
: string
| ALTER TABLE employees ADD PARTITION IF NOT EXISTS PARTITIONS 5
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD PARTITION IF NOT EXISTS LOCAL PARTITIONS 5
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD PARTITION IF NOT EXISTS LOCAL PARTITIONS 5
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t drop partition if exists p1, p2
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t drop partition if exists p1, p2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t MODIFY COLUMN IF EXISTS a varchar(255)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t MODIFY COLUMN IF EXISTS a varchar(255)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t CHANGE COLUMN IF EXISTS a b varchar(255)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t CHANGE COLUMN IF EXISTS a b varchar(255)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD INDEX IF NOT EXISTS (a) USING BTREE COMMENT 'a'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD INDEX IF NOT EXISTS (a) USING BTREE COMMENT 'a'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD KEY IF NOT EXISTS (a) USING HASH COMMENT 'a'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD KEY IF NOT EXISTS (a) USING HASH COMMENT 'a'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t DROP INDEX IF EXISTS a
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t DROP INDEX IF EXISTS a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t DROP KEY IF EXISTS a
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t DROP KEY IF EXISTS a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD CONSTRAINT fk_t2_id FOREIGN KEY IF NOT EXISTS (t2_id) REFERENCES t(id)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD COLUMN IF EXISTS a INT
< error
| 1064 42000 …line 1 column 34 near "EXISTS a INT" 
> exec ^
: string
| ALTER TABLE t MODIFY COLUMN IF NOT EXISTS a INT
< error
| 1064 42000 …line 1 column 34 near "NOT EXISTS a INT" 
> exec ^
: string
| ALTER TABLE t ADD UNIQUE IF NOT EXISTS (a)
< error
| 1064 42000 …line 1 column 27 near "IF NOT EXISTS (a)" 
> exec ^
: string
| ALTER TABLE t ADD UNIQUE KEY IF NOT EXISTS (a)
< error
| 1064 42000 …line 1 column 31 near "IF NOT EXISTS (a)" 
> exec ^
: string
| ALTER TABLE t ADD PRIMARY KEY IF NOT EXISTS (a)
< error
| 1064 42000 …line 1 column 32 near "IF NOT EXISTS (a)" 
> exec ^
: string
| ALTER TABLE t ADD FULLTEXT IF NOT EXISTS (a)
< error
| 1064 42000 …line 1 column 29 near "IF NOT EXISTS (a)" 
> exec ^
: string
| ALTER TABLE t ADD CONSTRAINT c CHECK IF NOT EXISTS (a > 1)
< error
| 1064 42000 …line 1 column 39 near "IF NOT EXISTS (a > 1)" 
> exec ^
: string
| ALTER TABLE t DROP CONSTRAINT IF EXISTS c
< error
| 1064 42000 …line 1 column 32 near "IF EXISTS c" 
> exec ^
: string
| ALTER TABLE t DROP CHECK IF EXISTS c
< error
| 1064 42000 …line 1 column 27 near "IF EXISTS c" 
> exec ^
: string
| ALTER TABLE t DROP FOREIGN KEY IF EXISTS c
< error
| 1064 42000 …line 1 column 33 near "IF EXISTS c" 
> exec ^
: string
| ALTER TABLE t DROP PRIMARY KEY IF EXISTS
< error
| 1064 42000 …line 1 column 33 near "IF EXISTS" 
> exec ^
: string
| alter table t add index (a) local
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add index (a) local
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add index (a) global
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add index (a) global
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add unique (a) local
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add unique (a) local
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add unique key (a) global
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add unique key (a) global
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add primary key (a) global
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add primary key (a) global
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD UNIQUE INDEX (a) LOCAL
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD UNIQUE INDEX (a) LOCAL
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD INDEX (a) GLOBAL LOCAL
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD INDEX (a) GLOBAL LOCAL
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD INDEX (a) COMMENT 'c' GLOBAL
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD INDEX (a) COMMENT 'c' GLOBAL
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD INDEX (a) GLOBAL COMMENT 'c'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD INDEX (a) GLOBAL COMMENT 'c'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD FULLTEXT (a) GLOBAL
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD FULLTEXT (a) GLOBAL
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD INDEX (a) NONCLUSTERED
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD INDEX (a) NONCLUSTERED
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD PRIMARY KEY (a) CLUSTERED GLOBAL
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD PRIMARY KEY (a) CLUSTERED GLOBAL
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD PRIMARY KEY (a) GLOBAL CLUSTERED
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD PRIMARY KEY (a) GLOBAL CLUSTERED
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a int, INDEX (a) GLOBAL)
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) GLOBAL
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i ON t (a) GLOBAL
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD a INT PRIMARY KEY GLOBAL
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD a INT PRIMARY KEY GLOBAL
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD a INT PRIMARY KEY CLUSTERED GLOBAL
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD a INT PRIMARY KEY CLUSTERED GLOBAL
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD a INT UNIQUE GLOBAL
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD a INT UNIQUE GLOBAL
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD a INT UNIQUE KEY LOCAL
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD a INT UNIQUE KEY LOCAL
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD a INT KEY GLOBAL
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD a INT KEY GLOBAL
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD a INT PRIMARY KEY GLOBAL CLUSTERED
< error
| 1064 42000 …line 1 column 52 near "CLUSTERED" 
> exec ^
: string
| ALTER TABLE t ADD a INT UNIQUE GLOBAL LOCAL
< error
| 1064 42000 …line 1 column 43 near "LOCAL" 
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX ((VEC_COSINE_DISTANCE(a, b))) USING HNSW COMMENT 'a'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX ((VEC_COSINE_DISTANCE(a, b))) USING HNSW COMMENT 'a'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX ((lower(a))) USING HNSW COMMENT 'a'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX ((lower(a))) USING HNSW COMMENT 'a'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX (a, (VEC_COSINE_DISTANCE(a))) USING HNSW COMMENT 'a'
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX (a, (VEC_COSINE_DISTANCE(a))) USING HNSW COMMENT 'a'
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX ((VEC_COSINE_DISTANCE(a))) USING HYPO COMMENT 'a'
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX ((VEC_COSINE_DISTANCE(a))) USING HYPO COMMENT 'a'
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX ((VEC_COSINE_DISTANCE(a))) COMMENT 'a'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX ((VEC_COSINE_DISTANCE(a))) COMMENT 'a'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX IF NOT EXISTS ((VEC_COSINE_DISTANCE(a))) USING HNSW COMMENT 'a'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX IF NOT EXISTS ((VEC_COSINE_DISTANCE(a))) USING HNSW COMMENT 'a'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX idx USING HNSW ((lower(a)))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX idx USING HNSW ((lower(a)))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX ((lower(a))) GLOBAL
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX ((lower(a))) GLOBAL
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD VECTOR KEY ((lower(a)))
< error
| 1064 42000 …line 1 column 28 near "KEY ((lower(a)))" 
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX ()
< error
| 1064 42000 …line 1 column 33 near ")" 
> exec ^
: string
| ALTER TABLE t ADD CONSTRAINT c VECTOR INDEX ((lower(a)))
< error
| 1064 42000 …line 1 column 37 near "VECTOR INDEX ((lower(a)))" 
> exec ^
: string
| CREATE INDEX i ON t (a) USING HYPO
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i ON t (a) USING HYPO
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX i ON t (a) USING `HYPO`
< error
| 1064 42000 …line 1 column 36 near "`HYPO`" 
> exec ^
: string
| ALTER TABLE db.t RENAME to db1.t1
< ok 0 0
> exec ^
: string
| ALTER TABLE db.t RENAME db1.t1
< ok 0 0
> exec ^
: string
| ALTER TABLE db.t RENAME = db1.t1
< ok 0 0
> exec ^
: string
| ALTER TABLE db.t RENAME as db1.t1
< ok 0 0
> exec ^
: string
| ALTER TABLE t RENAME = t1
< ok 0 0
> exec ^
: string
| ALTER TABLE t RENAME COLUMN a b
< error
| 1064 42000 …line 1 column 31 near "b" 
> exec ^
: string
| ALTER TABLE t RENAME INDEX a b
< error
| 1064 42000 …line 1 column 30 near "b" 
> exec ^
: string
| ALTER TABLE t CONVERT TO CHARSET UTF8
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t CONVERT TO CHARSET UTF8
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t CONVERT TO CHARSET UTF8 COLLATE UTF8_BIN
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t CONVERT TO CHARSET UTF8 COLLATE UTF8_BIN
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table d_n.t_n convert to charset default
< error
| 1146 42S02 Table 'd_n.t_n' doesn't exist
> exec ^
: string
| alter table d_n.t_n convert to charset default
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table d_n.t_n convert to char set default
< error
| 1146 42S02 Table 'd_n.t_n' doesn't exist
> exec ^
: string
| alter table d_n.t_n convert to char set default
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t CONVERT TO CHARACTER SET DEFAULT COLLATE utf8_bin
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t CONVERT TO CHARACTER SET DEFAULT COLLATE utf8_bin
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t CONVERT TO utf8
< error
| 1064 42000 …line 1 column 29 near "utf8" 
> exec ^
: string
| ALTER TABLE t CONVERT TO CHARSET
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| ALTER TABLE t CONVERT CHARSET utf8
< error
| 1064 42000 …line 1 column 29 near "CHARSET utf8" 
> exec ^
: string
| ALTER TABLE t CONVERT TO CHARSET utf8 COLLATE = utf8_bin
< error
| 1064 42000 …line 1 column 47 near "= utf8_bin" 
> exec ^
: string
| ALTER TABLE t ADD ( CHECK ( true ) )
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD ( CHECK ( true ) )
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD ( CONSTRAINT CHECK ( true ) )
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD ( CONSTRAINT CHECK ( true ) )
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add column (b int, constraint c unique key (b))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add column (b int, constraint c unique key (b))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN ( CONSTRAINT CHECK ( true ) ENFORCED , CHECK ( true ) )
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN ( CONSTRAINT CHECK ( true ) ENFORCED , CHECK ( true ) )
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN (a1 int, a2 int, CONSTRAINT b1 CHECK (a1>0), CONSTRAINT b2 CHECK (a2<10))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN (a1 int, a2 int, CONSTRAINT b1 CHECK (a1>0), CONSTRAINT b2 CHECK (a2<10))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE `t` ADD COLUMN (`a1` INT, PRIMARY KEY (`a1`))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE `t` ADD COLUMN (`a1` INT, PRIMARY KEY (`a1`))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD (a1 int, a2 int, PRIMARY KEY (a1), UNIQUE (a2))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD (a1 int, a2 int, PRIMARY KEY (a1), UNIQUE (a2))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN (INDEX (a))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN (INDEX (a))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN (FOREIGN KEY (a) REFERENCES u(b))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN (FOREIGN KEY (a) REFERENCES u(b))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN ()
< error
| 1064 42000 …line 1 column 27 near ")" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN (a INT,)
< error
| 1064 42000 …line 1 column 33 near ")" 
> exec ^
: string
| ALTER TABLE t ADD a1 int CONSTRAINT ident CHECK ( a1 > 1 ) REFERENCES b ON DELETE CASCADE ON UPDATE CASCADE
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD a1 int CONSTRAINT ident CHECK ( a1 > 1 ) REFERENCES b ON DELETE CASCADE ON UPDATE CASCADE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN a DATE CONSTRAINT CHECK ( a > 0 ) FIRST
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN a DATE CONSTRAINT CHECK ( a > 0 ) FIRST
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD a2 int CONSTRAINT ident CHECK (a2 > 1) ENFORCED
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD a2 int CONSTRAINT ident CHECK (a2 > 1) ENFORCED
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD a2 int CONSTRAINT ident CHECK (a2 > 1) NOT ENFORCED
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD a2 int CONSTRAINT ident CHECK (a2 > 1) NOT ENFORCED
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a INT CONSTRAINT c CHECK (a > 1))
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD a INT CONSTRAINT c NOT NULL
< error
| 1064 42000 …line 1 column 40 near "NOT NULL" 
> exec ^
: string
| ALTER TABLE t ADD a INT CONSTRAINT c
< error
| 1064 42000 …line 1 column 36 near "" 
> exec ^
: string
| ALTER TABLE t ADD a CONSTRAINT c CHECK (a>1)
< error
| 1064 42000 …line 1 column 30 near "CONSTRAINT c CHECK (a>1)" 
> exec ^
: string
| alter table d_n.t_n add constraint foreign key ident (ident(1)) references d_n.t_n match full on delete set null
< ok 0 0
> exec ^
: string
| alter table t_n add constraint ident foreign key (ident,ident(1)) references t_n match full on update set null on delete restrict
< ok 0 0
> exec ^
: string
| alter table d_n.t_n add foreign key ident (ident, ident(1) asc) references t_n match partial on delete cascade remove partitioning
< ok 0 0
> exec ^
: string
| alter table d_n.t_n add constraint foreign key (ident asc) references d_n.t_n match simple on update cascade on delete cascade
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD FOREIGN KEY (a) REFERENCES u(b) MATCH foo
< error
| 1064 42000 …line 1 column 59 near "foo" 
> exec ^
: string
| ALTER TABLE t ADD FOREIGN KEY ((a + 1)) REFERENCES u(b)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD CONSTRAINT c FOREIGN KEY ((a+1)) REFERENCES u(b)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, FOREIGN KEY ((a + 1)) REFERENCES u(b))
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX IF NOT EXISTS idx ((lower(a)))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX IF NOT EXISTS idx ((lower(a)))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE d_n.t_n WITHOUT VALIDATION , ADD PARTITION ( PARTITION ident VALUES LESS THAN ( MAXVALUE ) STORAGE ENGINE text_string MAX_ROWS 12 )
< error
| 1146 42S02 Table 'd_n.t_n' doesn't exist
> exec ^
: string
| ALTER TABLE d_n.t_n WITHOUT VALIDATION , ADD PARTITION ( PARTITION ident VALUES LESS THAN ( MAXVALUE ) STORAGE ENGINE text_string MAX_ROWS 12 )
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE d_n.t_n WITH VALIDATION , ADD PARTITION NO_WRITE_TO_BINLOG (PARTITION ident VALUES LESS THAN MAXVALUE STORAGE ENGINE = text_string, PARTITION ident VALUES LESS THAN ( MAXVALUE ) (SUBPARTITION text_string MIN_ROWS 11))
< error
| 1146 42S02 Table 'd_n.t_n' doesn't exist
> exec ^
: string
| ALTER TABLE d_n.t_n WITH VALIDATION , ADD PARTITION NO_WRITE_TO_BINLOG (PARTITION ident VALUES LESS THAN MAXVALUE STORAGE ENGINE = text_string, PARTITION ident VALUES LESS THAN ( MAXVALUE ) (SUBPARTITION text_string MIN_ROWS 11))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE d_n.t_n WITHOUT VALIDATION , ADD PARTITION ( PARTITION ident VALUES IN ( DEFAULT ) STORAGE ENGINE text_string MAX_ROWS 12 )
< error
| 1146 42S02 Table 'd_n.t_n' doesn't exist
> exec ^
: string
| ALTER TABLE d_n.t_n WITHOUT VALIDATION , ADD PARTITION ( PARTITION ident VALUES IN ( DEFAULT ) STORAGE ENGINE text_string MAX_ROWS 12 )
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES IN (DEFAULT, 1))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES IN (DEFAULT, 1))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES IN (1, DEFAULT))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES IN (1, DEFAULT))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p DEFAULT)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p DEFAULT)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (1) PLACEMENT POLICY x)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (1) PLACEMENT POLICY x)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (1) STORAGE DISK)
< error
| 1064 42000 …line 1 column 74 near "DISK)" 
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (1) TTL_ENABLE = 'ON')
< error
| 1064 42000 …line 1 column 72 near "TTL_ENABLE = 'ON')" 
> exec ^
: string
| alter table m partition t placement policy="ww"
< error
| 1146 42S02 Table 'test.m' doesn't exist
> exec ^
: string
| alter table m partition t placement policy="ww"
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table m partition t /*T![placement] placement policy="ww" */
< error
| 1146 42S02 Table 'test.m' doesn't exist
> exec ^
: string
| alter table m partition t /*T![placement] placement policy="ww" */
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 'x'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 'x'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT NULL
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT NULL
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT TRUE
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT TRUE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT (CURRENT_TIMESTAMP)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT (CURRENT_TIMESTAMP)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT (1+1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT (1+1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT CURRENT_TIMESTAMP
< error
| 1064 42000 …line 1 column 58 near "CURRENT_TIMESTAMP" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT NOW()
< error
| 1064 42000 …line 1 column 44 near "NOW()" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 1+1
< error
| 1064 42000 …line 1 column 43 near "+1" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -'x'
< error
| 1064 42000 …line 1 column 45 near "'x'" 
> exec ^
: string
| ALTER TABLE t ALTER CONSTRAINT ident ENFORCED
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER CONSTRAINT ident ENFORCED
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER CHECK ident NOT ENFORCED
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER CHECK ident NOT ENFORCED
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t_n ALTER CONSTRAINT ident
< error
| 1064 42000 …line 1 column 38 near "" 
> exec ^
: string
| ALTER TABLE t ALTER CHECK ident
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| ALTER TABLE attributes ADD COLUMN a INT
< error
| 1146 42S02 Table 'test.attributes' doesn't exist
> exec ^
: string
| ALTER TABLE attributes ADD COLUMN a INT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN attributes INT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN attributes INT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN stats_options INT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN stats_options INT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN vector INT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN vector INT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN compact INT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN compact INT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t DROP COLUMN cache
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t DROP COLUMN cache
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t DROP COLUMN nocache
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t DROP COLUMN nocache
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t DROP COLUMN last
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t DROP COLUMN last
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t DROP COLUMN global
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t DROP COLUMN global
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t RENAME COLUMN a TO attributes
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t RENAME COLUMN a TO attributes
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t RENAME TO attributes
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD INDEX i (attributes)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD INDEX i (attributes)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD INDEX attributes (a)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD INDEX attributes (a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD INDEX cmsketch (a)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD INDEX cmsketch (a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD INDEX i (global, local)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD INDEX i (global, local)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t MODIFY COLUMN attributes INT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t MODIFY COLUMN attributes INT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t CHANGE COLUMN a attributes INT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t CHANGE COLUMN a attributes INT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE attributes (a INT)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (attributes INT, stats_options INT, topn INT, samplerate INT, cmsketch INT, depth INT, width INT, samples INT)
< ok 0 0
> exec ^
: string
| CREATE INDEX attributes ON t (a)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX attributes ON t (a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t AS attributes
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t AS attributes
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT attributes.a FROM t attributes
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT attributes.a FROM t attributes
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t attributes WHERE attributes.id = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t attributes WHERE attributes.id = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| UPDATE t attributes SET attributes.a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| UPDATE t attributes SET attributes.a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ANALYZE PARTITION p WITH 4 buckets
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ANALYZE PARTITION p WITH 4 buckets
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ANALYZE PARTITION p WITH 4 topn
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ANALYZE PARTITION p WITH 4 topn
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD INDEX i (first, last, merge, split)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD INDEX i (first, last, merge, split)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN merge INT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN merge INT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN split INT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN split INT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN global INT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN global INT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN local INT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN local INT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN clustered INT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN clustered INT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN nonclustered INT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN nonclustered INT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD INDEX global (a)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD INDEX global (a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t DROP INDEX local
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t DROP INDEX local
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t MODIFY COLUMN compact INT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t MODIFY COLUMN compact INT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t CHANGE COLUMN a vector INT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t CHANGE COLUMN a vector INT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t RENAME COLUMN a TO last
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t RENAME COLUMN a TO last
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE merge ADD COLUMN a INT
< error
| 1146 42S02 Table 'test.merge' doesn't exist
> exec ^
: string
| ALTER TABLE merge ADD COLUMN a INT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION last VALUES LESS THAN (1))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION last VALUES LESS THAN (1))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t DROP PARTITION first, last
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t DROP PARTITION first, last
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t TRUNCATE PARTITION cache
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t TRUNCATE PARTITION cache
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ANALYZE PARTITION buckets
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ANALYZE PARTITION buckets
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ANALYZE PARTITION p INDEX topn
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ANALYZE PARTITION p INDEX topn
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t COMPACT PARTITION tiflash
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t COMPACT PARTITION tiflash
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t PARTITION attributes ENGINE = x
< error
| 1286 42000 Unknown storage engine 'x'
> exec ^
: string
| ALTER TABLE t PARTITION attributes ENGINE = x
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t PARTITION p PLACEMENT POLICY = attributes
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t PARTITION p PLACEMENT POLICY = attributes
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN global SET DEFAULT 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER COLUMN global SET DEFAULT 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER INDEX global VISIBLE
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER INDEX global VISIBLE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD CONSTRAINT global CHECK (a > 1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD CONSTRAINT global CHECK (a > 1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t DROP CONSTRAINT vector
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t DROP CONSTRAINT vector
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT first, last, merge, split, global, local FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT first, last, merge, split, global, local FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t `ATTRIBUTES` = 'a'
< error
| 1064 42000 …line 1 column 26 near "`ATTRIBUTES` = 'a'" 
> exec ^
: string
| ALTER TABLE t `STATS_OPTIONS` = 'a'
< error
| 1064 42000 …line 1 column 29 near "`STATS_OPTIONS` = 'a'" 
> exec ^
: string
| ALTER TABLE t PARTITION p GARBAGE
< error
| 1064 42000 …line 1 column 33 near "GARBAGE" 
> exec ^
: string
| ALTER TABLE t COMPACT GARBAGE
< error
| 1064 42000 …line 1 column 29 near "GARBAGE" 
> exec ^
: string
| ALTER TABLE t ANALYZE PARTITION p GARBAGE
< error
| 1064 42000 …line 1 column 41 near "GARBAGE" 
> exec ^
: string
| ALTER TABLE t CACHE GARBAGE
< error
| 1064 42000 …line 1 column 27 near "GARBAGE" 
> exec ^
: string
| ALTER TABLE t READ ONLY GARBAGE
< error
| 1064 42000 …line 1 column 31 near "GARBAGE" 
> exec ^
: string
| ALTER TABLE t ATTRIBUTES = 'a' GARBAGE
< error
| 1064 42000 …line 1 column 38 near "GARBAGE" 
> exec ^
: string
| ALTER TABLE t SPLIT MAXVALUE PARTITION LESS THAN (1) GARBAGE
< error
| 1064 42000 …line 1 column 60 near "GARBAGE" 
> exec ^
: string
| ALTER TABLE t LAST PARTITION LESS THAN (1) GARBAGE
< error
| 1064 42000 …line 1 column 50 near "GARBAGE" 
> exec ^
: string
| ALTER TABLE t ADD VECTOR INDEX ((lower(a))) GARBAGE
< error
| 1064 42000 …line 1 column 51 near "GARBAGE" 
> exec ^
: string
| ALTER TABLE t PARTITION BY HASH(a) UPDATE INDEXES (i GLOBAL) GARBAGE
< error
| 1064 42000 …line 1 column 68 near "GARBAGE" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN (a INT, PRIMARY KEY (a)) GARBAGE
< error
| 1064 42000 …line 1 column 57 near "GARBAGE" 
> exec ^
: string
| ALTER TABLE t RENAME = d.t1 GARBAGE
< error
| 1064 42000 …line 1 column 35 near "GARBAGE" 
> exec ^
: string
| ALTER TABLE t CONVERT TO CHARSET utf8 GARBAGE
< error
| 1064 42000 …line 1 column 45 near "GARBAGE" 
> exec ^
: string
| ALTER TABLE t UNION = () GARBAGE
< error
| 1064 42000 …line 1 column 32 near "GARBAGE" 
