-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestOptionEqualFormsMatchEngine default#1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) engine := innodb
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) secondary_engine := null
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'null' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) secondary_engine := x
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) max_rows := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) min_rows := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) avg_row_length := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) password := 'x'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) comment := 'x'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) compression := 'zlib'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) encryption := 'y'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) auto_increment := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) pack_keys := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) pack_keys := default
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) stats_auto_recalc := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) stats_auto_recalc := default
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) stats_persistent := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) stats_persistent := default
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) stats_sample_pages := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) stats_sample_pages := default
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) checksum := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) table_checksum := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) page_checksum := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) delay_key_write := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) row_format := dynamic
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) union := (u)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) insert_method := no
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) data directory := 'x'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) index directory := 'x'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) tablespace := ts
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= ts' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) connection := 'x'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) key_block_size := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) engine_attribute := '{}'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) secondary_engine_attribute := '{}'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) transactional := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) character set := utf8mb4
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) default character set := utf8mb4
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) charset := utf8mb4
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) default charset := utf8mb4
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) collate := utf8mb4_bin
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) default collate := utf8mb4_bin
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) engine := innodb, comment := 'x', max_rows = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) engine := innodb comment 'x'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter table t engine := innodb
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter table t comment := 'x'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter table t auto_increment := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter table t row_format := dynamic
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter table t character set := utf8mb4
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter table t default charset := utf8mb4
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter table t collate := utf8mb4_bin
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter table t key_block_size := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter table t union := (u)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) foo := 'x'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) foo := bar
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) foo := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) foo := default
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter table t algorithm := inplace, add column b int
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter table t algorithm := default, add column b int
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter table t lock := none, add column b int
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter table t lock := default, add column b int
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter table t add column b int, algorithm := inplace, lock := none
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create index i on t (a) algorithm := inplace
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create index i on t (a) lock := none
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create index i on t (a) algorithm := inplace lock := none
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| drop index i on t algorithm := inplace
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'algorithm := inplace' at line 1
> exec
| SET @omni_sweep = ?
: string
| drop index i on t lock := none
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'lock := none' at line 1
> exec
| SET @omni_sweep = ?
: string
| drop index i on t algorithm := inplace lock := none
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'algorithm := inplace lock := none' at line 1
> exec
| SET @omni_sweep = ?
: string
| create index i on t (a) key_block_size := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create index i on t (a) engine_attribute := '{}'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create index i on t (a) secondary_engine_attribute := '{}'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int, key k (a) key_block_size := 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int, key k (a) engine_attribute := '{}')
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int, unique key k (a) key_block_size := 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int, primary key (a) key_block_size := 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter table t add index i (a) key_block_size := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter table t add index i (a) engine_attribute := '{}'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create fulltext index i on t (a) key_block_size := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create spatial index i on t (a) key_block_size := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int engine_attribute := '{}')
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int secondary_engine_attribute := '{}')
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create database d character set := utf8mb4
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create database d default character set := utf8mb4
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create database d charset := utf8mb4
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create database d collate := utf8mb4_bin
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create database d default collate := utf8mb4_bin
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create database d encryption := 'n'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'encryption := 'n'' at line 1
> exec
| SET @omni_sweep = ?
: string
| create database d default encryption := 'n'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'encryption := 'n'' at line 1
> exec
| SET @omni_sweep = ?
: string
| create schema d character set := utf8mb4
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter database d character set := utf8mb4
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter database d collate := utf8mb4_bin
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter database d encryption := 'n'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'encryption := 'n'' at line 1
> exec
| SET @omni_sweep = ?
: string
| alter database d read only := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'read only := 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| alter database d read only := default
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'read only := default' at line 1
> exec
| SET @omni_sweep = ?
: string
| alter schema d default charset := utf8mb4
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by hash(a) (partition p0 engine := innodb)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by hash(a) (partition p0 storage engine := innodb)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by hash(a) (partition p0 tablespace := ts)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by hash(a) (partition p0 nodegroup := 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by hash(a) (partition p0 max_rows := 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by hash(a) (partition p0 min_rows := 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by hash(a) (partition p0 data directory := 'x')
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by hash(a) (partition p0 index directory := 'x')
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by hash(a) (partition p0 comment := 'x')
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by range (a) subpartition by hash (a) (partition p0 values less than (1) (subpartition s0 engine := innodb))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter table t add partition (partition p1 engine := innodb)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter table t reorganize partition p0 into (partition p1 values less than (1) comment := 'x')
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create tablespace ts add datafile 'a.ibd' initial_size := 1024
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'tablespace ts add datafile 'a.ibd' initial_size := 1024' at line 1
> exec
| SET @omni_sweep = ?
: string
| create tablespace ts add datafile 'a.ibd' autoextend_size := 4M
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'tablespace ts add datafile 'a.ibd' autoextend_size := 4M' at line 1
> exec
| SET @omni_sweep = ?
: string
| create tablespace ts add datafile 'a.ibd' max_size := 1024
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'tablespace ts add datafile 'a.ibd' max_size := 1024' at line 1
> exec
| SET @omni_sweep = ?
: string
| create tablespace ts add datafile 'a.ibd' extent_size := 1024
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'tablespace ts add datafile 'a.ibd' extent_size := 1024' at line 1
> exec
| SET @omni_sweep = ?
: string
| create tablespace ts add datafile 'a.ibd' nodegroup := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'tablespace ts add datafile 'a.ibd' nodegroup := 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| create tablespace ts add datafile 'a.ibd' comment := 'x'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'tablespace ts add datafile 'a.ibd' comment := 'x'' at line 1
> exec
| SET @omni_sweep = ?
: string
| create tablespace ts add datafile 'a.ibd' engine := innodb
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'tablespace ts add datafile 'a.ibd' engine := innodb' at line 1
> exec
| SET @omni_sweep = ?
: string
| create tablespace ts add datafile 'a.ibd' storage engine := innodb
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'tablespace ts add datafile 'a.ibd' storage engine := innodb' at line 1
> exec
| SET @omni_sweep = ?
: string
| create tablespace ts add datafile 'a.ibd' file_block_size := 1024
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'tablespace ts add datafile 'a.ibd' file_block_size := 1024' at line 1
> exec
| SET @omni_sweep = ?
: string
| create tablespace ts add datafile 'a.ibd' encryption := 'y'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'tablespace ts add datafile 'a.ibd' encryption := 'y'' at line 1
> exec
| SET @omni_sweep = ?
: string
| create tablespace ts add datafile 'a.ibd' engine_attribute := '{}'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'tablespace ts add datafile 'a.ibd' engine_attribute := '{}'' at line 1
> exec
| SET @omni_sweep = ?
: string
| alter tablespace ts autoextend_size := 4M
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'tablespace ts autoextend_size := 4M' at line 1
> exec
| SET @omni_sweep = ?
: string
| alter tablespace ts engine := innodb
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'tablespace ts engine := innodb' at line 1
> exec
| SET @omni_sweep = ?
: string
| alter tablespace ts encryption := 'y'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'tablespace ts encryption := 'y'' at line 1
> exec
| SET @omni_sweep = ?
: string
| alter tablespace ts initial_size := 1024
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'tablespace ts initial_size := 1024' at line 1
> exec
| SET @omni_sweep = ?
: string
| alter tablespace ts engine_attribute := '{}'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'tablespace ts engine_attribute := '{}'' at line 1
> exec
| SET @omni_sweep = ?
: string
| alter tablespace ts add datafile 'b' initial_size := 1024
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'tablespace ts add datafile 'b' initial_size := 1024' at line 1
> exec
| SET @omni_sweep = ?
: string
| drop tablespace ts engine := innodb
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'tablespace ts engine := innodb' at line 1
> exec
| SET @omni_sweep = ?
: string
| create logfile group lg add undofile 'u' initial_size := 1024
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'logfile group lg add undofile 'u' initial_size := 1024' at line 1
> exec
| SET @omni_sweep = ?
: string
| create logfile group lg add undofile 'u' undo_buffer_size := 1024
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'logfile group lg add undofile 'u' undo_buffer_size := 1024' at line 1
> exec
| SET @omni_sweep = ?
: string
| create logfile group lg add undofile 'u' redo_buffer_size := 1024
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'logfile group lg add undofile 'u' redo_buffer_size := 1024' at line 1
> exec
| SET @omni_sweep = ?
: string
| create logfile group lg add undofile 'u' nodegroup := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'logfile group lg add undofile 'u' nodegroup := 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| create logfile group lg add undofile 'u' comment := 'x'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'logfile group lg add undofile 'u' comment := 'x'' at line 1
> exec
| SET @omni_sweep = ?
: string
| create logfile group lg add undofile 'u' engine := ndb
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'logfile group lg add undofile 'u' engine := ndb' at line 1
> exec
| SET @omni_sweep = ?
: string
| alter logfile group lg add undofile 'u' initial_size := 1024
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'logfile group lg add undofile 'u' initial_size := 1024' at line 1
> exec
| SET @omni_sweep = ?
: string
| alter logfile group lg add undofile 'u' engine := ndb
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'logfile group lg add undofile 'u' engine := ndb' at line 1
> exec
| SET @omni_sweep = ?
: string
| drop logfile group lg engine := ndb
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'logfile group lg engine := ndb' at line 1
> exec
| SET @omni_sweep = ?
: string
| create resource group g type := user
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'resource group g type := user' at line 1
> exec
| SET @omni_sweep = ?
: string
| create resource group g type := user vcpu := 1 thread_priority := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'resource group g type := user vcpu := 1 thread_priority := 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| alter resource group g vcpu := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'resource group g vcpu := 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| alter resource group g thread_priority := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'resource group g thread_priority := 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| clone local data directory := 'x'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'clone local data directory := 'x'' at line 1
> exec
| SET @omni_sweep = ?
: string
| clone instance from u@h:1 identified by 'x' data directory := 'x'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'clone instance from u@h:1 identified by 'x' data directory := 'x'' at line 1
> exec
| SET @omni_sweep = ?
: string
| create sequence s start := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create sequence s increment := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create sequence s minvalue := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create sequence s maxvalue := 100
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create sequence s cache := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter sequence s restart := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter sequence s minvalue := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| set autocommit := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| set @@autocommit := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| set session autocommit := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| set global autocommit := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| set @@global.autocommit := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| set @a := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| set @a := 1, autocommit := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| set names := 'utf8mb4'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '' at line 1
> exec
| SET @omni_sweep = ?
: string
| set password := 'x'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| set password for u := 'x'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| set password for u@h := 'x'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| install component 'c' set global autocommit := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'component 'c' set global autocommit := 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| create definer := u view v as select 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= u view v as select 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| create algorithm := merge view v as select 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= merge view v as select 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| alter algorithm := merge view v as select 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= merge view v as select 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| create definer := u procedure p() select 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= u procedure p() select 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| alter definer := u view v as select 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= u view v as select 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| explain format := json select 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= json select 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| alter table t rename := u
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= u' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by key algorithm := 2 (a)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= 2 (a)' at line 1
> exec
| SET @omni_sweep = ?
: string
| signal sqlstate '45000' set message_text := 'x'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= 'x'' at line 1
> exec
| SET @omni_sweep = ?
: string
| resignal set message_text := 'x'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= 'x'' at line 1
> exec
| SET @omni_sweep = ?
: string
| get diagnostics @a := number
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= number' at line 1
> exec
| SET @omni_sweep = ?
: string
| get diagnostics condition 1 @a := mysql_errno
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= mysql_errno' at line 1
> exec
| SET @omni_sweep = ?
: string
| change replication source to source_host := 'x'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'replication source to source_host := 'x'' at line 1
> exec
| SET @omni_sweep = ?
: string
| change master to master_host := 'x'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= 'x'' at line 1
> exec
| SET @omni_sweep = ?
: string
| change master to master_port := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| change master to master_log_file := 'x'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= 'x'' at line 1
> exec
| SET @omni_sweep = ?
: string
| change master to ignore_server_ids := (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= (1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| change replication filter replicate_do_db := (d)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'replication filter replicate_do_db := (d)' at line 1
> exec
| SET @omni_sweep = ?
: string
| start replica until source_log_file := 'x', source_log_pos := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'source_log_file := 'x', source_log_pos := 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| start slave until master_log_file := 'x', master_log_pos := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= 'x', master_log_pos := 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| start replica user := 'x'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'user := 'x'' at line 1
> exec
| SET @omni_sweep = ?
: string
| start replica password := 'x'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'password := 'x'' at line 1
> exec
| SET @omni_sweep = ?
: string
| start replica default_auth := 'x'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'default_auth := 'x'' at line 1
> exec
| SET @omni_sweep = ?
: string
| start replica plugin_dir := 'x'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'plugin_dir := 'x'' at line 1
> exec
| SET @omni_sweep = ?
: string
| start replica until sql_before_gtids := 'x'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= 'x'' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) engine : = innodb
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ': = innodb' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) engine :=
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) engine :=  = innodb
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '= innodb' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) engine = := innodb
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= innodb' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) charset = binary
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) charset binary
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) charset := binary
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) default charset = binary
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) character set = binary
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) character set binary
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) default character set binary
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) collate = binary
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) collate binary
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) collate := binary
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) default collate = binary
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) charset = default
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) character set = default
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) collate = default
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) charset = 'utf8mb4'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) collate = 'utf8mb4_bin'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) charset = select
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'select' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) collate = select
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'select' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) charset = null
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'null' at line 1
> exec
| SET @omni_sweep = ?
: string
| alter table t charset = binary
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter table t character set = binary
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter table t collate = binary
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter table t default charset binary
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter table t convert to character set binary
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter table t convert to charset binary
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter table t convert to character set binary collate binary
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter table t convert to character set utf8mb4 collate binary
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1253 42000 COLLATION 'binary' is not valid for CHARACTER SET 'utf8mb4'
> exec
| SET @omni_sweep = ?
: string
| alter table t convert to character set default
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter table t convert to character set default collate utf8mb4_bin
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter table t convert to character set select
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'select' at line 1
> exec
| SET @omni_sweep = ?
: string
| alter table t convert to character set utf8mb4 collate select
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'select' at line 1
> exec
| SET @omni_sweep = ?
: string
| create database d charset binary
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create database d collate binary
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) engine = binary
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'binary' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) engine = select
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'select' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) row_format = default
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) row_format = binary
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'binary' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) row_format = select
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'select' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) compression = binary
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'binary' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) insert_method = first
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) insert_method = last
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) insert_method = no
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) tablespace = binary
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '= binary' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) tablespace binary
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'binary' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) storage binary
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'binary' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) secondary_engine = binary
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'binary' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) pack_keys = binary
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'binary' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) max_rows = binary
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'binary' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) foo = binary
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'binary' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) foo = default
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by hash(a) (partition p engine = binary)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'binary)' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by hash(a) (partition p tablespace = binary)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'binary)' at line 1
> exec
| SET @omni_sweep = ?
: string
| alter table t algorithm = default, add column b int
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter table t lock = default, add column b int
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter table t algorithm = binary, add column b int
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'binary, add column b int' at line 1
> exec
| SET @omni_sweep = ?
: string
| alter table t lock = binary, add column b int
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'binary, add column b int' at line 1
> exec
| SET @omni_sweep = ?
: string
| alter table t algorithm = select, add column b int
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'select, add column b int' at line 1
> exec
| SET @omni_sweep = ?
: string
| set password 'x'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ''x'' at line 1
> exec
| SET @omni_sweep = ?
: string
| set password for u 'x'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ''x'' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) foo := null
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'null' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) foo := select
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'select' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) pack_keys := null
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'null' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) max_rows := null
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'null' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) comment := null
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'null' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) secondary_engine := 'x'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) auto_increment := 1 engine := innodb
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create tablespace ts add datafile 'a.ibd' engine := ndb nodegroup := 1 comment := 'x'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'tablespace ts add datafile 'a.ibd' engine := ndb nodegroup := 1 comment := 'x'' at line 1
> exec
| SET @omni_sweep = ?
: string
| clone instance from 'u'@'h':1 identified by 'x' data directory := 'x'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'clone instance from 'u'@'h':1 identified by 'x' data directory := 'x'' at line 1
> exec
| SET @omni_sweep = ?
: string
| create index i on t (a) key_block_size := 1 comment 'x'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| drop index i on t algorithm := default
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'algorithm := default' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) storage disk
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) storage memory
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) storage := disk
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= disk' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) storage = disk
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '= disk' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) storage := memory
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= memory' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) storage engine := innodb
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'engine := innodb' at line 1
