-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestBRIEFormsOracle default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| BACKUP DATABASE * TO 'noop://'
< ok 0 0
> exec ^
: string
| BACKUP SCHEMA * TO 'noop://'
< ok 0 0
> exec ^
: string
| BACKUP DATABASE a TO 'noop://'
< ok 0 0
> exec ^
: string
| BACKUP DATABASE a, b, c TO 'noop://'
< ok 0 0
> exec ^
: string
| BACKUP SCHEMA a TO 'noop://'
< ok 0 0
> exec ^
: string
| BACKUP TABLE a TO 'noop://'
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| BACKUP TABLE a TO 'noop://'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE a.b TO 'noop://'
< error
| 1146 42S02 Table 'a.b' doesn't exist
> exec ^
: string
| BACKUP TABLE a.b TO 'noop://'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE a.b, c.d, e TO 'noop://'
< error
| 1146 42S02 Table 'a.b' doesn't exist
> exec ^
: string
| BACKUP TABLE a.b, c.d, e TO 'noop://'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE `a`.`b` TO 'noop://'
< error
| 1146 42S02 Table 'a.b' doesn't exist
> exec ^
: string
| BACKUP TABLE `a`.`b` TO 'noop://'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE *.a TO 'noop://'
< error
| 1146 42S02 Table '*.a' doesn't exist
> exec ^
: string
| BACKUP TABLE *.a TO 'noop://'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP DATABASE `*` TO 'noop://'
< ok 0 0
> exec ^
: string
| RESTORE DATABASE * FROM 'noop://'
< ok 0 0
> exec ^
: string
| RESTORE DATABASE a, b FROM 'noop://'
< ok 0 0
> exec ^
: string
| RESTORE SCHEMA a FROM 'noop://'
< ok 0 0
> exec ^
: string
| RESTORE TABLE a.b FROM 'noop://'
< ok 0 0
> exec ^
: string
| BACKUP DATABASE * TO "noop://"
< ok 0 0
> exec ^
: string
| BACKUP TO 'noop://'
< error
| 1064 42000 …line 1 column 9 near "TO 'noop://'" 
> exec ^
: string
| BACKUP DATABASE TO 'noop://'
< error
| 1064 42000 …line 1 column 18 near "TO 'noop://'" 
> exec ^
: string
| BACKUP TABLE TO 'noop://'
< error
| 1064 42000 …line 1 column 15 near "TO 'noop://'" 
> exec ^
: string
| BACKUP DATABASE *, a TO 'noop://'
< error
| 1064 42000 …line 1 column 18 near ", a TO 'noop://'" 
> exec ^
: string
| BACKUP DATABASE a, * TO 'noop://'
< error
| 1064 42000 …line 1 column 20 near "* TO 'noop://'" 
> exec ^
: string
| BACKUP DATABASE a.b TO 'noop://'
< error
| 1064 42000 …line 1 column 18 near ".b TO 'noop://'" 
> exec ^
: string
| BACKUP TABLE * TO 'noop://'
< error
| 1064 42000 …line 1 column 17 near "TO 'noop://'" 
> exec ^
: string
| BACKUP TABLE a.* TO 'noop://'
< error
| 1064 42000 …line 1 column 16 near "* TO 'noop://'" 
> exec ^
: string
| BACKUP DATABASE a, TO 'noop://'
< error
| 1064 42000 …line 1 column 21 near "TO 'noop://'" 
> exec ^
: string
| BACKUP DATABASES a TO 'noop://'
< error
| 1064 42000 …line 1 column 16 near "DATABASES a TO 'noop://'" 
> exec ^
: string
| BACKUP DATABASE a
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| BACKUP DATABASE a TO
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| BACKUP DATABASE a TO 'noop' '://'
< error
| 1064 42000 …line 1 column 33 near "'://'" 
> exec ^
: string
| BACKUP DATABASE a TO noop
< error
| 1064 42000 …line 1 column 25 near "noop" 
> exec ^
: string
| BACKUP DATABASE a TO 1
< error
| 1064 42000 …line 1 column 22 near "1" 
> exec ^
: string
| BACKUP DATABASE a FROM 'noop://'
< error
| 1064 42000 …line 1 column 22 near "FROM 'noop://'" 
> exec ^
: string
| RESTORE DATABASE a TO 'noop://'
< error
| 1064 42000 …line 1 column 21 near "TO 'noop://'" 
> exec ^
: string
| BACKUP DATABASE a TO 'noop://', b TO 'noop://'
< error
| 1064 42000 …line 1 column 31 near ", b TO 'noop://'" 
> exec ^
: string
| BACKUP LOGS TO 'noop://'
< ok 0 0
> exec ^
: string
| BACKUP LOGS TO 'noop://' start_ts = '20220304'
< ok 0 0
> exec ^
: string
| BACKUP LOGS
< error
| 1064 42000 …line 1 column 11 near "" 
> exec ^
: string
| BACKUP LOGS TO
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| BACKUP LOGS a TO 'noop://'
< error
| 1064 42000 …line 1 column 13 near "a TO 'noop://'" 
> exec ^
: string
| RESTORE POINT FROM 'noop://'
< ok 0 0
> exec ^
: string
| RESTORE POINT FROM 'noop://' full_backup_storage = 'noop://full'
< ok 0 0
> exec ^
: string
| RESTORE POINT FROM 'noop://' full_backup_storage='noop://full' start_ts='20230101' restored_ts='20230123'
< ok 0 0
> exec ^
: string
| RESTORE POINT
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| RESTORE POINT TO 'noop://'
< error
| 1064 42000 …line 1 column 16 near "TO 'noop://'" 
> exec ^
: string
| RESTORE POINTS FROM 'noop://'
< error
| 1064 42000 …line 1 column 14 near "POINTS FROM 'noop://'" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' concurrency 4
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| BACKUP TABLE a TO 'noop://' concurrency 4
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' concurrency = 4
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| BACKUP TABLE a TO 'noop://' concurrency = 4
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' CONCURRENCY 4
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| BACKUP TABLE a TO 'noop://' CONCURRENCY 4
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' concurrency 4 concurrency 8
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| BACKUP TABLE a TO 'noop://' concurrency 4 concurrency 8
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' checksum_concurrency 4 compression_level 4 ignore_stats 1 compression_type 'lz4'
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| BACKUP TABLE a TO 'noop://' checksum_concurrency 4 compression_level 4 ignore_stats 1 compression_type 'lz4'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| RESTORE TABLE a FROM 'noop://' concurrency 40 checksum 0 online 1
< ok 0 0
> exec ^
: string
| RESTORE TABLE a FROM 'noop://' checksum_concurrency 4 wait_tiflash_ready 1 with_sys_table 1
< ok 0 0
> exec ^
: string
| BACKUP TABLE a TO 'noop://' resume 1
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| BACKUP TABLE a TO 'noop://' resume 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' nosuchopt 1
< error
| 1064 42000 …line 1 column 37 near "nosuchopt 1" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' `concurrency` 4
< error
| 1064 42000 …line 1 column 41 near "`concurrency` 4" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' concurrency
< error
| 1064 42000 …line 1 column 39 near "" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' 4
< error
| 1064 42000 …line 1 column 29 near "4" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' select 1
< error
| 1064 42000 …line 1 column 34 near "select 1" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' concurrency 0
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| BACKUP TABLE a TO 'noop://' concurrency 0
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' concurrency 18446744073709551615
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| BACKUP TABLE a TO 'noop://' concurrency 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' concurrency 18446744073709551616
< error
| 1064 42000 …line 1 column 60 near "18446744073709551616" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' concurrency -1
< error
| 1064 42000 …line 1 column 41 near "-1" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' concurrency 'x'
< error
| 1064 42000 …line 1 column 43 near "'x'" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' concurrency 1.5
< error
| 1064 42000 …line 1 column 43 near "1.5" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' online 1
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| BACKUP TABLE a TO 'noop://' online 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' online true
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| BACKUP TABLE a TO 'noop://' online true
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' online FALSE
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| BACKUP TABLE a TO 'noop://' online FALSE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' online = TRUE
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| BACKUP TABLE a TO 'noop://' online = TRUE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' online 2
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| BACKUP TABLE a TO 'noop://' online 2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' online `true`
< error
| 1064 42000 …line 1 column 41 near "`true`" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' online 'x'
< error
| 1064 42000 …line 1 column 38 near "'x'" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' online -0
< error
| 1064 42000 …line 1 column 36 near "-0" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' send_credentials_to_tikv = 1
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| BACKUP TABLE a TO 'noop://' send_credentials_to_tikv = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' checkpoint 1 skip_schema_files 1 strict_format 1
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| BACKUP TABLE a TO 'noop://' checkpoint 1 skip_schema_files 1 strict_format 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' csv_not_null 1 csv_backslash_escape 1 csv_trim_last_separators 1
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| BACKUP TABLE a TO 'noop://' csv_not_null 1 csv_backslash_escape 1 csv_trim_last_separators 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' load_stats 0
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| BACKUP TABLE a TO 'noop://' load_stats 0
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' compression_type 'lz4'
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| BACKUP TABLE a TO 'noop://' compression_type 'lz4'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' compression_type lz4
< error
| 1064 42000 …line 1 column 48 near "lz4" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' tikv_importer = 'x'
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| BACKUP TABLE a TO 'noop://' tikv_importer = 'x'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' encryption_method 'aes128-ctr' encryption_keyfile 'f'
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| BACKUP TABLE a TO 'noop://' encryption_method 'aes128-ctr' encryption_keyfile 'f'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' csv_separator ',' csv_delimiter '"' csv_null 'N'
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| BACKUP TABLE a TO 'noop://' csv_separator ',' csv_delimiter '"' csv_null 'N'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' gc_ttl 'x' until_ts 'x' start_ts 'x' restored_ts 'x'
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| BACKUP TABLE a TO 'noop://' gc_ttl 'x' until_ts 'x' start_ts 'x' restored_ts 'x'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' full_backup_storage 'x'
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| BACKUP TABLE a TO 'noop://' full_backup_storage 'x'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' compression_type 1
< error
| 1064 42000 …line 1 column 46 near "1" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' backend = 'local'
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| BACKUP TABLE a TO 'noop://' backend = 'local'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' backend = local
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| BACKUP TABLE a TO 'noop://' backend = local
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' backend local
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| BACKUP TABLE a TO 'noop://' backend local
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' backend = `ignore`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| BACKUP TABLE a TO 'noop://' backend = `ignore`
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' backend = ignore
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| BACKUP TABLE a TO 'noop://' backend = ignore
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' backend = replace
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| BACKUP TABLE a TO 'noop://' backend = replace
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' `backend` = 'local'
< error
| 1064 42000 …line 1 column 37 near "`backend` = 'local'" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' backend = 1
< error
| 1064 42000 …line 1 column 39 near "1" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' on_duplicate = ignore
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| BACKUP TABLE a TO 'noop://' on_duplicate = ignore
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' ON DUPLICATE = ignore
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| BACKUP TABLE a TO 'noop://' ON DUPLICATE = ignore
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' ON `DUPLICATE` = ignore
< error
| 1064 42000 …line 1 column 42 near "`DUPLICATE` = ignore" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' ON = ignore
< error
| 1064 42000 …line 1 column 32 near "= ignore" 
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' snapshot = '2020-03-18 18:13:54'
< ok 0 0
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' snapshot = 1234567890
< ok 0 0
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' snapshot 5 minute ago
< ok 0 0
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' snapshot = 5 sql_tsi_hour ago
< ok 0 0
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' snapshot 5 microsecond ago
< ok 0 0
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' snapshot 5 week ago
< ok 0 0
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' snapshot 5 month ago
< error
| 1064 42000 …MONTH is not a constant time interval and cannot be used here
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' snapshot 5 year ago
< error
| 1064 42000 …YEAR is not a constant time interval and cannot be used here
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' snapshot 5 sql_tsi_month ago
< error
| 1064 42000 …MONTH is not a constant time interval and cannot be used here
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' snapshot 5 minute
< error
| 1064 42000 …line 1 column 48 near "" 
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' snapshot 5 `minute` ago
< error
| 1064 42000 …line 1 column 50 near "`minute` ago" 
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' snapshot 5 minute `ago`
< error
| 1064 42000 …line 1 column 54 near "`ago`" 
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' snapshot ago
< error
| 1064 42000 …line 1 column 43 near "ago" 
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' snapshot 18446744073709551616
< error
| 1064 42000 …line 1 column 60 near "18446744073709551616" 
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' SNAPSHOT 1234567890 CONCURRENCY 4
< ok 0 0
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' CONCURRENCY 4 SNAPSHOT 1234567890
< ok 0 0
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' SNAPSHOT 1234567890 ONLINE 1
< ok 0 0
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' SNAPSHOT 5 CONCURRENCY 4
< ok 0 0
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' SNAPSHOT 5 SNAPSHOT 6
< ok 0 0
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' SNAPSHOT = 5 minute ago CONCURRENCY 4
< ok 0 0
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' SNAPSHOT 5 day ago ONLINE 1 CONCURRENCY 4 checksum off
< ok 0 0
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' SNAPSHOT 1234567890 week ago
< ok 0 0
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' SNAPSHOT 5 nosuchopt 4
< error
| 1064 42000 …line 1 column 51 near "nosuchopt 4" 
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' SNAPSHOT 5 microsecond CONCURRENCY 4
< error
| 1064 42000 …line 1 column 65 near "CONCURRENCY 4" 
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' LAST_BACKUP 5 CONCURRENCY 4
< ok 0 0
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' csv_header 1 CONCURRENCY 4
< ok 0 0
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' csv_header fields CONCURRENCY 4
< ok 0 0
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' checksum 1 CONCURRENCY 4
< ok 0 0
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' checksum off CONCURRENCY 4
< ok 0 0
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' online 1 CONCURRENCY 4
< ok 0 0
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' rate_limit 500 MB/second CONCURRENCY 4
< ok 0 0
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' concurrency 4 backend local
< ok 0 0
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' backend local CONCURRENCY 4
< ok 0 0
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' backend CONCURRENCY 4
< error
| 1064 42000 …line 1 column 52 near "4" 
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' LAST_BACKUP = '2020-02-02 14:14:14'
< ok 0 0
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' LAST_BACKUP = 1234567890
< ok 0 0
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' LAST_BACKUP 1234567890
< ok 0 0
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' LAST_BACKUP 5 minute ago
< error
| 1064 42000 …line 1 column 51 near "minute ago" 
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' LAST_BACKUP
< error
| 1064 42000 …line 1 column 42 near "" 
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' rate_limit 500 MB/second
< ok 0 0
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' rate_limit = 500 MB / second
< ok 0 0
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' rate_limit 500 mb/second
< ok 0 0
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' rate_limit 500 MB
< error
| 1064 42000 …line 1 column 48 near "" 
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' rate_limit 500 GB/second
< error
| 1064 42000 …line 1 column 48 near "GB/second" 
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' rate_limit 500 MB/minute
< error
| 1064 42000 …line 1 column 55 near "minute" 
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' rate_limit 500 `MB`/second
< error
| 1064 42000 …line 1 column 50 near "`MB`/second" 
> exec ^
: string
| BACKUP DATABASE * TO 'noop://' rate_limit 500
< error
| 1064 42000 …line 1 column 45 near "" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' csv_header = 1
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| BACKUP TABLE a TO 'noop://' csv_header = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' csv_header = fields
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| BACKUP TABLE a TO 'noop://' csv_header = fields
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' csv_header COLUMNS
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| BACKUP TABLE a TO 'noop://' csv_header COLUMNS
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' csv_header = `fields`
< error
| 1064 42000 …line 1 column 49 near "`fields`" 
> exec ^
: string
| BACKUP TABLE a TO 'noop://' csv_header = 'fields'
< error
| 1064 42000 …line 1 column 49 near "'fields'" 
> exec ^
: string
| RESTORE TABLE g FROM 'noop://' checksum off
< ok 0 0
> exec ^
: string
| RESTORE TABLE g FROM 'noop://' checksum optional
< ok 0 0
> exec ^
: string
| RESTORE TABLE g FROM 'noop://' checksum required
< ok 0 0
> exec ^
: string
| RESTORE TABLE g FROM 'noop://' checksum 0
< ok 0 0
> exec ^
: string
| RESTORE TABLE g FROM 'noop://' checksum true
< ok 0 0
> exec ^
: string
| RESTORE TABLE g FROM 'noop://' analyze off
< ok 0 0
> exec ^
: string
| RESTORE TABLE g FROM 'noop://' analyze 1
< ok 0 0
> exec ^
: string
| RESTORE TABLE g FROM 'noop://' checksum `off`
< error
| 1064 42000 …line 1 column 45 near "`off`" 
> exec ^
: string
| RESTORE TABLE g FROM 'noop://' checksum 'off'
< error
| 1064 42000 …line 1 column 45 near "'off'" 
> exec ^
: string
| RESTORE TABLE g FROM 'noop://' checksum maybe
< error
| 1064 42000 …line 1 column 45 near "maybe" 
> exec ^
: string
| backup database * to 's3://bucket/path/?provider=alibaba&region=us-west-9&storage-class=glacier&sse=AES256&acl=authenticated-read&use-accelerate-endpoint=1' send_credentials_to_tikv = 1
< ok 0 0
> exec ^
: string
| backup table x to 's3://bucket/path/?endpoint=https://test-cluster-s3.local&access-key=aaaaaaaaa&secret-access-key=bbbbbbbb&force-path-style=1'
< error
| 1146 42S02 Table 'test.x' doesn't exist
> exec ^
: string
| backup table x to 's3://bucket/path/?endpoint=https://test-cluster-s3.local&access-key=aaaaaaaaa&secret-access-key=bbbbbbbb&force-path-style=1'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| restore database * from 'gcs://bucket/path/?endpoint=https://test-cluster.gcs.local&storage-class=coldline&predefined-acl=OWNER&credentials-file=/data/private/creds.json'
< ok 0 0
> exec ^
: string
| RECOVER TABLE BY JOB 11
< ok 0 0
> exec ^
: string
| recover table by job 0
< ok 0 0
> exec ^
: string
| RECOVER TABLE BY JOB 9223372036854775807
< ok 0 0
> exec ^
: string
| RECOVER TABLE BY JOB
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| RECOVER TABLE BY JOB -1
< error
| 1064 42000 …line 1 column 22 near "-1" 
> exec ^
: string
| RECOVER TABLE BY JOB 9223372036854775808
< error
| 1064 42000 …line 1 column 40 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| RECOVER TABLE BY JOB 11,12,13
< error
| 1064 42000 …line 1 column 24 near ",12,13" 
> exec ^
: string
| RECOVER TABLE BY JOB 'x'
< error
| 1064 42000 …line 1 column 24 near "'x'" 
> exec ^
: string
| RECOVER TABLE BY JOB t1
< error
| 1064 42000 …line 1 column 23 near "t1" 
> exec ^
: string
| RECOVER TABLE BY 11
< error
| 1064 42000 …line 1 column 19 near "11" 
> exec ^
: string
| RECOVER TABLE t1
< ok 0 0
> exec ^
: string
| RECOVER TABLE t1 100
< ok 0 0
> exec ^
: string
| RECOVER TABLE t1 0
< ok 0 0
> exec ^
: string
| RECOVER TABLE a.b
< ok 0 0
> exec ^
: string
| RECOVER TABLE a.b 100
< ok 0 0
> exec ^
: string
| RECOVER TABLE *.a
< ok 0 0
> exec ^
: string
| RECOVER TABLE t1,t2
< error
| 1064 42000 …line 1 column 17 near ",t2" 
> exec ^
: string
| RECOVER TABLE t1 100 200
< error
| 1064 42000 …line 1 column 24 near "200" 
> exec ^
: string
| RECOVER TABLE t1 abc
< error
| 1064 42000 …line 1 column 20 near "abc" 
> exec ^
: string
| RECOVER TABLE t1 -1
< error
| 1064 42000 …line 1 column 18 near "-1" 
> exec ^
: string
| RECOVER TABLE
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| RECOVER
< error
| 1064 42000 …line 1 column 7 near "" 
> exec ^
: string
| RECOVER DATABASE d
< error
| 1064 42000 …line 1 column 16 near "DATABASE d" 
> exec ^
: string
| RECOVER TABLE t1 'x'
< error
| 1064 42000 …line 1 column 20 near "'x'" 
> exec ^
: string
| flashback table t
< ok 0 0
> exec ^
: string
| flashback table t TO t1
< ok 0 0
> exec ^
: string
| flashback table a.b to c
< ok 0 0
> exec ^
: string
| flashback table `t` to `t1`
< ok 0 0
> exec ^
: string
| flashback database db1
< ok 0 0
> exec ^
: string
| flashback schema db1
< ok 0 0
> exec ^
: string
| flashback database db1 to db2
< ok 0 0
> exec ^
: string
| flashback schema db1 to db2
< ok 0 0
> exec ^
: string
| flashback table a.b to c.d
< error
| 1064 42000 …line 1 column 25 near ".d" 
> exec ^
: string
| flashback database a.b
< error
| 1064 42000 …line 1 column 21 near ".b" 
> exec ^
: string
| flashback schema d to e.f
< error
| 1064 42000 …line 1 column 24 near ".f" 
> exec ^
: string
| flashback table t,t1
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| flashback table t,t1 to t2
< error
| 1064 42000 …line 1 column 23 near "to t2" 
> exec ^
: string
| flashback table t to
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| flashback table
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| flashback database
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| flashback
< error
| 1064 42000 …line 1 column 9 near "" 
> exec ^
: string
| flashback table t to 'x'
< error
| 1064 42000 …line 1 column 24 near "'x'" 
> exec ^
: string
| flashback table t to 1
< error
| 1064 42000 …line 1 column 22 near "1" 
> exec ^
: string
| flashback cluster to timestamp '2021-05-26 16:45:26'
< ok 0 0
> exec ^
: string
| flashback cluster to tso 445494955052105721
< ok 0 0
> exec ^
: string
| flashback table t to timestamp '2021-05-26 16:45:26'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| flashback table t to timestamp '2021-05-26 16:45:26'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| flashback table t,t1 to timestamp '2021-05-26 16:45:26'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| flashback table t,t1 to timestamp '2021-05-26 16:45:26'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| flashback table t to tso 445494955052105722
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| flashback table t to tso 445494955052105722
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| flashback table t,t1 to tso 445494955052105723
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| flashback table t,t1 to tso 445494955052105723
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| flashback database test to timestamp '2021-05-26 16:45:26'
< ok 0 0
> exec ^
: string
| flashback schema test to timestamp '2021-05-26 16:45:26'
< ok 0 0
> exec ^
: string
| flashback database test to tso 445494955052105724
< ok 0 0
> exec ^
: string
| flashback schema test to tso 445494955052105725
< ok 0 0
> exec ^
: string
| flashback cluster to timestamp "2021-05-26 16:45:26"
< ok 0 0
> exec ^
: string
| flashback table t TO timestamp
< ok 0 0
> exec ^
: string
| flashback table t TO tso
< ok 0 0
> exec ^
: string
| flashback table t TO `timestamp`
< ok 0 0
> exec ^
: string
| flashback database d TO timestamp
< ok 0 0
> exec ^
: string
| flashback database d TO tso
< ok 0 0
> exec ^
: string
| flashback table t to timestamp 'a' 'b'
< error
| 1064 42000 …line 1 column 38 near "'b'" 
> exec ^
: string
| flashback table t to TIMESTAMP
< ok 0 0
> exec ^
: string
| flashback table to timestamp '2021-05-26 16:45:26'
< error
| 1064 42000 …line 1 column 28 near "timestamp '2021-05-26 16:45:26'" 
> exec ^
: string
| flashback database to timestamp '2021-05-26 16:45:26'
< error
| 1064 42000 …line 1 column 31 near "timestamp '2021-05-26 16:45:26'" 
> exec ^
: string
| flashback table to tso 445494955052105726
< error
| 1064 42000 …line 1 column 22 near "tso 445494955052105726" 
> exec ^
: string
| flashback database to tso 445494955052105727
< error
| 1064 42000 …line 1 column 25 near "tso 445494955052105727" 
> exec ^
: string
| flashback cluster
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| flashback cluster to timestamp
< error
| 1064 42000 …line 1 column 20 near "to timestamp" 
> exec ^
: string
| flashback cluster to tso
< error
| 1064 42000 …line 1 column 20 near "to tso" 
> exec ^
: string
| flashback cluster to c
< error
| 1064 42000 …line 1 column 20 near "to c" 
> exec ^
: string
| flashback cluster to timestamp 'x' to tso 1
< error
| 1064 42000 …line 1 column 41 near "tso 1" 
> exec ^
: string
| flashback cluster to timestamp DATE_SUB(NOW(), INTERVAL 3 SECOND)
< error
| 1064 42000 …line 1 column 20 near "to timestamp DATE_SUB(NOW(), INTERVAL 3 SECOND)" 
> exec ^
: string
| flashback cluster to timestamp TIDB_BOUNDED_STALENESS(DATE_SUB(NOW(), INTERVAL 3 SECOND), NOW())
< error
| 1064 42000 …line 1 column 20 near "to timestamp TIDB_BOUNDED_STALENESS(DATE_SUB(NOW(), INTERVAL 3 SECOND), NOW())" 
> exec ^
: string
| flashback `cluster` to timestamp '2021-05-26 16:45:26'
< error
| 1064 42000 …line 1 column 19 near "`cluster` to timestamp '2021-05-26 16:45:26'" 
> exec ^
: string
| flashback schema test to tso 18446744073709551615
< ok 0 0
> exec ^
: string
| flashback schema test to tso 0
< error
| 1064 42000 …line 1 column 30 near ""Invalid TSO value provided: 0 
> exec ^
: string
| flashback schema test to tso -100
< error
| 1064 42000 …line 1 column 30 near "-100" 
> exec ^
: string
| flashback schema test to tso 18446744073709551616
< error
| 1064 42000 …line 1 column 49 near "18446744073709551616" 
> exec ^
: string
| flashback table t to tso 1.5
< error
| 1064 42000 …line 1 column 28 near "1.5" 
> exec ^
: string
| flashback table t to tso +1
< error
| 1064 42000 …line 1 column 26 near "+1" 
> exec ^
: string
| CREATE TABLE flashback (flashback int)
< ok 0 0
> exec ^
: string
| SELECT flashback FROM flashback
< error
| 1146 42S02 Table 'test.flashback' doesn't exist
> exec ^
: string
| SELECT flashback FROM flashback
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t flashback
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t flashback
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t AS flashback
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t AS flashback
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t flashback WHERE flashback.a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t flashback WHERE flashback.a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| UPDATE t SET flashback = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| UPDATE t SET flashback = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX flashback ON t (a)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX flashback ON t (a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE DATABASE flashback
< ok 0 0
> exec ^
: string
| SELECT 1 AS flashback
< ok 0 0
> exec ^
: string
| SELECT * FROM flashback.t
< error
| 1146 42S02 Table 'flashback.t' doesn't exist
> exec ^
: string
| SELECT * FROM flashback.t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (backup int, restore int, recover int)
< ok 0 0
> exec ^
: string
| SELECT * FROM t backup
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t backup
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t restore
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t restore
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t recover
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t recover
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t backup WHERE backup.a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t backup WHERE backup.a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
