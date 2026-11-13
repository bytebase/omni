-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOceanBaseTimeUnitGrammarOracle default#1
> exec
| CREATE DATABASE IF NOT EXISTS obunit
< ok 1 0
> exec
| USE obunit
< ok 0 0
> exec
| DROP TABLE IF EXISTS obunit_t
< ok 0 0
> exec
| CREATE TABLE obunit_t (a INT, d DATETIME)
< ok 0 0
> exec
| PREPARE obunit_probe FROM 'SELECT TIMESTAMPDIFF(SQL_TSI_SECOND, ''2020-01-01'', ''2020-01-02'')'
< error
| 1064 42000 …near 'SQL_TSI_SECOND, '2020-01-01', '2020-01-02')' at line 1
> exec
| PREPARE obunit_probe FROM 'SELECT TIMESTAMPDIFF(SQL_TSI_MINUTE, ''2020-01-01'', ''2020-01-02'')'
< error
| 1064 42000 …near 'SQL_TSI_MINUTE, '2020-01-01', '2020-01-02')' at line 1
> exec
| PREPARE obunit_probe FROM 'SELECT TIMESTAMPDIFF(SQL_TSI_HOUR, ''2020-01-01'', ''2020-01-02'')'
< error
| 1064 42000 …near 'SQL_TSI_HOUR, '2020-01-01', '2020-01-02')' at line 1
> exec
| PREPARE obunit_probe FROM 'SELECT TIMESTAMPDIFF(SQL_TSI_DAY, ''2020-01-01'', ''2020-01-02'')'
< error
| 1064 42000 …near 'SQL_TSI_DAY, '2020-01-01', '2020-01-02')' at line 1
> exec
| PREPARE obunit_probe FROM 'SELECT TIMESTAMPDIFF(SQL_TSI_WEEK, ''2020-01-01'', ''2020-01-02'')'
< error
| 1064 42000 …near 'SQL_TSI_WEEK, '2020-01-01', '2020-01-02')' at line 1
> exec
| PREPARE obunit_probe FROM 'SELECT TIMESTAMPDIFF(SQL_TSI_MONTH, ''2020-01-01'', ''2020-01-02'')'
< error
| 1064 42000 …near 'SQL_TSI_MONTH, '2020-01-01', '2020-01-02')' at line 1
> exec
| PREPARE obunit_probe FROM 'SELECT TIMESTAMPDIFF(SQL_TSI_QUARTER, ''2020-01-01'', ''2020-01-02'')'
< error
| 1064 42000 …near 'SQL_TSI_QUARTER, '2020-01-01', '2020-01-02')' at line 1
> exec
| PREPARE obunit_probe FROM 'SELECT TIMESTAMPDIFF(SQL_TSI_YEAR, ''2020-01-01'', ''2020-01-02'')'
< error
| 1064 42000 …near 'SQL_TSI_YEAR, '2020-01-01', '2020-01-02')' at line 1
> exec
| PREPARE obunit_probe FROM 'SELECT TIMESTAMPDIFF(sql_tsi_second, ''2020-01-01'', ''2020-01-02'')'
< error
| 1064 42000 …near 'sql_tsi_second, '2020-01-01', '2020-01-02')' at line 1
> exec
| PREPARE obunit_probe FROM 'SELECT TIMESTAMPADD(SQL_TSI_DAY, 1, ''2020-01-01'')'
< error
| 1064 42000 …near 'SQL_TSI_DAY, 1, '2020-01-01')' at line 1
> exec
| PREPARE obunit_probe FROM 'SELECT TIMESTAMPADD(SQL_TSI_SECOND, 1, ''2020-01-01'')'
< error
| 1064 42000 …near 'SQL_TSI_SECOND, 1, '2020-01-01')' at line 1
> exec
| PREPARE obunit_probe FROM 'SELECT DATE_ADD(''2020-01-01'', INTERVAL 1 SQL_TSI_DAY)'
< error
| 1064 42000 …near 'SQL_TSI_DAY)' at line 1
> exec
| PREPARE obunit_probe FROM 'SELECT DATE_SUB(''2020-01-01'', INTERVAL 1 SQL_TSI_DAY)'
< error
| 1064 42000 …near 'SQL_TSI_DAY)' at line 1
> exec
| PREPARE obunit_probe FROM 'SELECT ADDDATE(''2020-01-01'', INTERVAL 1 SQL_TSI_DAY)'
< error
| 1064 42000 …near 'SQL_TSI_DAY)' at line 1
> exec
| PREPARE obunit_probe FROM 'SELECT SUBDATE(''2020-01-01'', INTERVAL 1 SQL_TSI_DAY)'
< error
| 1064 42000 …near 'SQL_TSI_DAY)' at line 1
> exec
| PREPARE obunit_probe FROM 'SELECT ''2020-01-01'' + INTERVAL 1 SQL_TSI_DAY'
< error
| 1064 42000 …near 'SQL_TSI_DAY' at line 1
> exec
| PREPARE obunit_probe FROM 'SELECT INTERVAL 1 SQL_TSI_DAY + ''2020-01-01'''
< error
| 1064 42000 …near 'SQL_TSI_DAY + '2020-01-01'' at line 1
> exec
| PREPARE obunit_probe FROM 'SELECT ''2020-01-01'' - INTERVAL 1 SQL_TSI_DAY'
< error
| 1064 42000 …near 'SQL_TSI_DAY' at line 1
> exec
| PREPARE obunit_probe FROM 'SELECT EXTRACT(SQL_TSI_DAY FROM ''2020-01-01'')'
< error
| 1064 42000 …near 'SQL_TSI_DAY FROM '2020-01-01')' at line 1
> exec
| PREPARE obunit_probe FROM 'SELECT EXTRACT(SQL_TSI_YEAR FROM ''2020-01-01'')'
< error
| 1064 42000 …near 'SQL_TSI_YEAR FROM '2020-01-01')' at line 1
> exec
| PREPARE obunit_probe FROM 'SELECT SUM(a) OVER (ORDER BY d RANGE INTERVAL 1 SQL_TSI_DAY PRECEDING) FROM obunit_t'
< error
| 1064 42000 …near 'SQL_TSI_DAY PRECEDING) FROM obunit_t' at line 1
> exec
| PREPARE obunit_probe FROM 'SELECT SUM(a) OVER (ORDER BY d RANGE BETWEEN INTERVAL 1 SQL_TSI_DAY PRECEDING AND CURRENT ROW) FROM obunit_t'
< error
| 1064 42000 …near 'SQL_TSI_DAY PRECEDING AND CURRENT ROW) FROM obunit_t' at line 1
> exec
| PREPARE obunit_probe FROM 'CREATE TABLE t (a DATETIME, b INT AS (TIMESTAMPDIFF(SQL_TSI_DAY, a, ''2020-01-01'')))'
< error
| 1064 42000 …near 'SQL_TSI_DAY, a, '2020-01-01')))' at line 1
> exec
| PREPARE obunit_probe FROM 'CREATE VIEW v AS SELECT TIMESTAMPDIFF(SQL_TSI_DAY, d, d) FROM obunit_t'
< error
| 1064 42000 …near 'SQL_TSI_DAY, d, d) FROM obunit_t' at line 1
> exec
| PREPARE obunit_probe FROM 'SELECT TIMESTAMPDIFF(SECOND, ''2020-01-01'', ''2020-01-02'')'
< ok 0 0
> exec
| PREPARE obunit_probe FROM 'SELECT TIMESTAMPDIFF(MICROSECOND, ''2020-01-01'', ''2020-01-02'')'
< ok 0 0
> exec
| PREPARE obunit_probe FROM 'SELECT TIMESTAMPADD(WEEK, 1, ''2020-01-01'')'
< ok 0 0
> exec
| PREPARE obunit_probe FROM 'SELECT DATE_ADD(''2020-01-01'', INTERVAL 1 DAY)'
< ok 0 0
> exec
| PREPARE obunit_probe FROM 'SELECT DATE_ADD(''2020-01-01'', INTERVAL ''1 1'' DAY_HOUR)'
< ok 0 0
> exec
| PREPARE obunit_probe FROM 'SELECT DATE_ADD(''2020-01-01'', INTERVAL 1 MICROSECOND)'
< ok 0 0
> exec
| PREPARE obunit_probe FROM 'SELECT EXTRACT(DAY_HOUR FROM ''2020-01-01'')'
< ok 0 0
> exec
| PREPARE obunit_probe FROM 'SELECT SUM(a) OVER (ORDER BY d RANGE INTERVAL 1 DAY PRECEDING) FROM obunit_t'
< ok 0 0
> exec
| PREPARE obunit_probe FROM 'SELECT sql_tsi_second FROM (SELECT 1 AS sql_tsi_second) t'
< ok 0 0
> exec
| PREPARE obunit_probe FROM 'CREATE EVENT e ON SCHEDULE EVERY 1 SQL_TSI_DAY DO SELECT 1'
< error
| 1064 42000 …near 'SQL_TSI_DAY' at line 1
> exec
| PREPARE obunit_probe FROM 'CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1'
< ok 0 0
> exec
| PREPARE obunit_probe FROM 'SELECT TIMESTAMPDIFF(DAY_HOUR, ''2020-01-01'', ''2020-01-02'')'
< ok 0 0
> exec
| PREPARE obunit_probe FROM 'SELECT TIMESTAMPDIFF(YEAR_MONTH, ''2020-01-01'', ''2020-01-02'')'
< ok 0 0
> exec
| PREPARE obunit_probe FROM 'SELECT TIMESTAMPADD(DAY_HOUR, 1, ''2020-01-01'')'
< ok 0 0
> exec
| PREPARE obunit_probe FROM 'SELECT TIMESTAMPADD(MINUTE_SECOND, 1, ''2020-01-01'')'
< ok 0 0
> exec
| PREPARE obunit_probe FROM 'SELECT TIMESTAMPDIFF(`SECOND`, ''2020-01-01'', ''2020-01-02'')'
< error
| 1064 42000 …near '`SECOND`, '2020-01-01', '2020-01-02')' at line 1
> exec
| PREPARE obunit_probe FROM 'SELECT TIMESTAMPDIFF(''SECOND'', ''2020-01-01'', ''2020-01-02'')'
< error
| 1064 42000 …near ''SECOND', '2020-01-01', '2020-01-02')' at line 1
> exec
| DROP DATABASE IF EXISTS obunit
< ok 1 0
