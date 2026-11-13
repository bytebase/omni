-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestIntervalSignedValue default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| SELECT DATE_ADD('2007-03-28 22:08:28', INTERVAL -2.2 MINUTE_MICROSECOND)
< ok 0 0
> exec ^
: string
| SELECT DATE_ADD('2007-03-28 22:08:28', INTERVAL -2.2 SECOND)
< ok 0 0
> exec ^
: string
| SELECT DATE_ADD('2007-03-28 22:08:28', INTERVAL -2 YEAR_MONTH)
< ok 0 0
> exec ^
: string
| SELECT DATE_ADD('2007-03-28 22:08:28', INTERVAL -2 DAY_SECOND)
< ok 0 0
> exec ^
: string
| SELECT '1900-01-01 00:00:00' + INTERVAL -1 DAY
< ok 0 0
> exec ^
: string
| SELECT '1900-01-01 00:00:00' + INTERVAL +1 HOUR
< ok 0 0
> exec ^
: string
| SELECT DATE_ADD('x', INTERVAL 2.2 MINUTE_MICROSECOND)
< ok 0 0
> exec ^
: string
| SELECT '1900-01-01' + INTERVAL 1 DAY
< ok 0 0
