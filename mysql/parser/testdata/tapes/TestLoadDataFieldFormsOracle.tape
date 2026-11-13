-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestLoadDataFieldFormsOracle default#1
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| load data infile '/x.csv' into table t fields terminated by X'6B6B'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t fields terminated by X'6B6B'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| load data infile '/x.csv' into table t character set utf8mb4 fields terminated by X'09'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t character set utf8mb4 fields terminated by X'09'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| load data infile '/x.csv' replace into table t fields terminated by X'09'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' replace into table t fields terminated by X'09'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| load data infile '/x.csv' into table t fields terminated by X'616'
< error
| 1064 42000 …near 'X'616'' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| load data infile '/x.csv' into table t charset binary
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t charset binary
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| load data infile '/x.csv' into table t charset 'utf8mb4'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t charset 'utf8mb4'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| load data infile '/x.csv' into table t character set 1
< error
| 1064 42000 …near '1' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| load data infile '/x.csv' into table t lines terminated by 'b' starting by 'c'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t lines terminated by 'b' starting by 'c'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| load data infile '/x.csv' into table t lines starting by 'a' starting by 'b'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t lines starting by 'a' starting by 'b'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| load data infile '/x.csv' into table t fields optionally enclosed by '"' terminated by ','
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x.csv' into table t fields optionally enclosed by '"' terminated by ','
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| load data infile '/x.csv' into table t fields
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| load data infile '/x.csv' into table t lines
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| load data infile '/x.csv' into table t fields lines terminated by 'a'
< error
| 1064 42000 …near 'lines terminated by 'a'' at line 1
> exec ^
: string
| load data infile '/x.csv' into table t fields optionally enclosed ignore 1 lines
< error
| 1064 42000 …near 'ignore 1 lines' at line 1
> exec ^
: string
| load data infile '/x.csv' into table t fields optionally `enclosed` by 'a'
< error
| 1064 42000 …near '`enclosed` by 'a'' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| load data infile '/x.csv' into table t fields defined null by 'nil'
< error
| 1064 42000 …near 'defined null by 'nil'' at line 1
> exec ^
: string
| load data infile '/x.csv' into table t fields defined null by 'NULL' optionally enclosed
< error
| 1064 42000 …near 'defined null by 'NULL' optionally enclosed' at line 1
> exec ^
: string
| load data infile '/x.csv' format 'csv' into table t
< error
| 1064 42000 …near 'format 'csv' into table t' at line 1
> exec ^
: string
| load data infile '/x.csv' into table t with detached
< error
| 1064 42000 …near 'with detached' at line 1
> exec ^
: string
| load data infile '/x.csv' into table t with threads=10
< error
| 1064 42000 …near 'with threads=10' at line 1
> exec ^
: string
| import into t from '/file.csv'
< error
| 1064 42000 …near 'into t from '/file.csv'' at line 1
