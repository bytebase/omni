-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestIntroducerLiteralFormsTiDBOracle default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| select _utf8 0xD0B1
< ok 0 0
> exec ^
: string
| select _utf8 X'D0B1'
< ok 0 0
> exec ^
: string
| select _utf8 0b1101000010110001
< ok 0 0
> exec ^
: string
| select _utf8 B'1101000010110001'
< ok 0 0
> exec ^
: string
| select _latin1 0x61
< ok 0 0
> exec ^
: string
| select _utf8mb4 x'61'
< ok 0 0
> exec ^
: string
| select _utf8 0x61 + 1
< ok 0 0
> exec ^
: string
| select _utf8 0xD0B1 from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select _utf8 0xD0B1 from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select _notacharset 0x61
< error
| 1064 42000 …line 1 column 24 near "0x61" 
> exec ^
: string
| select `_utf8` 0x61
< error
| 1064 42000 …line 1 column 19 near "0x61" 
> exec ^
: string
| select _utf8 'str'
< ok 0 0
