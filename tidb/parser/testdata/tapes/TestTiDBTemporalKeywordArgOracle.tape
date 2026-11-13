-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestTiDBTemporalKeywordArgOracle default#1
> exec
| CREATE DATABASE IF NOT EXISTS omni_tka
< ok 0 0
> exec
| USE omni_tka
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS tu (a DATETIME, b DATETIME)
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS tu2 (second INT, quarter INT)
< ok 0 0
> exec
| SELECT TIMESTAMPDIFF(MICROSECOND, a, b) FROM tu
< ok 0 0
> exec
| SELECT TIMESTAMPDIFF(SECOND, a, b) FROM tu
< ok 0 0
> exec
| SELECT TIMESTAMPDIFF(MINUTE, a, b) FROM tu
< ok 0 0
> exec
| SELECT TIMESTAMPDIFF(HOUR, a, b) FROM tu
< ok 0 0
> exec
| SELECT TIMESTAMPDIFF(DAY, a, b) FROM tu
< ok 0 0
> exec
| SELECT TIMESTAMPDIFF(WEEK, a, b) FROM tu
< ok 0 0
> exec
| SELECT TIMESTAMPDIFF(MONTH, a, b) FROM tu
< ok 0 0
> exec
| SELECT TIMESTAMPDIFF(QUARTER, a, b) FROM tu
< ok 0 0
> exec
| SELECT TIMESTAMPDIFF(YEAR, a, b) FROM tu
< ok 0 0
> exec
| SELECT TIMESTAMPADD(SECOND, 1, a) FROM tu
< ok 0 0
> exec
| SELECT TIMESTAMPADD(YEAR, 1, a) FROM tu
< ok 0 0
> exec
| SELECT TIMESTAMPDIFF(SQL_TSI_SECOND, a, b) FROM tu
< ok 0 0
> exec
| SELECT TIMESTAMPADD(SQL_TSI_YEAR, 1, a) FROM tu
< ok 0 0
> exec
| SELECT a + INTERVAL 1 SQL_TSI_DAY FROM tu
< ok 0 0
> exec
| SELECT a + INTERVAL 1 DAY FROM tu
< ok 0 0
> exec
| SELECT GET_FORMAT(DATE, 'USA') FROM tu
< ok 0 0
> exec
| SELECT GET_FORMAT(TIME, 'USA') FROM tu
< ok 0 0
> exec
| SELECT GET_FORMAT(DATETIME, 'USA') FROM tu
< ok 0 0
> exec
| SELECT GET_FORMAT(TIMESTAMP, 'USA') FROM tu
< ok 0 0
> exec
| SELECT second FROM tu2
< ok 0 0
> exec
| CREATE TABLE tu3 (second INT, quarter INT)
< ok 0 0
> exec
| SELECT s.timestampdiff(1, 2, 3)
< error
| 1305 42000 FUNCTION s.timestampdiff does not exist
> exec
| SELECT s.get_format(1, 2)
< error
| 1305 42000 FUNCTION s.get_format does not exist
> exec
| SELECT `timestampdiff`(1, 2, 3)
< ok 0 0
> exec
| SELECT `TIMESTAMPDIFF`(SECOND, a, b) FROM tu
< error
| 1054 42S22 Unknown column 'second' in 'field list'
> exec
| SELECT `timestampadd`(a, 2, 3) FROM tu
< error
| 1105 HY000 should not happened
> exec
| SELECT `get_format`(1, 2)
< ok 0 0
> exec
| SELECT `timestampdiff`()
< error
| 1582 42000 Incorrect parameter count in the call to native function 'timestampdiff'
> exec
| SELECT TIMESTAMPDIFF(DAY_HOUR, a, b) FROM tu
< error
| 1064 42000 …line 1 column 29 near "DAY_HOUR, a, b) FROM tu" 
> exec
| SELECT TIMESTAMPDIFF(YEAR_MONTH, a, b) FROM tu
< error
| 1064 42000 …line 1 column 31 near "YEAR_MONTH, a, b) FROM tu" 
> exec
| SELECT TIMESTAMPADD(MINUTE_SECOND, 1, a) FROM tu
< error
| 1064 42000 …line 1 column 33 near "MINUTE_SECOND, 1, a) FROM tu" 
> exec
| SELECT TIMESTAMPDIFF(`SECOND`, a, b) FROM tu
< error
| 1064 42000 …line 1 column 29 near "`SECOND`, a, b) FROM tu" 
> exec
| SELECT TIMESTAMPDIFF('SECOND', a, b) FROM tu
< error
| 1064 42000 …line 1 column 29 near "'SECOND', a, b) FROM tu" 
> exec
| SELECT GET_FORMAT(`DATE`, 'USA') FROM tu
< error
| 1064 42000 …line 1 column 24 near "`DATE`, 'USA') FROM tu" 
> exec
| SELECT GET_FORMAT('DATE', 'USA') FROM tu
< error
| 1064 42000 …line 1 column 24 near "'DATE', 'USA') FROM tu" 
> exec
| SELECT TIMESTAMPDIFF(SQL_TSI_MICROSECOND, a, b) FROM tu
< error
| 1064 42000 …line 1 column 40 near "SQL_TSI_MICROSECOND, a, b) FROM tu" 
> exec
| SELECT a + INTERVAL 1 SQL_TSI_MICROSECOND FROM tu
< error
| 1064 42000 …line 1 column 41 near "SQL_TSI_MICROSECOND FROM tu" 
> exec
| SELECT TIMESTAMPDIFF(NOTAUNIT, a, b) FROM tu
< error
| 1064 42000 …line 1 column 29 near "NOTAUNIT, a, b) FROM tu" 
> exec
| SELECT GET_FORMAT(YEAR, 'USA') FROM tu
< error
| 1064 42000 …line 1 column 22 near "YEAR, 'USA') FROM tu" 
> exec
| SELECT TIMESTAMPDIFF(SECOND, a) FROM tu
< error
| 1064 42000 …line 1 column 31 near ") FROM tu" 
> exec
| SELECT TIMESTAMPDIFF(SECOND, a, b, a) FROM tu
< error
| 1064 42000 …line 1 column 34 near ", a) FROM tu" 
> exec
| SELECT GET_FORMAT(DATE) FROM tu
< error
| 1064 42000 …line 1 column 23 near ") FROM tu" 
> exec
| SELECT `timestampdiff`(DAY_HOUR, a, b) FROM tu
< error
| 1064 42000 …line 1 column 31 near "DAY_HOUR, a, b) FROM tu" 
> exec
| DROP DATABASE IF EXISTS omni_tka
< ok 0 0
