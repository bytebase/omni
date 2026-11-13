-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestOutfileFieldsFormsMatchTiDB default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| select * into outfile '/x' fields terminated by 'a'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| select * into outfile '/x' fields terminated by 'a'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * into outfile '/x' fields terminated by 'a' terminated by 'b'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| select * into outfile '/x' fields terminated by 'a' terminated by 'b'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * into outfile '/x' fields optionally enclosed by 'a' optionally enclosed by 'b'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| select * into outfile '/x' fields optionally enclosed by 'a' optionally enclosed by 'b'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * into outfile '/x' fields terminated by 'a' enclosed by 'b' escaped by 'c'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| select * into outfile '/x' fields terminated by 'a' enclosed by 'b' escaped by 'c'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * into outfile '/x' fields enclosed by 'a' optionally enclosed by 'b'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| select * into outfile '/x' fields enclosed by 'a' optionally enclosed by 'b'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * into outfile '/x' fields
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| select * into outfile '/x' fields lines terminated by 'a'
< error
| 1064 42000 …line 1 column 39 near "lines terminated by 'a'" 
> exec ^
: string
| select * into outfile '/x' fields optionally
< error
| 1064 42000 …line 1 column 44 near "optionally" 
> exec ^
: string
| select * into outfile '/x' fields optionally enclosed
< error
| 1064 42000 …line 1 column 44 near "optionally enclosed" 
> exec ^
: string
| select * into outfile '/x' fields optionally enclosed by 'ab'
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| select * into outfile '/x' fields optionally enclosed by 'ab'
| ; SELECT FROM
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| select * into outfile '/x' fields defined null by 'n'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| select * into outfile '/x' fields defined null by 'n'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * into outfile '/x' fields defined null by 'n' optionally enclosed
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| select * into outfile '/x' fields defined null by 'n' optionally enclosed
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * into outfile '/x' fields defined null by 'a' defined null by 'b'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| select * into outfile '/x' fields defined null by 'a' defined null by 'b'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * into outfile '/x' fields terminated by 'a' defined null by 'b' enclosed by 'c'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| select * into outfile '/x' fields terminated by 'a' defined null by 'b' enclosed by 'c'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * into outfile '/x' fields defined null by 'a' optionally enclosed by 'b'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| select * into outfile '/x' fields defined null by 'a' optionally enclosed by 'b'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * into outfile '/x' fields defined null by 'a' optionally enclosed by 'ab'
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| select * into outfile '/x' fields defined null by 'a' optionally enclosed by 'ab'
| ; SELECT FROM
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| select * into outfile '/x' fields defined null by 'a' optionally
< error
| 1064 42000 …line 1 column 64 near "" 
> exec ^
: string
| select * into outfile '/x' fields `defined` null by 'a'
< error
| 1064 42000 …line 1 column 43 near "`defined` null by 'a'" 
> exec ^
: string
| select * into outfile
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| select * into outfile 's' lines starting 's'
< error
| 1064 42000 …line 1 column 44 near "'s'" 
> exec ^
: string
| select * into outfile 's' lines starting by
< error
| 1064 42000 …line 1 column 43 near "" 
> exec ^
: string
| select * into outfile 's' lines terminated 's'
< error
| 1064 42000 …line 1 column 46 near "'s'" 
> exec ^
: string
| select * into outfile 's' lines terminated by
< error
| 1064 42000 …line 1 column 45 near "" 
> exec ^
: string
| select * into outfile 's' lines
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| select * into outfile 's' lines
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * into outfile 's' lines starting by 'a'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| select * into outfile 's' lines starting by 'a'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * into outfile 's' lines terminated by 'b'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| select * into outfile 's' lines terminated by 'b'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * into outfile 's' lines starting by 'a' terminated by 'b'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| select * into outfile 's' lines starting by 'a' terminated by 'b'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
