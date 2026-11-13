-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestLoadDataFormsMatchTiDB default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| load data infile '/x.csv' into table t
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile into table t
< error
| 1064 42000 …line 1 column 21 near "into table t" 
> exec ^
: string
| load data infile '/x.csv' table t
< error
| 1064 42000 …line 1 column 31 near "table t" 
> exec ^
: string
| load data infile '/x.csv' into t
< error
| 1064 42000 …line 1 column 32 near "t" 
> exec ^
: string
| load data infile '/x.csv' into table t fields terminated by X'6B6B'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t fields terminated by X'6B6B'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t fields terminated by x'6b6b'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t fields terminated by x'6b6b'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t fields terminated by 0x6B6B
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t fields terminated by 0x6B6B
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t fields terminated by B'110101101101011'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t fields terminated by B'110101101101011'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t fields terminated by b'1101'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t fields terminated by b'1101'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t fields terminated by X''
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t fields terminated by X''
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t fields enclosed by X'0D'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t fields enclosed by X'0D'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t fields enclosed by X'0D0D'
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| load data infile '/x.csv' into table t fields enclosed by X'0D0D'
| ; SELECT FROM
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| load data infile '/x.csv' into table t fields escaped by X'5C'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t fields escaped by X'5C'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t fields optionally enclosed by X'22'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t fields optionally enclosed by X'22'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t columns terminated by X'09'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t columns terminated by X'09'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t lines starting by X'6B6B' terminated by X'6B6B'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t lines starting by X'6B6B' terminated by X'6B6B'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t lines starting by B'1101' terminated by B'1101'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t lines starting by B'1101' terminated by B'1101'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t fields terminated by X'616'
< error
| 1064 42000 …line 1 column 66 near "X'616'"hex literal: invalid hexadecimal format, must even numbers, but 3 
> exec ^
: string
| load data infile '/x.csv' into table t fields terminated by 'a' terminated by 'b'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t fields terminated by 'a' terminated by 'b'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t fields escaped by 'a' terminated by 'b'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t fields escaped by 'a' terminated by 'b'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t fields enclosed by 'a' escaped by 'b' terminated by 'c'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t fields enclosed by 'a' escaped by 'b' terminated by 'c'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t fields
< error
| 1064 42000 …line 1 column 45 near "" 
> exec ^
: string
| load data infile '/x.csv' into table t fields lines terminated by 'a'
< error
| 1064 42000 …line 1 column 51 near "lines terminated by 'a'" 
> exec ^
: string
| load data infile '/x.csv' into table t fields optionally enclosed ignore 1 lines
< error
| 1064 42000 …line 1 column 56 near "optionally enclosed ignore 1 lines" 
> exec ^
: string
| load data infile '/x.csv' into table t fields optionally enclosed by '"'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t fields optionally enclosed by '"'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t fields optionally `enclosed` by 'a'
< error
| 1064 42000 …line 1 column 56 near "optionally `enclosed` by 'a'" 
> exec ^
: string
| load data infile '/x.csv' into table t fields defined null by 'nil'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t fields defined null by 'nil'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t fields defined null by X'00'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t fields defined null by X'00'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t fields defined null by B'1101'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t fields defined null by B'1101'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t fields defined null by 'a' defined null by 'b'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t fields defined null by 'a' defined null by 'b'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t fields defined null by 'a' optionally enclosed defined null by 'b'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t fields defined null by 'a' optionally enclosed defined null by 'b'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t fields optionally enclosed by 'a' enclosed by 'b'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t fields optionally enclosed by 'a' enclosed by 'b'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t fields enclosed by 'a' optionally enclosed by 'b'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t fields enclosed by 'a' optionally enclosed by 'b'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t fields terminated by 'a' defined null by 'b' enclosed by 'c'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t fields terminated by 'a' defined null by 'b' enclosed by 'c'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t fields defined null by 'a' terminated by 'b'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t fields defined null by 'a' terminated by 'b'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t fields defined null by 'NULL' optionally enclosed
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t fields defined null by 'NULL' optionally enclosed
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t fields defined null by 'NULL' optionally enclosed ignore 1 lines
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t fields defined null by 'NULL' optionally enclosed ignore 1 lines
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t fields defined null by 'a' optionally enclosed by 'b'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t fields defined null by 'a' optionally enclosed by 'b'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t fields defined null by 'a' optionally enclosed by X'62'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t fields defined null by 'a' optionally enclosed by X'62'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t fields defined null by 'a' optionally enclosed by 'bc'
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| load data infile '/x.csv' into table t fields defined null by 'a' optionally enclosed by 'bc'
| ; SELECT FROM
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| load data infile '/x.csv' into table t fields defined null by 'a' optionally enclosed by X'6263'
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| load data infile '/x.csv' into table t fields defined null by 'a' optionally enclosed by X'6263'
| ; SELECT FROM
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| load data infile '/x.csv' into table t fields defined null by 'a' optionally
< error
| 1064 42000 …line 1 column 76 near "" 
> exec ^
: string
| load data infile '/x.csv' into table t fields `defined` null by 'a'
< error
| 1064 42000 …line 1 column 55 near "`defined` null by 'a'" 
> exec ^
: string
| load data infile '/x.csv' into table t fields defined `null` by 'a'
< error
| 1064 42000 …line 1 column 60 near "`null` by 'a'" 
> exec ^
: string
| load data infile '/x.csv' into table t fields defined null `by` 'a'
< error
| 1064 42000 …line 1 column 63 near "`by` 'a'" 
> exec ^
: string
| load data infile '/x.csv' into table t character set binary
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t character set binary
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t character set utf8mb4
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t character set utf8mb4
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t character set 'utf8mb4'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t character set 'utf8mb4'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t character set 1
< error
| 1064 42000 …line 1 column 54 near "1" 
> exec ^
: string
| load data infile '/x.csv' into table t character binary
< error
| 1064 42000 …line 1 column 55 near "binary" 
> exec ^
: string
| load data infile '/x.csv' into table t lines
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t lines
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t lines starting by 'a'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t lines starting by 'a'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t lines terminated by 'b'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t lines terminated by 'b'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t lines starting by 'a' terminated by 'b'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t lines starting by 'a' terminated by 'b'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t lines terminated by 'b' starting by 'c'
< error
| 1064 42000 …line 1 column 71 near "starting by 'c'" 
> exec ^
: string
| load data infile '/x.csv' into table t lines starting by 'a' starting by 'b'
< error
| 1064 42000 …line 1 column 69 near "starting by 'b'" 
> exec ^
: string
| load data infile '/x.csv' into table t lines starting by 'a' terminated by 'b' starting by 'c'
< error
| 1064 42000 …line 1 column 87 near "starting by 'c'" 
> exec ^
: string
| load data infile '/x.csv' into table t ignore 2 lines
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t ignore 2 lines
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t ignore 2 rows
< error
| 1064 42000 …line 1 column 52 near "rows" 
> exec ^
: string
| load data infile '/x.csv' into table t ignore lines
< error
| 1064 42000 …line 1 column 51 near "lines" 
> exec ^
: string
| load data infile '/x.csv' into table t ignore 2 lines fields terminated by ','
< error
| 1064 42000 …line 1 column 60 near "fields terminated by ','" 
> exec ^
: string
| load data infile '/x.csv' format 'delimited data' into table t
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' format 'delimited data' into table t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.sql' format 'sql file' replace into table t
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.sql' format 'sql file' replace into table t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' format 'csv' ignore into table t
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' format 'csv' ignore into table t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data local infile '/x.parquet' format 'parquet' into table t
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data local infile '/x.parquet' format 'parquet' into table t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t format 'csv'
< error
| 1064 42000 …line 1 column 45 near "format 'csv'" 
> exec ^
: string
| load data infile '/x.csv' format 'csv' format 'csv' into table t
< error
| 1064 42000 …line 1 column 45 near "format 'csv' into table t" 
> exec ^
: string
| load data infile '/x.csv' format X'6161' into table t
< error
| 1064 42000 …line 1 column 40 near "X'6161' into table t" 
> exec ^
: string
| load data infile '/x.csv' `format` 'csv' into table t
< error
| 1064 42000 …line 1 column 34 near "`format` 'csv' into table t" 
> exec ^
: string
| load data infile '/x.csv' into table t with detached
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t with detached
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t with threads=10
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t with threads=10
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t with threads=10, detached
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t with threads=10, detached
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t with threads=-1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t with threads=-1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t with threads=+2
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t with threads=+2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t with threads=1.5
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t with threads=1.5
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t with threads='abc'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t with threads='abc'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t with threads=X'3130'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t with threads=X'3130'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t with threads=true
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t with threads=true
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t with threads=null
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t with threads=null
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t with character_set='utf8'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t with character_set='utf8'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t with Threads=10
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t with Threads=10
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t with detached, detached
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t with detached, detached
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t with `detached`
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t with `detached`
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t with `threads`=10
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t with `threads`=10
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t fields terminated by 'a' with detached
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t fields terminated by 'a' with detached
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t (a) set b=1 with detached
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t (a) set b=1 with detached
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t with
< error
| 1064 42000 …line 1 column 43 near "" 
> exec ^
: string
| load data infile '/x.csv' into table t with detached,
< error
| 1064 42000 …line 1 column 53 near "" 
> exec ^
: string
| load data infile '/x.csv' into table t with threads=
< error
| 1064 42000 …line 1 column 52 near "" 
> exec ^
: string
| load data infile '/x.csv' into table t with threads=@a
< error
| 1064 42000 …line 1 column 54 near "@a" 
> exec ^
: string
| load data infile '/x.csv' into table t with select=1
< error
| 1064 42000 …line 1 column 50 near "select=1" 
> exec ^
: string
| load data infile '/x.csv' into table t with detached with detached
< error
| 1064 42000 …line 1 column 57 near "with detached" 
> exec ^
: string
| load data infile '/x.csv' into table t (a, @v) set b = @v
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t (a, @v) set b = @v
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t (`@v`)
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t (`@v`)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t (`@v`, a)
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t (`@v`, a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t (a,@1) with detached
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t (a,@1) with detached
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x.csv' into table t set b = default
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t set b = default
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t from '/file.csv'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t from '/file.csv'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into `t` from '/file.csv'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into `t` from '/file.csv'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into d.t from '/file.csv'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into d.t from '/file.csv'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t (a,b) from '/file.csv'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t (a,b) from '/file.csv'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t (a,@1) from '/file.csv'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t (a,@1) from '/file.csv'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t (@1) from '/file.csv'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t (@1) from '/file.csv'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t () from '/file.csv'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t () from '/file.csv'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t (a,@1) set b=@1+100 from '/file.csv'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t (a,@1) set b=@1+100 from '/file.csv'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t set a=1 from '/file.csv'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t set a=1 from '/file.csv'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t (a) set b=default from '/file.csv'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t (a) set b=default from '/file.csv'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t from '/file.csv' format 'sql file'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t from '/file.csv' format 'sql file'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t from '/file.csv' with detached
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t from '/file.csv' with detached
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t from '/file.csv' with thread=1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t from '/file.csv' with thread=1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t from '/file.csv' with detached, thread=1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t from '/file.csv' with detached, thread=1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t from '/file.csv' with detached=1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t from '/file.csv' with detached=1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t from '/file.csv' with thread
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t from '/file.csv' with thread
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t from '/file.csv' with __starting_ts=1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t from '/file.csv' with __starting_ts=1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t from '/file.csv' format 'csv' with detached
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t from '/file.csv' format 'csv' with detached
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t from '/file.csv' format 'csv' with cloud_storage_uri='s3://x'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t from '/file.csv' format 'csv' with cloud_storage_uri='s3://x'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t (a,b) set c=1 from '/file.csv' format 'csv' with thread=2, detached
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t (a,b) set c=1 from '/file.csv' format 'csv' with thread=2, detached
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| import into t from
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| import into t from '/file.csv' with detached format 'csv'
< error
| 1064 42000 …line 1 column 51 near "format 'csv'" 
> exec ^
: string
| import into t from 'f' format 'csv' format 'csv'
< error
| 1064 42000 …line 1 column 42 near "format 'csv'" 
> exec ^
: string
| import into t from '/file.csv' format X'6363'
< error
| 1064 42000 …line 1 column 45 near "X'6363'" 
> exec ^
: string
| import into t from '/file.csv' `format` 'csv'
< error
| 1064 42000 …line 1 column 39 near "`format` 'csv'" 
> exec ^
: string
| import into t from '/file.csv' with thread=1 ,
< error
| 1064 42000 …line 1 column 46 near "" 
> exec ^
: string
| import into t partition (p0) from '/file.csv'
< error
| 1064 42000 …line 1 column 23 near "partition (p0) from '/file.csv'" 
> exec ^
: string
| import into t from select * from s
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t from select * from s
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t from (select * from s)
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t from (select * from s)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t from with cte as (select 1) select * from cte
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t from with cte as (select 1) select * from cte
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t from table s
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t from table s
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t (a,b) from select * from s
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t (a,b) from select * from s
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t (a) from select a from s with thread=2
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t (a) from select a from s with thread=2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t (`@v`) from select 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t (`@v`) from select 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t from select a from s order by a with thread=2
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t from select a from s order by a with thread=2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t from select a from s limit 1 with thread=2
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t from select a from s limit 1 with thread=2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t from select a from s group by a with rollup
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t from select a from s group by a with rollup
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t from select a from s group by a with rollup with thread=2
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t from select a from s group by a with rollup with thread=2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t from select a from s group by a with thread=2
< error
| 1064 42000 …line 1 column 57 near "thread=2" 
> exec ^
: string
| import into t (a,@1) from select a from s
< error
| 1064 42000 …line 1 column 41 near ""Cannot use user variable(1) in IMPORT INTO FROM SELECT statement. 
> exec ^
: string
| import into t (a,@1) from table s
< error
| 1064 42000 …line 1 column 33 near ""Cannot use user variable(1) in IMPORT INTO FROM SELECT statement. 
> exec ^
: string
| import into t set a=1 from select * from s
< error
| 1064 42000 …line 1 column 42 near ""Cannot use SET clause in IMPORT INTO FROM SELECT statement. 
> exec ^
: string
| import into t (a) set b=1 from select * from s
< error
| 1064 42000 …line 1 column 46 near ""Cannot use SET clause in IMPORT INTO FROM SELECT statement. 
> exec ^
: string
| import into t from values (1)
< error
| 1064 42000 …line 1 column 27 near "(1)" 
> exec ^
: string
| import into t from values row(1)
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t from values row(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t from values row(1,1)
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t from values row(1,1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t from values row(1), row(2)
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t from values row(1), row(2)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t from values row(1) union select 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t from values row(1) union select 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t from table t1 union table t2
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t from table t1 union table t2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t from table t1 union all table t2
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t from table t1 union all table t2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t from table t1 intersect table t2
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t from table t1 intersect table t2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t from table t1 except table t2
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t from table t1 except table t2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t from table t1 order by a limit 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t from table t1 order by a limit 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t from table t1 union
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| import into t from table t1 union table
< error
| 1064 42000 …line 1 column 39 near "" 
> exec ^
: string
| import into t from table
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| import into t from values
< error
| 1064 42000 …line 1 column 25 near "" 
> exec ^
: string
| import table from '/tmp/a.sdi'
< error
| 1064 42000 …line 1 column 12 near "table from '/tmp/a.sdi'" 
