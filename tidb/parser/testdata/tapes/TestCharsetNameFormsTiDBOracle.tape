-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestCharsetNameFormsTiDBOracle default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| alter table t add c char charset 'binary'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add c char charset 'binary'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add c char charset 'binary' binary
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add c char charset 'binary' binary
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add c char binary charset 'binary'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add c char binary charset 'binary'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add c char character set 'utf8mb4'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add c char character set 'utf8mb4'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add c char char set 'utf8mb4'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add c char char set 'utf8mb4'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add c char charset "utf8mb4"
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add c char charset "utf8mb4"
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add c char charset 'UTF8MB4'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add c char charset 'UTF8MB4'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add c char charset utf8mb4
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add c char charset utf8mb4
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add c char charset `utf8mb4`
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add c char charset `utf8mb4`
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add c char charset binary
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add c char charset binary
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add c char collate 'utf8mb4_bin'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add c char collate 'utf8mb4_bin'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add c char collate 'binary'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add c char collate 'binary'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add c char collate binary
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add c char collate binary
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add c char charset 'utf8mb4' collate 'utf8mb4_bin'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add c char charset 'utf8mb4' collate 'utf8mb4_bin'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add c char charset utf8mb4 collate 'utf8mb4_bin'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add c char charset utf8mb4 collate 'utf8mb4_bin'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add c char charset 'utf8mb4' collate utf8mb4_bin
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add c char charset 'utf8mb4' collate utf8mb4_bin
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add c char binary collate 'utf8mb4_bin'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add c char binary collate 'utf8mb4_bin'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add c char charset 'utf8mb4' binary collate 'utf8mb4_bin'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add c char charset 'utf8mb4' binary collate 'utf8mb4_bin'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add c char charset 'utf8mb4'collate'utf8mb4_bin'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add c char charset 'utf8mb4'collate'utf8mb4_bin'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add c varchar(10) charset 'utf8mb4'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add c varchar(10) charset 'utf8mb4'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add c text charset 'utf8mb4'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add c text charset 'utf8mb4'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add c text collate 'utf8mb4_bin'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add c text collate 'utf8mb4_bin'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add c text ascii collate 'latin1_bin'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add c text ascii collate 'latin1_bin'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add c text byte collate 'binary'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add c text byte collate 'binary'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add c enum('a') charset 'utf8mb4'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add c enum('a') charset 'utf8mb4'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add c set('a') charset 'utf8mb4'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add c set('a') charset 'utf8mb4'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add c long charset 'utf8mb4'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add c long charset 'utf8mb4'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add c long varchar charset 'utf8mb4'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add c long varchar charset 'utf8mb4'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add c nchar(1) charset 'utf8mb4'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add c nchar(1) charset 'utf8mb4'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add c national char(1) charset 'utf8mb4'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add c national char(1) charset 'utf8mb4'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t modify c char(10) charset 'utf8mb4'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t modify c char(10) charset 'utf8mb4'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t change c c char(10) charset 'utf8mb4'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t change c c char(10) charset 'utf8mb4'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t2 (c char charset 'utf8mb4')
< ok 0 0
> exec ^
: string
| create table t2 (c char(1) charset 'utf8mb4' collate 'utf8mb4_bin')
< ok 0 0
> exec ^
: string
| create table t2 (c char(1) charset 'utf8mb4' not null)
< ok 0 0
> exec ^
: string
| alter table t add c char charset 'utf8mb4' default 'x'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add c char charset 'utf8mb4' default 'x'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add c char(1) charset 'utf8mb4' collate 'utf8mb4_bin' collate 'utf8mb4_bin'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add c char(1) charset 'utf8mb4' collate 'utf8mb4_bin' collate 'utf8mb4_bin'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add c char charset 'utf8' 'mb4'
< error
| 1064 42000 …line 1 column 45 near "'mb4'" 
> exec ^
: string
| alter table t add c char charset 'utf8mb4' 'x'
< error
| 1064 42000 …line 1 column 46 near "'x'" 
> exec ^
: string
| alter table t add c char charset N'utf8mb4'
< error
| 1064 42000 …line 1 column 34 near "N'utf8mb4'" 
> exec ^
: string
| alter table t add c char charset _utf8mb4'utf8mb4'
< error
| 1064 42000 …line 1 column 41 near "_utf8mb4'utf8mb4'" 
> exec ^
: string
| alter table t add c char charset 0x62
< error
| 1064 42000 …line 1 column 37 near "0x62" 
> exec ^
: string
| alter table t add c char charset x'62'
< error
| 1064 42000 …line 1 column 38 near "x'62'" 
> exec ^
: string
| alter table t add c char charset 1
< error
| 1064 42000 …line 1 column 34 near "1" 
> exec ^
: string
| alter table t add c char charset null
< error
| 1064 42000 …line 1 column 37 near "null" 
> exec ^
: string
| alter table t add c char charset ?
< error
| 1064 42000 …line 1 column 34 near "?" 
> exec ^
: string
| alter table t add c char charset @a
< error
| 1064 42000 …line 1 column 35 near "@a" 
> exec ^
: string
| alter table t add c char charset _utf8mb4
< error
| 1064 42000 …line 1 column 41 near "_utf8mb4" 
> exec ^
: string
| alter table t add c char collate @a
< error
| 1064 42000 …line 1 column 35 near "@a" 
> exec ^
: string
| alter table t add c char collate _utf8mb4
< error
| 1064 42000 …line 1 column 41 near "_utf8mb4" 
> exec ^
: string
| alter table t add c char charset
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| alter table t add c char charset default
< error
| 1064 42000 …line 1 column 40 near "default" 
> exec ^
: string
| alter table t add c char charset = 'utf8mb4'
< error
| 1064 42000 …line 1 column 34 near "= 'utf8mb4'" 
> exec ^
: string
| alter table t add c char charset 'utf8mb4' charset 'utf8mb4'
< error
| 1064 42000 …line 1 column 50 near "charset 'utf8mb4'" 
> exec ^
: string
| alter table t add c char collate 'utf8mb4_bin' charset 'utf8mb4'
< error
| 1064 42000 …line 1 column 54 near "charset 'utf8mb4'" 
> exec ^
: string
| alter table t add c char collate 'utf8mb4_bin' binary
< error
| 1064 42000 …line 1 column 53 near "binary" 
> exec ^
: string
| alter table t add c char charset 'utf8mb4' binary binary
< error
| 1064 42000 …line 1 column 56 near "binary" 
> exec ^
: string
| alter table t add c char binary charset 'utf8mb4' binary
< error
| 1064 42000 …line 1 column 56 near "binary" 
> exec ^
: string
| alter table t add c int charset 'utf8mb4'
< error
| 1064 42000 …line 1 column 31 near "charset 'utf8mb4'" 
