-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestMariaDBTemporalKeywordArgOracle default#1
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
| 1630 42000 FUNCTION s.timestampdiff does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT s.timestampdiff(1, 2, 3)
< error
| 1630 42000 FUNCTION s.timestampdiff does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT s.timestampdiff(1, 2, 3)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT s.get_format(1, 2)
< error
| 1630 42000 FUNCTION s.get_format does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT s.get_format(1, 2)
< error
| 1630 42000 FUNCTION s.get_format does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT s.get_format(1, 2)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT `timestampdiff`(1, 2, 3)
< error
| 1630 42000 FUNCTION test.timestampdiff does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT `timestampdiff`(1, 2, 3)
< error
| 1630 42000 FUNCTION test.timestampdiff does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT `timestampdiff`(1, 2, 3)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT `TIMESTAMPDIFF`(SECOND, a, b) FROM tu
< error
| 1630 42000 FUNCTION test.TIMESTAMPDIFF does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT `TIMESTAMPDIFF`(SECOND, a, b) FROM tu
< error
| 1630 42000 FUNCTION test.TIMESTAMPDIFF does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT `TIMESTAMPDIFF`(SECOND, a, b) FROM tu
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT `timestampadd`(a, 2, 3) FROM tu
< error
| 1630 42000 FUNCTION test.timestampadd does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT `timestampadd`(a, 2, 3) FROM tu
< error
| 1630 42000 FUNCTION test.timestampadd does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT `timestampadd`(a, 2, 3) FROM tu
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT `get_format`(1, 2)
< error
| 1630 42000 FUNCTION test.get_format does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT `get_format`(1, 2)
< error
| 1630 42000 FUNCTION test.get_format does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT `get_format`(1, 2)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT `timestampdiff`()
< error
| 1630 42000 FUNCTION test.timestampdiff does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT `timestampdiff`()
< error
| 1630 42000 FUNCTION test.timestampdiff does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT `timestampdiff`()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT TIMESTAMPDIFF(DAY_HOUR, a, b) FROM tu
< error
| 1064 42000 …near 'DAY_HOUR, a, b) FROM tu' at line 1
> exec
| SELECT TIMESTAMPDIFF(YEAR_MONTH, a, b) FROM tu
< error
| 1064 42000 …near 'YEAR_MONTH, a, b) FROM tu' at line 1
> exec
| SELECT TIMESTAMPADD(MINUTE_SECOND, 1, a) FROM tu
< error
| 1064 42000 …near 'MINUTE_SECOND, 1, a) FROM tu' at line 1
> exec
| SELECT TIMESTAMPDIFF(`SECOND`, a, b) FROM tu
< error
| 1064 42000 …near '`SECOND`, a, b) FROM tu' at line 1
> exec
| SELECT TIMESTAMPDIFF('SECOND', a, b) FROM tu
< error
| 1064 42000 …near ''SECOND', a, b) FROM tu' at line 1
> exec
| SELECT GET_FORMAT(`DATE`, 'USA') FROM tu
< error
| 1064 42000 …near '`DATE`, 'USA') FROM tu' at line 1
> exec
| SELECT GET_FORMAT('DATE', 'USA') FROM tu
< error
| 1064 42000 …near ''DATE', 'USA') FROM tu' at line 1
> exec
| SELECT TIMESTAMPDIFF(SQL_TSI_MICROSECOND, a, b) FROM tu
< error
| 1064 42000 …near 'SQL_TSI_MICROSECOND, a, b) FROM tu' at line 1
> exec
| SELECT a + INTERVAL 1 SQL_TSI_MICROSECOND FROM tu
< error
| 1064 42000 …near 'SQL_TSI_MICROSECOND FROM tu' at line 1
> exec
| SELECT TIMESTAMPDIFF(NOTAUNIT, a, b) FROM tu
< error
| 1064 42000 …near 'NOTAUNIT, a, b) FROM tu' at line 1
> exec
| SELECT GET_FORMAT(YEAR, 'USA') FROM tu
< error
| 1064 42000 …near 'YEAR, 'USA') FROM tu' at line 1
> exec
| SELECT TIMESTAMPDIFF(SECOND, a) FROM tu
< error
| 1064 42000 …near ') FROM tu' at line 1
> exec
| SELECT TIMESTAMPDIFF(SECOND, a, b, a) FROM tu
< error
| 1064 42000 …near ' a) FROM tu' at line 1
> exec
| SELECT GET_FORMAT(DATE) FROM tu
< error
| 1064 42000 …near ') FROM tu' at line 1
> exec
| SELECT `timestampdiff`(DAY_HOUR, a, b) FROM tu
< error
| 1064 42000 …near 'DAY_HOUR, a, b) FROM tu' at line 1
