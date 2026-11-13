-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOptionEqualFormsMatchEngine default#1
> exec
| PREPARE omni_sweep_stmt FROM ?
: string
| create table t (a int) engine := innodb
< ok 0 0
> exec ^
: string
| create table t (a int) secondary_engine := null
< ok 0 0
> exec ^
: string
| create table t (a int) secondary_engine := x
< ok 0 0
> exec ^
: string
| create table t (a int) max_rows := 1
< ok 0 0
> exec ^
: string
| create table t (a int) min_rows := 1
< ok 0 0
> exec ^
: string
| create table t (a int) avg_row_length := 1
< ok 0 0
> exec ^
: string
| create table t (a int) password := 'x'
< ok 0 0
> exec ^
: string
| create table t (a int) comment := 'x'
< ok 0 0
> exec ^
: string
| create table t (a int) compression := 'zlib'
< ok 0 0
> exec ^
: string
| create table t (a int) encryption := 'y'
< ok 0 0
> exec ^
: string
| create table t (a int) auto_increment := 1
< ok 0 0
> exec ^
: string
| create table t (a int) pack_keys := 1
< ok 0 0
> exec ^
: string
| create table t (a int) pack_keys := default
< ok 0 0
> exec ^
: string
| create table t (a int) stats_auto_recalc := 1
< ok 0 0
> exec ^
: string
| create table t (a int) stats_auto_recalc := default
< ok 0 0
> exec ^
: string
| create table t (a int) stats_persistent := 1
< ok 0 0
> exec ^
: string
| create table t (a int) stats_persistent := default
< ok 0 0
> exec ^
: string
| create table t (a int) stats_sample_pages := 1
< ok 0 0
> exec ^
: string
| create table t (a int) stats_sample_pages := default
< ok 0 0
> exec ^
: string
| create table t (a int) checksum := 1
< ok 0 0
> exec ^
: string
| create table t (a int) table_checksum := 1
< ok 0 0
> exec ^
: string
| create table t (a int) page_checksum := 1
< error
| 1064 42000 …near 'page_checksum := 1' at line 1
> exec ^
: string
| create table t (a int) delay_key_write := 1
< ok 0 0
> exec ^
: string
| create table t (a int) row_format := dynamic
< ok 0 0
> exec ^
: string
| create table t (a int) union := (u)
< ok 0 0
> exec ^
: string
| create table t (a int) insert_method := no
< ok 0 0
> exec ^
: string
| create table t (a int) data directory := 'x'
< ok 0 0
> exec ^
: string
| create table t (a int) index directory := 'x'
< ok 0 0
> exec ^
: string
| create table t (a int) tablespace := ts
< ok 0 0
> exec ^
: string
| create table t (a int) connection := 'x'
< ok 0 0
> exec ^
: string
| create table t (a int) key_block_size := 1
< ok 0 0
> exec ^
: string
| create table t (a int) engine_attribute := '{}'
< ok 0 0
> exec ^
: string
| create table t (a int) secondary_engine_attribute := '{}'
< ok 0 0
> exec ^
: string
| create table t (a int) transactional := 1
< error
| 1064 42000 …near 'transactional := 1' at line 1
> exec ^
: string
| create table t (a int) sequence := 1
< error
| 1064 42000 …near 'sequence := 1' at line 1
> exec ^
: string
| create table t (a int) character set := utf8mb4
< ok 0 0
> exec ^
: string
| create table t (a int) default character set := utf8mb4
< ok 0 0
> exec ^
: string
| create table t (a int) charset := utf8mb4
< ok 0 0
> exec ^
: string
| create table t (a int) default charset := utf8mb4
< ok 0 0
> exec ^
: string
| create table t (a int) collate := utf8mb4_bin
< ok 0 0
> exec ^
: string
| create table t (a int) default collate := utf8mb4_bin
< ok 0 0
> exec ^
: string
| create table t (a int) engine := innodb, comment := 'x', max_rows = 1
< ok 0 0
> exec ^
: string
| create table t (a int) engine := innodb comment 'x'
< ok 0 0
> exec ^
: string
| alter table t engine := innodb
< ok 0 0
> exec ^
: string
| alter table t comment := 'x'
< ok 0 0
> exec ^
: string
| alter table t auto_increment := 1
< ok 0 0
> exec ^
: string
| alter table t row_format := dynamic
< ok 0 0
> exec ^
: string
| alter table t character set := utf8mb4
< ok 0 0
> exec ^
: string
| alter table t default charset := utf8mb4
< ok 0 0
> exec ^
: string
| alter table t collate := utf8mb4_bin
< ok 0 0
> exec ^
: string
| alter table t key_block_size := 1
< ok 0 0
> exec ^
: string
| alter table t union := (u)
< ok 0 0
> exec ^
: string
| create table t (a int) foo := 'x'
< error
| 1064 42000 …near 'foo := 'x'' at line 1
> exec ^
: string
| create table t (a int) foo := bar
< error
| 1064 42000 …near 'foo := bar' at line 1
> exec ^
: string
| create table t (a int) foo := 1
< error
| 1064 42000 …near 'foo := 1' at line 1
> exec ^
: string
| create table t (a int) foo := default
< error
| 1064 42000 …near 'foo := default' at line 1
> exec ^
: string
| alter table t algorithm := inplace, add column b int
< ok 0 0
> exec ^
: string
| alter table t algorithm := default, add column b int
< ok 0 0
> exec ^
: string
| alter table t lock := none, add column b int
< ok 0 0
> exec ^
: string
| alter table t lock := default, add column b int
< ok 0 0
> exec ^
: string
| alter table t add column b int, algorithm := inplace, lock := none
< ok 0 0
> exec ^
: string
| create index i on t (a) algorithm := inplace
< ok 0 0
> exec ^
: string
| create index i on t (a) lock := none
< ok 0 0
> exec ^
: string
| create index i on t (a) algorithm := inplace lock := none
< ok 0 0
> exec ^
: string
| drop index i on t algorithm := inplace
< ok 0 0
> exec ^
: string
| drop index i on t lock := none
< ok 0 0
> exec ^
: string
| drop index i on t algorithm := inplace lock := none
< ok 0 0
> exec ^
: string
| create index i on t (a) key_block_size := 1
< ok 0 0
> exec ^
: string
| create index i on t (a) engine_attribute := '{}'
< ok 0 0
> exec ^
: string
| create index i on t (a) secondary_engine_attribute := '{}'
< ok 0 0
> exec ^
: string
| create index i on t (a) foo := 1
< error
| 1064 42000 …near 'foo := 1' at line 1
> exec ^
: string
| create table t (a int, key k (a) key_block_size := 1)
< ok 0 0
> exec ^
: string
| create table t (a int, key k (a) engine_attribute := '{}')
< ok 0 0
> exec ^
: string
| create table t (a int, unique key k (a) key_block_size := 1)
< ok 0 0
> exec ^
: string
| create table t (a int, primary key (a) key_block_size := 1)
< ok 0 0
> exec ^
: string
| alter table t add index i (a) key_block_size := 1
< ok 0 0
> exec ^
: string
| alter table t add index i (a) engine_attribute := '{}'
< ok 0 0
> exec ^
: string
| create fulltext index i on t (a) key_block_size := 1
< ok 0 0
> exec ^
: string
| create spatial index i on t (a) key_block_size := 1
< ok 0 0
> exec ^
: string
| create table t (a int engine_attribute := '{}')
< ok 0 0
> exec ^
: string
| create table t (a int secondary_engine_attribute := '{}')
< ok 0 0
> exec ^
: string
| create table t (a int foo := 1)
< error
| 1064 42000 …near 'foo := 1)' at line 1
> exec ^
: string
| create table t (a blob compressed := zlib)
< error
| 1064 42000 …near 'compressed := zlib)' at line 1
> exec ^
: string
| create database d character set := utf8mb4
< ok 0 0
> exec ^
: string
| create database d default character set := utf8mb4
< ok 0 0
> exec ^
: string
| create database d charset := utf8mb4
< ok 0 0
> exec ^
: string
| create database d collate := utf8mb4_bin
< ok 0 0
> exec ^
: string
| create database d default collate := utf8mb4_bin
< ok 0 0
> exec ^
: string
| create database d encryption := 'n'
< ok 0 0
> exec ^
: string
| create database d default encryption := 'n'
< ok 0 0
> exec ^
: string
| create database d comment := 'x'
< error
| 1064 42000 …near 'comment := 'x'' at line 1
> exec ^
: string
| create schema d character set := utf8mb4
< ok 0 0
> exec ^
: string
| alter database d character set := utf8mb4
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| alter database d collate := utf8mb4_bin
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| alter database d encryption := 'n'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| alter database d read only := 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| alter database d read only := default
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| alter database d comment := 'x'
< error
| 1064 42000 …near 'comment := 'x'' at line 1
> exec ^
: string
| alter schema d default charset := utf8mb4
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create table t (a int) partition by hash(a) (partition p0 engine := innodb)
< ok 0 0
> exec ^
: string
| create table t (a int) partition by hash(a) (partition p0 storage engine := innodb)
< ok 0 0
> exec ^
: string
| create table t (a int) partition by hash(a) (partition p0 tablespace := ts)
< ok 0 0
> exec ^
: string
| create table t (a int) partition by hash(a) (partition p0 nodegroup := 1)
< ok 0 0
> exec ^
: string
| create table t (a int) partition by hash(a) (partition p0 max_rows := 1)
< ok 0 0
> exec ^
: string
| create table t (a int) partition by hash(a) (partition p0 min_rows := 1)
< ok 0 0
> exec ^
: string
| create table t (a int) partition by hash(a) (partition p0 data directory := 'x')
< ok 0 0
> exec ^
: string
| create table t (a int) partition by hash(a) (partition p0 index directory := 'x')
< ok 0 0
> exec ^
: string
| create table t (a int) partition by hash(a) (partition p0 comment := 'x')
< ok 0 0
> exec ^
: string
| create table t (a int) partition by hash(a) (partition p0 connection := 'x')
< error
| 1064 42000 …near 'connection := 'x')' at line 1
> exec ^
: string
| create table t (a int) partition by range (a) subpartition by hash (a) (partition p0 values less than (1) (subpartition s0 engine := innodb))
< ok 0 0
> exec ^
: string
| alter table t add partition (partition p1 engine := innodb)
< ok 0 0
> exec ^
: string
| alter table t reorganize partition p0 into (partition p1 values less than (1) comment := 'x')
< ok 0 0
> exec ^
: string
| create tablespace ts add datafile 'a.ibd' initial_size := 1024
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create tablespace ts add datafile 'a.ibd' autoextend_size := 4M
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create tablespace ts add datafile 'a.ibd' max_size := 1024
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create tablespace ts add datafile 'a.ibd' extent_size := 1024
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create tablespace ts add datafile 'a.ibd' nodegroup := 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create tablespace ts add datafile 'a.ibd' comment := 'x'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create tablespace ts add datafile 'a.ibd' engine := innodb
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create tablespace ts add datafile 'a.ibd' file_block_size := 1024
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create tablespace ts add datafile 'a.ibd' encryption := 'y'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create tablespace ts add datafile 'a.ibd' engine_attribute := '{}'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| alter tablespace ts autoextend_size := 4M
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| alter tablespace ts engine := innodb
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| alter tablespace ts encryption := 'y'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| alter tablespace ts initial_size := 1024
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| alter tablespace ts engine_attribute := '{}'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| alter tablespace ts add datafile 'b' initial_size := 1024
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| drop tablespace ts engine := innodb
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create logfile group lg add undofile 'u' initial_size := 1024
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create logfile group lg add undofile 'u' undo_buffer_size := 1024
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create logfile group lg add undofile 'u' redo_buffer_size := 1024
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create logfile group lg add undofile 'u' nodegroup := 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create logfile group lg add undofile 'u' comment := 'x'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create logfile group lg add undofile 'u' engine := ndb
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| alter logfile group lg add undofile 'u' initial_size := 1024
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| alter logfile group lg add undofile 'u' engine := ndb
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| drop logfile group lg engine := ndb
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create resource group g type := user
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create resource group g type := user vcpu := 1 thread_priority := 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| alter resource group g vcpu := 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| alter resource group g thread_priority := 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| clone local data directory := 'x'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create sequence s start := 1
< error
| 1064 42000 …near 'sequence s start := 1' at line 1
> exec ^
: string
| create sequence s increment := 1
< error
| 1064 42000 …near 'sequence s increment := 1' at line 1
> exec ^
: string
| create sequence s minvalue := 1
< error
| 1064 42000 …near 'sequence s minvalue := 1' at line 1
> exec ^
: string
| create sequence s maxvalue := 1
< error
| 1064 42000 …near 'sequence s maxvalue := 1' at line 1
> exec ^
: string
| create sequence s cache := 1
< error
| 1064 42000 …near 'sequence s cache := 1' at line 1
> exec ^
: string
| alter sequence s restart := 1
< error
| 1064 42000 …near 'sequence s restart := 1' at line 1
> exec ^
: string
| alter sequence s minvalue := 1
< error
| 1064 42000 …near 'sequence s minvalue := 1' at line 1
> exec ^
: string
| set autocommit := 1
< ok 0 0
> exec ^
: string
| set @@autocommit := 1
< ok 0 0
> exec ^
: string
| set session autocommit := 1
< ok 0 0
> exec ^
: string
| set global autocommit := 1
< ok 0 0
> exec ^
: string
| set @@global.autocommit := 1
< ok 0 0
> exec ^
: string
| set persist autocommit := 1
< ok 0 0
> exec ^
: string
| set @a := 1
< ok 0 0
> exec ^
: string
| set @a := 1, autocommit := 1
< ok 0 0
> exec ^
: string
| set names := 'utf8mb4'
< error
| 1064 42000 …near ':= 'utf8mb4'' at line 1
> exec ^
: string
| set password := 'x'
< ok 0 0
> exec ^
: string
| set password for u := 'x'
< error
| 1044 42000 Access denied for user 'omni_test'@'%' to database 'mysql'
> exec ^
: string
| set password for u@h := 'x'
< error
| 1044 42000 Access denied for user 'omni_test'@'%' to database 'mysql'
> exec ^
: string
| set statement max_statement_time := 1 for select 1
< error
| 1064 42000 …near 'max_statement_time := 1 for select 1' at line 1
> exec ^
: string
| create definer := u view v as select 1
< error
| 1064 42000 …near ':= u view v as select 1' at line 1
> exec ^
: string
| create algorithm := merge view v as select 1
< error
| 1064 42000 …near ':= merge view v as select 1' at line 1
> exec ^
: string
| alter algorithm := merge view v as select 1
< error
| 1064 42000 …near ':= merge view v as select 1' at line 1
> exec ^
: string
| create definer := u procedure p() select 1
< error
| 1064 42000 …near ':= u procedure p() select 1' at line 1
> exec ^
: string
| alter definer := u view v as select 1
< error
| 1064 42000 …near ':= u view v as select 1' at line 1
> exec ^
: string
| explain format := json select 1
< error
| 1064 42000 …near ':= json select 1' at line 1
> exec ^
: string
| alter table t rename := u
< error
| 1064 42000 …near ':= u' at line 1
> exec ^
: string
| create table t (a int) partition by key algorithm := 2 (a)
< error
| 1064 42000 …near ':= 2 (a)' at line 1
> exec ^
: string
| signal sqlstate '45000' set message_text := 'x'
< error
| 1064 42000 …near ':= 'x'' at line 1
> exec ^
: string
| resignal set message_text := 'x'
< error
| 1064 42000 …near ':= 'x'' at line 1
> exec ^
: string
| get diagnostics @a := number
< error
| 1064 42000 …near ':= number' at line 1
> exec ^
: string
| get diagnostics condition 1 @a := mysql_errno
< error
| 1064 42000 …near ':= mysql_errno' at line 1
> exec ^
: string
| change replication source to source_host := 'x'
< error
| 1064 42000 …near ':= 'x'' at line 1
> exec ^
: string
| change master to master_host := 'x'
< error
| 1064 42000 …near ':= 'x'' at line 1
> exec ^
: string
| change master to master_port := 1
< error
| 1064 42000 …near ':= 1' at line 1
> exec ^
: string
| change master to master_log_file := 'x'
< error
| 1064 42000 …near ':= 'x'' at line 1
> exec ^
: string
| change master to ignore_server_ids := (1)
< error
| 1064 42000 …near ':= (1)' at line 1
> exec ^
: string
| change replication filter replicate_do_db := (d)
< error
| 1064 42000 …near ':= (d)' at line 1
> exec ^
: string
| start replica until source_log_file := 'x', source_log_pos := 1
< error
| 1064 42000 …near ':= 'x', source_log_pos := 1' at line 1
> exec ^
: string
| start slave until master_log_file := 'x', master_log_pos := 1
< error
| 1064 42000 …near ':= 'x', master_log_pos := 1' at line 1
> exec ^
: string
| start replica user := 'x'
< error
| 1064 42000 …near ':= 'x'' at line 1
> exec ^
: string
| start replica password := 'x'
< error
| 1064 42000 …near ':= 'x'' at line 1
> exec ^
: string
| start replica default_auth := 'x'
< error
| 1064 42000 …near ':= 'x'' at line 1
> exec ^
: string
| start replica plugin_dir := 'x'
< error
| 1064 42000 …near ':= 'x'' at line 1
> exec ^
: string
| start replica until sql_before_gtids := 'x'
< error
| 1064 42000 …near ':= 'x'' at line 1
> exec ^
: string
| create table t (a int) engine : = innodb
< error
| 1064 42000 …near ': = innodb' at line 1
> exec ^
: string
| create table t (a int) engine :=
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| create table t (a int) engine :=  = innodb
< error
| 1064 42000 …near '= innodb' at line 1
> exec ^
: string
| create table t (a int) engine = := innodb
< error
| 1064 42000 …near ':= innodb' at line 1
> exec ^
: string
| create table t (a int) charset = binary
< ok 0 0
> exec ^
: string
| create table t (a int) charset binary
< ok 0 0
> exec ^
: string
| create table t (a int) charset := binary
< ok 0 0
> exec ^
: string
| create table t (a int) default charset = binary
< ok 0 0
> exec ^
: string
| create table t (a int) character set = binary
< ok 0 0
> exec ^
: string
| create table t (a int) character set binary
< ok 0 0
> exec ^
: string
| create table t (a int) default character set binary
< ok 0 0
> exec ^
: string
| create table t (a int) collate = binary
< ok 0 0
> exec ^
: string
| create table t (a int) collate binary
< ok 0 0
> exec ^
: string
| create table t (a int) collate := binary
< ok 0 0
> exec ^
: string
| create table t (a int) default collate = binary
< ok 0 0
> exec ^
: string
| create table t (a int) charset = default
< error
| 1064 42000 …near 'default' at line 1
> exec ^
: string
| create table t (a int) character set = default
< error
| 1064 42000 …near 'default' at line 1
> exec ^
: string
| create table t (a int) collate = default
< error
| 1064 42000 …near 'default' at line 1
> exec ^
: string
| create table t (a int) charset = 'utf8mb4'
< ok 0 0
> exec ^
: string
| create table t (a int) collate = 'utf8mb4_bin'
< ok 0 0
> exec ^
: string
| create table t (a int) charset = select
< error
| 1064 42000 …near 'select' at line 1
> exec ^
: string
| create table t (a int) collate = select
< error
| 1064 42000 …near 'select' at line 1
> exec ^
: string
| create table t (a int) charset = null
< error
| 1064 42000 …near 'null' at line 1
> exec ^
: string
| alter table t charset = binary
< ok 0 0
> exec ^
: string
| alter table t character set = binary
< ok 0 0
> exec ^
: string
| alter table t collate = binary
< ok 0 0
> exec ^
: string
| alter table t default charset binary
< ok 0 0
> exec ^
: string
| alter table t convert to character set binary
< ok 0 0
> exec ^
: string
| alter table t convert to charset binary
< ok 0 0
> exec ^
: string
| alter table t convert to character set binary collate binary
< ok 0 0
> exec ^
: string
| alter table t convert to character set utf8mb4 collate binary
< error
| 1253 42000 COLLATION 'binary' is not valid for CHARACTER SET 'utf8mb4'
> exec ^
: string
| alter table t convert to character set default
< ok 0 0
> exec ^
: string
| alter table t convert to character set default collate utf8mb4_bin
< ok 0 0
> exec ^
: string
| alter table t convert to character set select
< error
| 1064 42000 …near 'select' at line 1
> exec ^
: string
| alter table t convert to character set utf8mb4 collate select
< error
| 1064 42000 …near 'select' at line 1
> exec ^
: string
| create database d charset binary
< ok 0 0
> exec ^
: string
| create database d collate binary
< ok 0 0
> exec ^
: string
| create table t (a int) engine = binary
< error
| 1064 42000 …near 'binary' at line 1
> exec ^
: string
| create table t (a int) engine = select
< error
| 1064 42000 …near 'select' at line 1
> exec ^
: string
| create table t (a int) row_format = default
< ok 0 0
> exec ^
: string
| create table t (a int) row_format = binary
< error
| 1064 42000 …near 'binary' at line 1
> exec ^
: string
| create table t (a int) row_format = select
< error
| 1064 42000 …near 'select' at line 1
> exec ^
: string
| create table t (a int) compression = binary
< error
| 1064 42000 …near 'binary' at line 1
> exec ^
: string
| create table t (a int) insert_method = first
< ok 0 0
> exec ^
: string
| create table t (a int) insert_method = last
< ok 0 0
> exec ^
: string
| create table t (a int) insert_method = no
< ok 0 0
> exec ^
: string
| create table t (a int) tablespace = binary
< error
| 1064 42000 …near 'binary' at line 1
> exec ^
: string
| create table t (a int) tablespace binary
< error
| 1064 42000 …near 'binary' at line 1
> exec ^
: string
| create table t (a int) storage binary
< error
| 1064 42000 …near 'binary' at line 1
> exec ^
: string
| create table t (a int) secondary_engine = binary
< error
| 1064 42000 …near 'binary' at line 1
> exec ^
: string
| create table t (a int) pack_keys = binary
< error
| 1064 42000 …near 'binary' at line 1
> exec ^
: string
| create table t (a int) max_rows = binary
< error
| 1064 42000 …near 'binary' at line 1
> exec ^
: string
| create table t (a int) foo = binary
< error
| 1064 42000 …near 'foo = binary' at line 1
> exec ^
: string
| create table t (a int) foo = default
< error
| 1064 42000 …near 'foo = default' at line 1
> exec ^
: string
| create table t (a int) partition by hash(a) (partition p engine = binary)
< error
| 1064 42000 …near 'binary)' at line 1
> exec ^
: string
| create table t (a int) partition by hash(a) (partition p tablespace = binary)
< error
| 1064 42000 …near 'binary)' at line 1
> exec ^
: string
| alter table t algorithm = default, add column b int
< ok 0 0
> exec ^
: string
| alter table t lock = default, add column b int
< ok 0 0
> exec ^
: string
| alter table t algorithm = binary, add column b int
< error
| 1064 42000 …near 'binary, add column b int' at line 1
> exec ^
: string
| alter table t lock = binary, add column b int
< error
| 1064 42000 …near 'binary, add column b int' at line 1
> exec ^
: string
| alter table t algorithm = select, add column b int
< error
| 1064 42000 …near 'select, add column b int' at line 1
> exec ^
: string
| set password 'x'
< error
| 1064 42000 …near ''x'' at line 1
> exec ^
: string
| set password for u 'x'
< error
| 1064 42000 …near ''x'' at line 1
> exec ^
: string
| create table t (a int) foo := null
< error
| 1064 42000 …near 'foo := null' at line 1
> exec ^
: string
| create table t (a int) foo := select
< error
| 1064 42000 …near 'foo := select' at line 1
> exec ^
: string
| create table t (a int) pack_keys := null
< error
| 1064 42000 …near 'null' at line 1
> exec ^
: string
| create table t (a int) max_rows := null
< error
| 1064 42000 …near 'null' at line 1
> exec ^
: string
| create table t (a int) comment := null
< error
| 1064 42000 …near 'null' at line 1
> exec ^
: string
| create table t (a int) secondary_engine := 'x'
< ok 0 0
> exec ^
: string
| create table t (a int) auto_increment := 1 engine := innodb
< ok 0 0
> exec ^
: string
| create tablespace ts add datafile 'a.ibd' engine := ndb nodegroup := 1 comment := 'x'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| clone instance from 'u'@'h':1 identified by 'x' data directory := 'x'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create index i on t (a) key_block_size := 1 comment 'x'
< ok 0 0
> exec ^
: string
| drop index i on t algorithm := default
< ok 0 0
> exec ^
: string
| create table t (a int) storage disk
< ok 0 0
> exec ^
: string
| create table t (a int) storage memory
< ok 0 0
> exec ^
: string
| create table t (a int) storage := disk
< error
| 1064 42000 …near ':= disk' at line 1
> exec ^
: string
| create table t (a int) storage = disk
< error
| 1064 42000 …near '= disk' at line 1
> exec ^
: string
| create table t (a int) storage := memory
< error
| 1064 42000 …near ':= memory' at line 1
> exec ^
: string
| create table t (a int) storage engine := innodb
< error
| 1064 42000 …near 'engine := innodb' at line 1
