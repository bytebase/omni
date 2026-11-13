-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestFunctionArgumentOracle default#1
> exec
| CREATE TABLE t (c1 INT, c2 INT, b JSON, a DATETIME, profit INT)
< ok 0 0
> exec
| CREATE TABLE sales (profit INT)
< ok 0 0
> exec
| SELECT SUM(c1) FROM t
< ok 0 0
> exec
| SELECT SUM(DISTINCT c1) FROM t
< ok 0 0
> exec
| SELECT SUM(DISTINCTROW c1) FROM t
< ok 0 0
> exec
| SELECT SUM(ALL c1) FROM t
< ok 0 0
> exec
| SELECT SUM(DISTINCT ALL c1) FROM t
< ok 0 0
> exec
| SELECT SUM(DISTINCTROW ALL c1) FROM t
< ok 0 0
> exec
| SELECT AVG(DISTINCTROW ALL c1) FROM t
< ok 0 0
> exec
| SELECT MAX(DISTINCT ALL c1) FROM t
< ok 0 0
> exec
| SELECT MIN(DISTINCTROW c1) FROM t
< ok 0 0
> exec
| SELECT STD(ALL c1) FROM t
< ok 0 0
> exec
| SELECT STD(DISTINCT c1) FROM t
< error
| 1064 42000 …near 'DISTINCT c1) FROM t' at line 1
> exec
| SELECT STD(DISTINCTROW ALL c1) FROM t
< error
| 1064 42000 …near 'DISTINCTROW ALL c1) FROM t' at line 1
> exec
| SELECT STDDEV(DISTINCT c1) FROM t
< error
| 1064 42000 …near 'DISTINCT c1) FROM t' at line 1
> exec
| SELECT STDDEV_POP(ALL c1) FROM t
< ok 0 0
> exec
| SELECT STDDEV_SAMP(DISTINCT ALL c1) FROM t
< error
| 1064 42000 …near 'DISTINCT ALL c1) FROM t' at line 1
> exec
| SELECT VARIANCE(ALL c1) FROM t
< ok 0 0
> exec
| SELECT VAR_POP(DISTINCTROW c1) FROM t
< error
| 1064 42000 …near 'DISTINCTROW c1) FROM t' at line 1
> exec
| SELECT VAR_SAMP(ALL c1) FROM t
< ok 0 0
> exec
| SELECT GROUP_CONCAT(DISTINCTROW c2, c1) FROM t GROUP BY c1
< ok 0 0
> exec
| SELECT GROUP_CONCAT(ALL c2, c1) FROM t GROUP BY c1
< error
| 1064 42000 …near 'ALL c2, c1) FROM t GROUP BY c1' at line 1
> exec
| SELECT GROUP_CONCAT(DISTINCT ALL c2 ORDER BY c1 SEPARATOR '-') FROM t
< error
| 1064 42000 …near 'ALL c2 ORDER BY c1 SEPARATOR '-') FROM t' at line 1
> exec
| SELECT GROUP_CONCAT(c1) OVER () FROM t
< error
| 1235 42000 This version of MySQL doesn't yet support 'group_concat as window function'
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| SELECT GROUP_CONCAT(c1) OVER () FROM t
< error
| 1235 42000 This version of MySQL doesn't yet support 'group_concat as window function'
> exec ^
: string
| SELECT GROUP_CONCAT(c1) OVER () FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT BIT_AND(ALL c1) FROM t
< ok 0 0
> exec
| SELECT BIT_OR(ALL c1) FROM t
< ok 0 0
> exec
| SELECT BIT_XOR(ALL c1) FROM t
< ok 0 0
> exec
| SELECT BIT_AND(DISTINCT c1) FROM t
< error
| 1064 42000 …near 'DISTINCT c1) FROM t' at line 1
> exec
| SELECT BIT_OR(DISTINCTROW c1) FROM t
< error
| 1064 42000 …near 'DISTINCTROW c1) FROM t' at line 1
> exec
| SELECT BIT_XOR(DISTINCT ALL c1) FROM t
< error
| 1064 42000 …near 'DISTINCT ALL c1) FROM t' at line 1
> exec
| SELECT JSON_ARRAYAGG(ALL c1) FROM t
< ok 0 0
> exec
| SELECT JSON_ARRAYAGG(DISTINCT c1) FROM t
< error
| 1064 42000 …near 'DISTINCT c1) FROM t' at line 1
> exec
| SELECT JSON_OBJECTAGG(c1, c2) FROM t
< ok 0 0
> exec
| SELECT JSON_OBJECTAGG(ALL c1, c2) FROM t
< ok 0 0
> exec
| SELECT JSON_OBJECTAGG(c1, ALL c2) FROM t
< ok 0 0
> exec
| SELECT JSON_OBJECTAGG(ALL c1, ALL c2) FROM t
< ok 0 0
> exec
| SELECT JSON_OBJECTAGG(DISTINCT c1, c2) FROM t
< error
| 1064 42000 …near 'DISTINCT c1, c2) FROM t' at line 1
> exec
| SELECT JSON_OBJECTAGG(ALL ALL c1, c2) FROM t
< error
| 1064 42000 …near 'ALL c1, c2) FROM t' at line 1
> exec
| SELECT COUNT(*) FROM t
< ok 0 0
> exec
| SELECT COUNT(c1) FROM t
< ok 0 0
> exec
| SELECT COUNT(ALL c1) FROM t
< ok 0 0
> exec
| SELECT COUNT(DISTINCT c1, c2) FROM t
< ok 0 0
> exec
| SELECT COUNT(DISTINCTROW c1, c2) FROM t
< ok 0 0
> exec
| SELECT COUNT(ALL profit) OVER () FROM sales
< ok 0 0
> exec
| SELECT SUM(DISTINCT c1) OVER () FROM t
< error
| 1235 42000 This version of MySQL doesn't yet support '<window function>(DISTINCT ..)'
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| SELECT SUM(DISTINCT c1) OVER () FROM t
< error
| 1235 42000 This version of MySQL doesn't yet support '<window function>(DISTINCT ..)'
> exec ^
: string
| SELECT SUM(DISTINCT c1) OVER () FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT COUNT(DISTINCT c1) OVER () FROM t
< error
| 1235 42000 This version of MySQL doesn't yet support '<window function>(DISTINCT ..)'
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| SELECT COUNT(DISTINCT c1) OVER () FROM t
< error
| 1235 42000 This version of MySQL doesn't yet support '<window function>(DISTINCT ..)'
> exec ^
: string
| SELECT COUNT(DISTINCT c1) OVER () FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT COUNT(DISTINCT *) FROM t
< error
| 1064 42000 …near '*) FROM t' at line 1
> exec
| SELECT COUNT(ALL *) FROM t
< ok 0 0
> exec
| SELECT COUNT(c1, c2) FROM t
< error
| 1064 42000 …near ', c2) FROM t' at line 1
> exec
| SELECT SUM(ALL ALL c1) FROM t
< error
| 1064 42000 …near 'ALL c1) FROM t' at line 1
> exec
| SELECT SUM(ALL DISTINCT c1) FROM t
< error
| 1064 42000 …near 'DISTINCT c1) FROM t' at line 1
> exec
| SELECT SUM(DISTINCT DISTINCT c1) FROM t
< error
| 1064 42000 …near 'DISTINCT c1) FROM t' at line 1
> exec
| SELECT SUM(c1 ALL) FROM t
< error
| 1064 42000 …near 'ALL) FROM t' at line 1
> exec
| SELECT ABS(DISTINCT c1) FROM t
< error
| 1064 42000 …near 'DISTINCT c1) FROM t' at line 1
> exec
| SELECT ANY_VALUE(DISTINCT c1) FROM t
< error
| 1064 42000 …near 'DISTINCT c1) FROM t' at line 1
> exec
| SELECT APPROX_COUNT_DISTINCT(DISTINCT c1) FROM t
< error
| 1064 42000 …near 'DISTINCT c1) FROM t' at line 1
> exec
| SELECT `sum`(DISTINCT c1) FROM t
< error
| 1064 42000 …near 'DISTINCT c1) FROM t' at line 1
> exec
| SELECT `count`(DISTINCT c1) FROM t
< error
| 1064 42000 …near 'DISTINCT c1) FROM t' at line 1
> exec
| SELECT AVG() FROM t
< error
| 1064 42000 …near ') FROM t' at line 1
> exec
| SELECT AVG(c1, c2) FROM t
< error
| 1064 42000 …near ', c2) FROM t' at line 1
> exec
| SELECT BIT_AND() FROM t
< error
| 1064 42000 …near ') FROM t' at line 1
> exec
| SELECT STD(c1, c2) FROM t
< error
| 1064 42000 …near ', c2) FROM t' at line 1
> exec
| SELECT GROUP_CONCAT() FROM t
< error
| 1064 42000 …near ') FROM t' at line 1
> exec
| SELECT JSON_ARRAYAGG(c1, c2) FROM t
< error
| 1064 42000 …near ', c2) FROM t' at line 1
> exec
| SELECT JSON_OBJECTAGG(c1, c2, c1) FROM t
< error
| 1064 42000 …near ', c1) FROM t' at line 1
> exec
| SELECT SUM(*) FROM t
< error
| 1064 42000 …near '*) FROM t' at line 1
> exec
| SELECT SUM(t.*) FROM t
< error
| 1064 42000 …near '*) FROM t' at line 1
> exec
| SELECT COUNT(t.*) FROM t
< error
| 1064 42000 …near '*) FROM t' at line 1
> exec
| SELECT ABS(*) FROM t
< error
| 1064 42000 …near '*) FROM t' at line 1
> exec
| SELECT GROUP_CONCAT(*) FROM t
< error
| 1064 42000 …near '*) FROM t' at line 1
> exec
| SELECT EXTRACT(MICROSECOND FROM a) FROM t
< ok 0 0
> exec
| SELECT EXTRACT(SECOND FROM a) FROM t
< ok 0 0
> exec
| SELECT EXTRACT(MINUTE FROM a) FROM t
< ok 0 0
> exec
| SELECT EXTRACT(HOUR FROM a) FROM t
< ok 0 0
> exec
| SELECT EXTRACT(DAY FROM a) FROM t
< ok 0 0
> exec
| SELECT EXTRACT(WEEK FROM a) FROM t
< ok 0 0
> exec
| SELECT EXTRACT(MONTH FROM a) FROM t
< ok 0 0
> exec
| SELECT EXTRACT(QUARTER FROM a) FROM t
< ok 0 0
> exec
| SELECT EXTRACT(YEAR FROM a) FROM t
< ok 0 0
> exec
| SELECT EXTRACT(SECOND_MICROSECOND FROM a) FROM t
< ok 0 0
> exec
| SELECT EXTRACT(MINUTE_MICROSECOND FROM a) FROM t
< ok 0 0
> exec
| SELECT EXTRACT(MINUTE_SECOND FROM a) FROM t
< ok 0 0
> exec
| SELECT EXTRACT(HOUR_MICROSECOND FROM a) FROM t
< ok 0 0
> exec
| SELECT EXTRACT(HOUR_SECOND FROM a) FROM t
< ok 0 0
> exec
| SELECT EXTRACT(HOUR_MINUTE FROM a) FROM t
< ok 0 0
> exec
| SELECT EXTRACT(DAY_MICROSECOND FROM a) FROM t
< ok 0 0
> exec
| SELECT EXTRACT(DAY_SECOND FROM a) FROM t
< ok 0 0
> exec
| SELECT EXTRACT(DAY_MINUTE FROM a) FROM t
< ok 0 0
> exec
| SELECT EXTRACT(DAY_HOUR FROM a) FROM t
< ok 0 0
> exec
| SELECT EXTRACT(YEAR_MONTH FROM a) FROM t
< ok 0 0
> exec
| SELECT EXTRACT(SQL_TSI_DAY FROM a) FROM t
< ok 0 0
> exec
| SELECT EXTRACT(`DAY_HOUR` FROM a) FROM t
< error
| 1064 42000 …near '`DAY_HOUR` FROM a) FROM t' at line 1
> exec
| SELECT EXTRACT('DAY_HOUR' FROM a) FROM t
< error
| 1064 42000 …near ''DAY_HOUR' FROM a) FROM t' at line 1
> exec
| SELECT EXTRACT(c1 FROM a) FROM t
< error
| 1064 42000 …near 'c1 FROM a) FROM t' at line 1
> exec
| SELECT EXTRACT(NOTAUNIT FROM a) FROM t
< error
| 1064 42000 …near 'NOTAUNIT FROM a) FROM t' at line 1
> exec
| SELECT EXTRACT(SQL_TSI_MICROSECOND FROM a) FROM t
< error
| 1064 42000 …near 'SQL_TSI_MICROSECOND FROM a) FROM t' at line 1
> exec
| SELECT EXTRACT(DAY_HOUR, a) FROM t
< error
| 1064 42000 …near ', a) FROM t' at line 1
> exec
| SELECT EXTRACT(DAY_HOUR FROM a) OVER () FROM t
< error
| 1064 42000 …near 'OVER () FROM t' at line 1
> exec
| SELECT `extract`(DAY_HOUR FROM a) FROM t
< error
| 1064 42000 …near 'DAY_HOUR FROM a) FROM t' at line 1
> exec
| SELECT INTERVAL(1, 0, 1, 2)
< ok 0 0
> exec
| SELECT (INTERVAL(1, 0, 1, 2)+5)*7+INTERVAL(1, 0, 1, 2)/2
< ok 0 0
> exec
| SELECT INTERVAL(0, (1*5)/2)+INTERVAL(5, 4, 3)
< ok 0 0
> exec
| SELECT INTERVAL((SELECT 1), 2)
< ok 0 0
> exec
| SELECT INTERVAL(c1, 1, 2) FROM t
< ok 0 0
> exec
| SELECT 1 + INTERVAL(1, 2)
< ok 0 0
> exec
| SELECT INTERVAL 1 DAY + INTERVAL(1, 2)
< ok 0 0
> exec
| INSERT INTO t (c1) VALUES (INTERVAL(1, 2))
< ok 1 0
> exec
| SELECT INTERVAL(1)
< error
| 1064 42000 …near '' at line 1
> exec
| SELECT INTERVAL()
< error
| 1064 42000 …near ')' at line 1
> exec
| SELECT a + INTERVAL (1) DAY FROM t
< ok 0 0
> exec
| SELECT DATE_ADD('2008-01-02', INTERVAL INTERVAL(1, 0, 1) DAY)
< ok 0 0
> exec
| SELECT a + INTERVAL -1 DAY FROM t
< ok 0 0
> exec
| SELECT INTERVAL(ALL 1, 2)
< error
| 1064 42000 …near 'ALL 1, 2)' at line 1
> exec
| SELECT INTERVAL(SELECT 1, 2)
< error
| 1064 42000 …near '' at line 1
> exec
| SELECT INTERVAL(c1) OVER () FROM t
< error
| 1064 42000 …near 'OVER () FROM t' at line 1
> exec
| SELECT CAST(b AS SIGNED ARRAY) FROM t
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| SELECT CAST(b AS SIGNED ARRAY) FROM t
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec ^
: string
| SELECT CAST(b AS SIGNED ARRAY) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT CAST(b AS UNSIGNED ARRAY) FROM t
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| SELECT CAST(b AS UNSIGNED ARRAY) FROM t
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec ^
: string
| SELECT CAST(b AS UNSIGNED ARRAY) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT CAST(b AS SIGNED INTEGER ARRAY) FROM t
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| SELECT CAST(b AS SIGNED INTEGER ARRAY) FROM t
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec ^
: string
| SELECT CAST(b AS SIGNED INTEGER ARRAY) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT CAST(b AS CHAR(10) ARRAY) FROM t
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| SELECT CAST(b AS CHAR(10) ARRAY) FROM t
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec ^
: string
| SELECT CAST(b AS CHAR(10) ARRAY) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT CAST(b AS DECIMAL(10,2) ARRAY) FROM t
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| SELECT CAST(b AS DECIMAL(10,2) ARRAY) FROM t
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec ^
: string
| SELECT CAST(b AS DECIMAL(10,2) ARRAY) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT CAST(b AS DATE ARRAY) FROM t
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| SELECT CAST(b AS DATE ARRAY) FROM t
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec ^
: string
| SELECT CAST(b AS DATE ARRAY) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT CAST(b AS DATETIME(3) ARRAY) FROM t
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| SELECT CAST(b AS DATETIME(3) ARRAY) FROM t
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec ^
: string
| SELECT CAST(b AS DATETIME(3) ARRAY) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT CAST(b AS TIME ARRAY) FROM t
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| SELECT CAST(b AS TIME ARRAY) FROM t
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec ^
: string
| SELECT CAST(b AS TIME ARRAY) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT CAST(b AS BINARY(10) ARRAY) FROM t
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| SELECT CAST(b AS BINARY(10) ARRAY) FROM t
< error
| 1235 42000 This version of MySQL doesn't yet support 'Use of CAST( .. AS .. ARRAY) outside of functional index in CREATE(non-SELECT)/ALTER TABLE or in general expressions'
> exec ^
: string
| SELECT CAST(b AS BINARY(10) ARRAY) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| CREATE INDEX idx ON t ((CAST(b AS SIGNED ARRAY)))
< ok 0 0
> exec
| SELECT CAST(b AS SIGNED ARRAY ARRAY) FROM t
< error
| 1064 42000 …near 'ARRAY) FROM t' at line 1
> exec
| SELECT CAST(b AS ARRAY) FROM t
< error
| 1064 42000 …near 'ARRAY) FROM t' at line 1
> exec
| SELECT CONVERT(b, SIGNED ARRAY) FROM t
< error
| 1064 42000 …near 'ARRAY) FROM t' at line 1
> exec
| SELECT POSITION('bar' IN 'foobarbar')
< ok 0 0
> exec
| SELECT POSITION(c1 IN c2) FROM t
< ok 0 0
> exec
| SELECT POSITION('bar' IN 'foobarbar') + 1
< ok 0 0
> exec
| SELECT POSITION('a' + 1 IN 'b')
< ok 0 0
> exec
| SELECT POSITION('a' IN 'b' = 1)
< ok 0 0
> exec
| SELECT POSITION('bar', 'foobarbar')
< error
| 1064 42000 …near ', 'foobarbar')' at line 1
> exec
| SELECT POSITION(1 = 1 IN 'x')
< error
| 1064 42000 …near '= 1 IN 'x')' at line 1
> exec
| SELECT POSITION(DISTINCT 'a' IN 'b')
< error
| 1064 42000 …near 'DISTINCT 'a' IN 'b')' at line 1
> exec
| SELECT POSITION('a' IN 'b') OVER ()
< error
| 1064 42000 …near 'OVER ()' at line 1
> exec
| SELECT POSITION()
< error
| 1064 42000 …near ')' at line 1
> exec
| SELECT POSITION('a')
< error
| 1064 42000 …near ')' at line 1
> exec
| SELECT POSITION(* IN 'a')
< error
| 1064 42000 …near '* IN 'a')' at line 1
> exec
| SELECT `position`('bar', 'foo')
< error
| 1630 42000 FUNCTION omni_t_mysql_parser.position does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| SELECT `position`('bar', 'foo')
< error
| 1630 42000 FUNCTION omni_t_mysql_parser.position does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT `position`('bar', 'foo')
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT `position`('bar' IN 'foo')
< error
| 1064 42000 …near ''foo')' at line 1
