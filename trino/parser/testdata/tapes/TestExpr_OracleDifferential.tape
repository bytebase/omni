-- pin: default trinodb/trino:482 sha256:90b35b7c603eaa1f889bf03981a62b75f998ee6c0f851d9f4e341b49a57022b6 482

== TestExpr_OracleDifferential default#1
> exec
| SELECT NULL
< ok - -
> exec
| SELECT TRUE
< ok - -
> exec
| SELECT FALSE
< ok - -
> exec
| SELECT 'a string'
< ok - -
> exec
| SELECT 'it''s escaped'
< ok - -
> exec
| SELECT U&'\0041'
< ok - -
> exec
| SELECT 123
< ok - -
> exec
| SELECT 1.5
< ok - -
> exec
| SELECT 1.5e10
< ok - -
> exec
| SELECT -3
< ok - -
> exec
| SELECT +4
< ok - -
> exec
| SELECT X'00ff'
< ok - -
> exec
| SELECT ?
< error
| INVALID_PARAMETER_USAGE 80 1:1
| line 1:1: Incorrect number of parameters: expected 1 but found 0
> exec
| SELECT DATE '2020-01-01'
< ok - -
> exec
| SELECT TIMESTAMP '2012-10-31 01:00:00'
< ok - -
> exec
| SELECT BIGINT '7'
< ok - -
> exec
| SELECT json '[1,2]'
< ok - -
> exec
| SELECT DOUBLE PRECISION '1.5'
< ok - -
> exec
| SELECT INTERVAL '1' DAY
< ok - -
> exec
| SELECT INTERVAL '1-2' YEAR TO MONTH
< ok - -
> exec
| SELECT INTERVAL '3' HOUR
< ok - -
> exec
| SELECT INTERVAL -'5' MINUTE
< ok - -
> exec
| SELECT INTERVAL '1 02:03:04' DAY TO SECOND
< ok - -
> exec
| SELECT a
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT a.b.c
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a.b.c' cannot be resolved
> exec
| SELECT a[1]
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT m['k'][1].f
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'm' cannot be resolved
> exec
| SELECT a.b[1]
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a.b' cannot be resolved
> exec
| SELECT a + b - c * d / e % f
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT a || b || c
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT - - 1
< ok - -
> exec
| SELECT -a * b
< error
| COLUMN_NOT_FOUND 47 1:9
| line 1:9: Column 'a' cannot be resolved
> exec
| SELECT ts AT TIME ZONE 'UTC'
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'ts' cannot be resolved
> exec
| SELECT ts AT TIME ZONE INTERVAL '1' HOUR
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'ts' cannot be resolved
> exec
| SELECT a * b AT TIME ZONE 'UTC'
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT TRUE AND FALSE OR NULL
< ok - -
> exec
| SELECT NOT a
< error
| COLUMN_NOT_FOUND 47 1:12
| line 1:12: Column 'a' cannot be resolved
> exec
| SELECT NOT NOT a
< error
| COLUMN_NOT_FOUND 47 1:16
| line 1:16: Column 'a' cannot be resolved
> exec
| SELECT NOT a = b
< error
| COLUMN_NOT_FOUND 47 1:12
| line 1:12: Column 'a' cannot be resolved
> exec
| SELECT NOT a AND b
< error
| COLUMN_NOT_FOUND 47 1:12
| line 1:12: Column 'a' cannot be resolved
> exec
| SELECT a <> b
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT a != b
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT a <= b
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT a >= b
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT a < b
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT a > b
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT a = b
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT a = ANY (SELECT 1)
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT a < ALL (SELECT 1)
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT a > SOME (SELECT x FROM t)
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT a BETWEEN 1 AND 2
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT a NOT BETWEEN 1 AND 2
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT a IN (1, 2, 3)
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT a NOT IN (1, 2, 3)
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT a IN (SELECT 1)
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT a NOT IN (SELECT 1)
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT a LIKE 'x%'
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT a NOT LIKE 'x%'
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT a LIKE 'x%' ESCAPE '!'
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT a IS NULL
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT a IS NOT NULL
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT a IS DISTINCT FROM b
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT a IS NOT DISTINCT FROM b
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT a IS TRUE
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT a IS NOT TRUE
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT a IS FALSE
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT a IS NOT FALSE
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT a IS UNKNOWN
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT a IS NOT UNKNOWN
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT a || b = c
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT (1, 2)
< ok - -
> exec
| SELECT (1, 2, 3)
< ok - -
> exec
| SELECT ROW(1, 'a')
< ok - -
> exec
| SELECT ROW(1)
< ok - -
> exec
| SELECT ARRAY[1, 2, 3]
< ok - -
> exec
| SELECT ARRAY[]
< ok - -
> exec
| SELECT (1)
< ok - -
> exec
| SELECT (SELECT 1)
< ok - -
> exec
| SELECT (SELECT max(x) FROM (SELECT a AS x FROM t) s)
< error
| TABLE_NOT_FOUND 46 1:48
| line 1:48: Table 'memory.default.t' does not exist
> exec
| SELECT EXISTS (SELECT 1)
< ok - -
> exec
| SELECT EXISTS (SELECT * FROM t WHERE a > 0)
< error
| TABLE_NOT_FOUND 46 1:30
| line 1:30: Table 'memory.default.t' does not exist
> exec
| SELECT EXISTS (TABLE t)
< error
| TABLE_NOT_FOUND 46 1:16
| line 1:16: Table 'memory.default.t' does not exist
> exec
| SELECT EXISTS (VALUES 1)
< ok - -
> exec
| SELECT a IN (SELECT b FROM t WHERE c IN (1, 2))
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'a' cannot be resolved
> exec
| SELECT CASE WHEN a > 1 THEN 'hi' ELSE 'lo' END
< error
| COLUMN_NOT_FOUND 47 1:18
| line 1:18: Column 'a' cannot be resolved
> exec
| SELECT CASE a WHEN 1 THEN 'one' WHEN 2 THEN 'two' END
< error
| COLUMN_NOT_FOUND 47 1:13
| line 1:13: Column 'a' cannot be resolved
> exec
| SELECT CASE WHEN a THEN 1 END
< error
| COLUMN_NOT_FOUND 47 1:18
| line 1:18: Column 'a' cannot be resolved
> exec
| SELECT CAST(a AS varchar)
< error
| COLUMN_NOT_FOUND 47 1:13
| line 1:13: Column 'a' cannot be resolved
> exec
| SELECT TRY_CAST(b AS bigint)
< error
| COLUMN_NOT_FOUND 47 1:17
| line 1:17: Column 'b' cannot be resolved
> exec
| SELECT CAST(x AS ARRAY(INTEGER))
< error
| COLUMN_NOT_FOUND 47 1:13
| line 1:13: Column 'x' cannot be resolved
> exec
| SELECT CAST(x AS ROW(a bigint, b varchar))
< error
| COLUMN_NOT_FOUND 47 1:13
| line 1:13: Column 'x' cannot be resolved
> exec
| SELECT CAST(x AS MAP(VARCHAR, ARRAY(INTEGER)))
< error
| COLUMN_NOT_FOUND 47 1:13
| line 1:13: Column 'x' cannot be resolved
> exec
| SELECT CAST(x AS INTERVAL DAY TO SECOND)
< error
| COLUMN_NOT_FOUND 47 1:13
| line 1:13: Column 'x' cannot be resolved
> exec
| SELECT CAST(ts AS timestamp) AT TIME ZONE 'UTC'
< error
| COLUMN_NOT_FOUND 47 1:13
| line 1:13: Column 'ts' cannot be resolved
> exec
| SELECT count(*)
< ok - -
> exec
| SELECT count(a)
< error
| COLUMN_NOT_FOUND 47 1:14
| line 1:14: Column 'a' cannot be resolved
> exec
| SELECT count(DISTINCT a)
< error
| COLUMN_NOT_FOUND 47 1:23
| line 1:23: Column 'a' cannot be resolved
> exec
| SELECT my.catalog.func(1)
< error
| FUNCTION_NOT_FOUND 6 1:8
| line 1:8: Function 'my.catalog.func' not registered
> exec
| SELECT "f"(1)
< error
| FUNCTION_NOT_FOUND 6 1:8
| line 1:8: Function 'f' not registered
> exec
| SELECT coalesce(a, b, 0)
< error
| COLUMN_NOT_FOUND 47 1:17
| line 1:17: Column 'a' cannot be resolved
> exec
| SELECT nullif(a, b)
< error
| COLUMN_NOT_FOUND 47 1:15
| line 1:15: Column 'a' cannot be resolved
> exec
| SELECT if(a > 0, 1, 2)
< error
| COLUMN_NOT_FOUND 47 1:11
| line 1:11: Column 'a' cannot be resolved
> exec
| SELECT count(*) FILTER (WHERE a > 0)
< error
| COLUMN_NOT_FOUND 47 1:31
| line 1:31: Column 'a' cannot be resolved
> exec
| SELECT count(t.*)
< error
| INVALID_FUNCTION_ARGUMENT 7 1:14
| line 1:14: label.* syntax is only supported as the only argument of row pattern count function
> exec
| SELECT f(t.*)
< error
| INVALID_FUNCTION_ARGUMENT 7 1:10
| line 1:10: label.* syntax is only supported as the only argument of row pattern count function
> exec
| SELECT array_agg(a ORDER BY b)
< error
| COLUMN_NOT_FOUND 47 1:18
| line 1:18: Column 'a' cannot be resolved
> exec
| SELECT array_agg(DISTINCT a ORDER BY a)
< error
| COLUMN_NOT_FOUND 47 1:27
| line 1:27: Column 'a' cannot be resolved
> exec
| SELECT array_agg(a ORDER BY b DESC NULLS LAST)
< error
| COLUMN_NOT_FOUND 47 1:18
| line 1:18: Column 'a' cannot be resolved
> exec
| SELECT sum(x) OVER (PARTITION BY a ORDER BY b)
< error
| COLUMN_NOT_FOUND 47 1:34
| line 1:34: Column 'a' cannot be resolved
> exec
| SELECT sum(x) OVER ()
< error
| COLUMN_NOT_FOUND 47 1:12
| line 1:12: Column 'x' cannot be resolved
> exec
| SELECT sum(x) OVER w
< error
| INVALID_WINDOW_REFERENCE 93 1:20
| line 1:20: Cannot resolve WINDOW name w
> exec
| SELECT row_number() OVER ()
< ok - -
> exec
| SELECT lag(x, 1) IGNORE NULLS OVER (ORDER BY y)
< error
| COLUMN_NOT_FOUND 47 1:46
| line 1:46: Column 'y' cannot be resolved
> exec
| SELECT lag(x) RESPECT NULLS OVER (ORDER BY y)
< error
| COLUMN_NOT_FOUND 47 1:44
| line 1:44: Column 'y' cannot be resolved
> exec
| SELECT count(DISTINCT a ORDER BY a) FILTER (WHERE a > 0) OVER (PARTITION BY b)
< error
| COLUMN_NOT_FOUND 47 1:77
| line 1:77: Column 'b' cannot be resolved
> exec
| SELECT FINAL first(x)
< error
| INVALID_PROCESSING_MODE 97 1:8
| line 1:8: FINAL semantics is not supported out of pattern recognition context
> exec
| SELECT sum(x) OVER (PARTITION BY a ORDER BY b ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW)
< error
| COLUMN_NOT_FOUND 47 1:34
| line 1:34: Column 'a' cannot be resolved
> exec
| SELECT sum(x) OVER (ORDER BY b RANGE 5 PRECEDING)
< error
| COLUMN_NOT_FOUND 47 1:30
| line 1:30: Column 'b' cannot be resolved
> exec
| SELECT sum(x) OVER (ORDER BY b GROUPS BETWEEN 1 PRECEDING AND 1 FOLLOWING)
< error
| COLUMN_NOT_FOUND 47 1:30
| line 1:30: Column 'b' cannot be resolved
> exec
| SELECT sum(x) OVER (ORDER BY b ROWS UNBOUNDED PRECEDING)
< error
| COLUMN_NOT_FOUND 47 1:30
| line 1:30: Column 'b' cannot be resolved
> exec
| SELECT sum(x) OVER (ROWS CURRENT ROW)
< error
| COLUMN_NOT_FOUND 47 1:12
| line 1:12: Column 'x' cannot be resolved
> exec
| SELECT sum(x) OVER (ORDER BY b RANGE BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED FOLLOWING)
< error
| COLUMN_NOT_FOUND 47 1:30
| line 1:30: Column 'b' cannot be resolved
> exec
| SELECT sum(x) OVER (ORDER BY b ROWS BETWEEN CURRENT ROW AND UNBOUNDED FOLLOWING)
< error
| COLUMN_NOT_FOUND 47 1:30
| line 1:30: Column 'b' cannot be resolved
> exec
| SELECT sum(x) OVER (w ORDER BY b)
< error
| INVALID_WINDOW_REFERENCE 93 1:21
| line 1:21: Cannot resolve WINDOW name w
> exec
| SELECT sum(x) OVER (w PARTITION BY a)
< error
| INVALID_WINDOW_REFERENCE 93 1:21
| line 1:21: Cannot resolve WINDOW name w
> exec
| SELECT sum(x) OVER (w ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW)
< error
| INVALID_WINDOW_REFERENCE 93 1:21
| line 1:21: Cannot resolve WINDOW name w
> exec
| SELECT a OVER (ORDER BY b)
< error
| INVALID_WINDOW_MEASURE 107 1:8
| line 1:8: Measure a is not defined in the corresponding window
> exec
| SELECT x -> x + 1
< error
| TYPE_MISMATCH 58 1:8
| line 1:8: Lambda expression should always be used inside a function
> exec
| SELECT () -> 1
< error
| TYPE_MISMATCH 58 1:8
| line 1:8: Lambda expression should always be used inside a function
> exec
| SELECT (x) -> x
< error
| TYPE_MISMATCH 58 1:8
| line 1:8: Lambda expression should always be used inside a function
> exec
| SELECT (k, v) -> k + v
< error
| TYPE_MISMATCH 58 1:8
| line 1:8: Lambda expression should always be used inside a function
> exec
| SELECT filter(arr, x -> x > 0)
< error
| COLUMN_NOT_FOUND 47 1:15
| line 1:15: Column 'arr' cannot be resolved
> exec
| SELECT reduce(a, 0, (s, x) -> s + x, s -> s)
< error
| COLUMN_NOT_FOUND 47 1:15
| line 1:15: Column 'a' cannot be resolved
> exec
| SELECT EXTRACT(YEAR FROM a)
< error
| COLUMN_NOT_FOUND 47 1:26
| line 1:26: Column 'a' cannot be resolved
> exec
| SELECT SUBSTRING('abc' FROM 1 FOR 2)
< ok - -
> exec
| SELECT SUBSTRING('abc' FROM 1)
< ok - -
> exec
| SELECT TRIM(BOTH 'x' FROM 'xax')
< ok - -
> exec
| SELECT TRIM(LEADING FROM '  x')
< ok - -
> exec
| SELECT TRIM('y' FROM '  x  ')
< ok - -
> exec
| SELECT TRIM('  x  ')
< ok - -
> exec
| SELECT TRIM('abc', 'a')
< ok - -
> exec
| SELECT NORMALIZE('abc')
< ok - -
> exec
| SELECT NORMALIZE('abc', NFC)
< ok - -
> exec
| SELECT POSITION('a' IN 'abc')
< ok - -
> exec
| SELECT substring
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'substring' cannot be resolved
> exec
| SELECT position
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'position' cannot be resolved
> exec
| SELECT substring('abc', 1, 2)
< ok - -
> exec
| SELECT substring('abc', 1)
< ok - -
> exec
| SELECT position('a', 'abc')
< error
| FUNCTION_NOT_FOUND 6 1:8
| line 1:8: Function 'position' not registered
> exec
| SELECT GROUPING(a, b)
< error
| COLUMN_NOT_FOUND 47 0:0
| Column 'a' cannot be resolved
> exec
| SELECT GROUPING(a)
< error
| COLUMN_NOT_FOUND 47 0:0
| Column 'a' cannot be resolved
> exec
| SELECT listagg(x, ',') WITHIN GROUP (ORDER BY x)
< error
| COLUMN_NOT_FOUND 47 1:16
| line 1:16: Column 'x' cannot be resolved
> exec
| SELECT listagg(DISTINCT x) WITHIN GROUP (ORDER BY x)
< error
| COLUMN_NOT_FOUND 47 1:25
| line 1:25: Column 'x' cannot be resolved
> exec
| SELECT listagg(x ON OVERFLOW ERROR) WITHIN GROUP (ORDER BY x)
< error
| COLUMN_NOT_FOUND 47 1:16
| line 1:16: Column 'x' cannot be resolved
> exec
| SELECT listagg(x ON OVERFLOW TRUNCATE '...' WITH COUNT) WITHIN GROUP (ORDER BY x)
< error
| COLUMN_NOT_FOUND 47 1:16
| line 1:16: Column 'x' cannot be resolved
> exec
| SELECT CURRENT_DATE
< ok - -
> exec
| SELECT CURRENT_TIME
< ok - -
> exec
| SELECT CURRENT_TIME(3)
< ok - -
> exec
| SELECT CURRENT_TIMESTAMP
< ok - -
> exec
| SELECT CURRENT_TIMESTAMP(6)
< ok - -
> exec
| SELECT LOCALTIME
< ok - -
> exec
| SELECT LOCALTIMESTAMP
< ok - -
> exec
| SELECT CURRENT_USER
< ok - -
> exec
| SELECT CURRENT_CATALOG
< ok - -
> exec
| SELECT CURRENT_SCHEMA
< ok - -
> exec
| SELECT CURRENT_PATH
< ok - -
> exec
| SELECT 1 +
< error
| SYNTAX_ERROR 1 1:11
| line 1:11: mismatched input '<EOF>'. Expecting: <expression>
> exec
| SELECT a AT TIME ZONE
< error
| SYNTAX_ERROR 1 1:22
| line 1:22: mismatched input '<EOF>'. Expecting: 'INTERVAL', <string>
> exec
| SELECT CAST(1 AS)
< error
| SYNTAX_ERROR 1 1:17
| line 1:17: mismatched input ')'. Expecting: <type>
> exec
| SELECT a BETWEEN 1
< error
| SYNTAX_ERROR 1 1:19
| line 1:19: mismatched input '<EOF>'. Expecting: '%', '*', '+', '-', '.', '/', 'AND', 'AT', '[', '||'
> exec
| SELECT EXTRACT(FROM a)
< error
| SYNTAX_ERROR 1 1:16
| line 1:16: mismatched input 'FROM'. Expecting: <identifier>
> exec
| SELECT TRIM(LEADING FROM)
< error
| SYNTAX_ERROR 1 1:25
| line 1:25: mismatched input ')'. Expecting: <expression>
> exec
| SELECT TRIM(FROM 'x')
< error
| SYNTAX_ERROR 1 1:8
| line 1:8: The 'trim' function must have specification, char or both arguments when it takes FROM
> exec
| SELECT NORMALIZE()
< error
| SYNTAX_ERROR 1 1:18
| line 1:18: mismatched input ')'. Expecting: <expression>
> exec
| SELECT LISTAGG(x)
< error
| SYNTAX_ERROR 1 1:18
| line 1:18: mismatched input '<EOF>'. Expecting: 'WITHIN'
> exec
| SELECT ARRAY[1, 2,
< error
| SYNTAX_ERROR 1 1:19
| line 1:19: mismatched input '<EOF>'. Expecting: <expression>
> exec
| SELECT CASE END
< error
| SYNTAX_ERROR 1 1:13
| line 1:13: mismatched input 'END'. Expecting: 'WHEN', <expression>
> exec
| SELECT CASE WHEN THEN 1 END
< error
| SYNTAX_ERROR 1 1:18
| line 1:18: mismatched input 'THEN'. Expecting: <expression>
> exec
| SELECT a = b = c
< error
| SYNTAX_ERROR 1 1:14
| line 1:14: mismatched input '='. Expecting: '%', '*', '+', ',', '-', '.', '/', 'AND', 'AS', 'AT', 'EXCEPT', 'FETCH', 'FROM', 'GROUP', 'HAVING', 'INTERSECT', 'LIMIT', 'OFFSET', 'OR', 'ORDER', 'UNION', 'WHERE', 'WINDOW', '[', '||', <EOF>, <identifier>
> exec
| SELECT a < b < c
< error
| SYNTAX_ERROR 1 1:14
| line 1:14: mismatched input '<'. Expecting: '%', '*', '+', ',', '-', '.', '/', 'AND', 'AS', 'AT', 'EXCEPT', 'FETCH', 'FROM', 'GROUP', 'HAVING', 'INTERSECT', 'LIMIT', 'OFFSET', 'OR', 'ORDER', 'UNION', 'WHERE', 'WINDOW', '[', '||', <EOF>, <identifier>
> exec
| SELECT a IN (1) IN (2)
< error
| SYNTAX_ERROR 1 1:17
| line 1:17: mismatched input 'IN'. Expecting: ',', 'AND', 'AS', 'EXCEPT', 'FETCH', 'FROM', 'GROUP', 'HAVING', 'INTERSECT', 'LIMIT', 'OFFSET', 'OR', 'ORDER', 'UNION', 'WHERE', 'WINDOW', <EOF>, <identifier>
> exec
| SELECT a IS NULL IS NULL
< error
| SYNTAX_ERROR 1 1:18
| line 1:18: mismatched input 'IS'. Expecting: ',', 'AND', 'AS', 'EXCEPT', 'FETCH', 'FROM', 'GROUP', 'HAVING', 'INTERSECT', 'LIMIT', 'OFFSET', 'OR', 'ORDER', 'UNION', 'WHERE', 'WINDOW', <EOF>, <identifier>
> exec
| SELECT a BETWEEN 1 AND 2 BETWEEN 3 AND 4
< error
| SYNTAX_ERROR 1 1:26
| line 1:26: mismatched input 'BETWEEN'. Expecting: '%', '*', '+', ',', '-', '.', '/', 'AND', 'AS', 'AT', 'EXCEPT', 'FETCH', 'FROM', 'GROUP', 'HAVING', 'INTERSECT', 'LIMIT', 'OFFSET', 'OR', 'ORDER', 'UNION', 'WHERE', 'WINDOW', '[', '||', <EOF>, <identifier>
> exec
| SELECT trim
< error
| SYNTAX_ERROR 1 1:12
| line 1:12: mismatched input '<EOF>'. Expecting: '('
> exec
| SELECT extract
< error
| SYNTAX_ERROR 1 1:15
| line 1:15: mismatched input '<EOF>'. Expecting: '('
> exec
| SELECT extract('a', 'b')
< error
| SYNTAX_ERROR 1 1:16
| line 1:16: mismatched input ''a''. Expecting: <identifier>
> exec
| SELECT listagg('a')
< error
| SYNTAX_ERROR 1 1:20
| line 1:20: mismatched input '<EOF>'. Expecting: 'WITHIN'
> exec
| SELECT grouping('a')
< error
| SYNTAX_ERROR 1 1:17
| line 1:17: mismatched input ''a''. Expecting: ')', <identifier>
> exec
| SELECT EXISTS ()
< error
| SYNTAX_ERROR 1 1:16
| line 1:16: mismatched input ')'. Expecting: <query>
> exec
| SELECT EXISTS (1)
< error
| SYNTAX_ERROR 1 1:16
| line 1:16: mismatched input '1'. Expecting: <query>
> exec
| SELECT EXISTS (1 + 2)
< error
| SYNTAX_ERROR 1 1:16
| line 1:16: mismatched input '1'. Expecting: <query>
> exec
| SELECT ()
< error
| SYNTAX_ERROR 1 1:10
| line 1:10: mismatched input '<EOF>'. Expecting: '->'
> exec
| SELECT a IN ()
< error
| SYNTAX_ERROR 1 1:10
| line 1:10: mismatched input 'IN'. Expecting: ',', 'AND', 'AS', 'EXCEPT', 'FETCH', 'FROM', 'GROUP', 'HAVING', 'INTERSECT', 'LIMIT', 'OFFSET', 'OR', 'ORDER', 'UNION', 'WHERE', 'WINDOW', <EOF>, <identifier>
