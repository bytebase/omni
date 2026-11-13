-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestLoadDataFieldFormsOracle default#1
> exec
| load data infile '/x.csv' into table t fields terminated by X'6B6B'
< error
| 1017 HY000 File not exist
> exec
| load data infile '/x.csv' into table t fields terminated by x'6b6b'
< error
| 1017 HY000 File not exist
> exec
| load data infile '/x.csv' into table t fields terminated by 0x6B6B
< error
| 1017 HY000 File not exist
> exec
| load data infile '/x.csv' into table t fields terminated by B'110101101101011'
< error
| 1017 HY000 File not exist
> exec
| load data infile '/x.csv' into table t fields terminated by b'1101'
< error
| 1017 HY000 File not exist
> exec
| load data infile '/x.csv' into table t fields terminated by X''
< error
| 1017 HY000 File not exist
> exec
| load data infile '/x.csv' into table t fields enclosed by X'0D'
< error
| 1017 HY000 File not exist
> exec
| load data infile '/x.csv' into table t fields escaped by X'5C'
< error
| 1017 HY000 File not exist
> exec
| load data infile '/x.csv' into table t fields optionally enclosed by X'22'
< error
| 1017 HY000 File not exist
> exec
| load data infile '/x.csv' into table t columns terminated by X'09'
< error
| 1017 HY000 File not exist
> exec
| load data infile '/x.csv' into table t lines starting by X'6B6B' terminated by X'6B6B'
< error
| 1017 HY000 File not exist
> exec
| load data infile '/x.csv' into table t lines starting by B'1101' terminated by B'1101'
< error
| 1017 HY000 File not exist
> exec
| load data infile '/x.csv' into table t character set utf8mb4 fields terminated by X'09'
< error
| 1017 HY000 File not exist
> exec
| load data infile '/x.csv' replace into table t fields terminated by X'09'
< error
| 1017 HY000 File not exist
> exec
| load data infile '/x.csv' into table t fields terminated by X'616'
< error
| 1064 42000 …near 'X'616'' at line 1
> exec
| load data infile '/x.csv' into table t character set binary
< error
| 1017 HY000 File not exist
> exec
| load data infile '/x.csv' into table t character set 'utf8mb4'
< error
| 1017 HY000 File not exist
> exec
| load data infile '/x.csv' into table t character set 1
< error
| 1064 42000 …near '1' at line 1
> exec
| load data infile '/x.csv' into table t fields terminated by 'a' terminated by 'b'
< error
| 1017 HY000 File not exist
> exec
| load data infile '/x.csv' into table t fields escaped by 'a' terminated by 'b'
< error
| 1017 HY000 File not exist
> exec
| load data infile '/x.csv' into table t fields enclosed by 'a' escaped by 'b' terminated by 'c'
< error
| 1017 HY000 File not exist
> exec
| load data infile '/x.csv' into table t lines terminated by 'b' starting by 'c'
< error
| 1017 HY000 File not exist
> exec
| load data infile '/x.csv' into table t lines starting by 'a' starting by 'b'
< error
| 1017 HY000 File not exist
> exec
| load data infile '/x.csv' into table t fields optionally enclosed by 'a' enclosed by 'b'
< error
| 1017 HY000 File not exist
> exec
| load data infile '/x.csv' into table t fields enclosed by 'a' optionally enclosed by 'b'
< error
| 1017 HY000 File not exist
> exec
| load data infile '/x.csv' into table t fields optionally enclosed by '"' terminated by ','
< error
| 1017 HY000 File not exist
> exec
| load data infile '/x.csv' into table t fields
< error
| 1064 42000 …near 'fields' at line 1
> exec
| load data infile '/x.csv' into table t lines
< error
| 1064 42000 …near 'lines' at line 1
> exec
| load data infile '/x.csv' into table t fields lines terminated by 'a'
< error
| 1064 42000 …near 'lines terminated by 'a'' at line 1
> exec
| load data infile '/x.csv' into table t fields optionally enclosed ignore 1 lines
< error
| 1064 42000 …near 'ignore 1 lines' at line 1
> exec
| load data infile '/x.csv' into table t fields optionally `enclosed` by 'a'
< error
| 1064 42000 …near '`enclosed` by 'a'' at line 1
> exec
| load data infile '/x.csv' into table t (a, @v) set b = @v
< error
| 1017 HY000 File not exist
> exec
| load data infile '/x.csv' into table t (`@v`)
< error
| 1017 HY000 File not exist
> exec
| load data infile '/x.csv' into table t (`@v`, a)
< error
| 1017 HY000 File not exist
> exec
| load data infile '/x.csv' into table t fields defined null by 'nil'
< error
| 1064 42000 …near 'defined null by 'nil'' at line 1
> exec
| load data infile '/x.csv' into table t fields defined null by 'NULL' optionally enclosed
< error
| 1064 42000 …near 'defined null by 'NULL' optionally enclosed' at line 1
> exec
| load data infile '/x.csv' format 'csv' into table t
< error
| 1064 42000 …near 'format 'csv' into table t' at line 1
> exec
| load data infile '/x.csv' into table t with detached
< error
| 1064 42000 …near 'with detached' at line 1
> exec
| load data infile '/x.csv' into table t with threads=10
< error
| 1064 42000 …near 'with threads=10' at line 1
> exec
| import into t from '/file.csv'
< error
| 1064 42000 …near 'import into t from '/file.csv'' at line 1
