-- pin: default trinodb/trino:482 sha256:90b35b7c603eaa1f889bf03981a62b75f998ee6c0f851d9f4e341b49a57022b6 482

== TestDataType_OracleDifferential default#1
> exec
| SELECT CAST(NULL AS boolean)
< ok - -
> exec
| SELECT CAST(NULL AS tinyint)
< ok - -
> exec
| SELECT CAST(NULL AS smallint)
< ok - -
> exec
| SELECT CAST(NULL AS integer)
< ok - -
> exec
| SELECT CAST(NULL AS int)
< ok - -
> exec
| SELECT CAST(NULL AS bigint)
< ok - -
> exec
| SELECT CAST(NULL AS real)
< ok - -
> exec
| SELECT CAST(NULL AS double)
< ok - -
> exec
| SELECT CAST(NULL AS double precision)
< ok - -
> exec
| SELECT CAST(NULL AS decimal)
< ok - -
> exec
| SELECT CAST(NULL AS decimal(38))
< ok - -
> exec
| SELECT CAST(NULL AS decimal(38, 10))
< ok - -
> exec
| SELECT CAST(NULL AS varchar)
< ok - -
> exec
| SELECT CAST(NULL AS varchar(100))
< ok - -
> exec
| SELECT CAST(NULL AS char)
< ok - -
> exec
| SELECT CAST(NULL AS char(10))
< ok - -
> exec
| SELECT CAST(NULL AS varbinary)
< ok - -
> exec
| SELECT CAST(NULL AS json)
< ok - -
> exec
| SELECT CAST(NULL AS variant)
< ok - -
> exec
| SELECT CAST(NULL AS date)
< ok - -
> exec
| SELECT CAST(NULL AS uuid)
< ok - -
> exec
| SELECT CAST(NULL AS ipaddress)
< ok - -
> exec
| SELECT CAST(NULL AS hyperloglog)
< ok - -
> exec
| SELECT CAST(NULL AS p4hyperloglog)
< ok - -
> exec
| SELECT CAST(NULL AS tdigest)
< ok - -
> exec
| SELECT CAST(NULL AS setdigest)
< ok - -
> exec
| SELECT CAST(NULL AS qdigest(bigint))
< ok - -
> exec
| SELECT CAST(NULL AS qdigest(double))
< ok - -
> exec
| SELECT CAST(NULL AS "varchar"(10))
< ok - -
> exec
| SELECT CAST(NULL AS time)
< ok - -
> exec
| SELECT CAST(NULL AS time(3))
< ok - -
> exec
| SELECT CAST(NULL AS time with time zone)
< ok - -
> exec
| SELECT CAST(NULL AS time(6) with time zone)
< ok - -
> exec
| SELECT CAST(NULL AS time without time zone)
< ok - -
> exec
| SELECT CAST(NULL AS timestamp)
< ok - -
> exec
| SELECT CAST(NULL AS timestamp(6))
< ok - -
> exec
| SELECT CAST(NULL AS timestamp with time zone)
< ok - -
> exec
| SELECT CAST(NULL AS timestamp(9) with time zone)
< ok - -
> exec
| SELECT CAST(NULL AS timestamp without time zone)
< ok - -
> exec
| SELECT CAST(NULL AS interval year)
< error
| NOT_SUPPORTED 13 0:0
| INTERVAL YEAR type not supported
> exec
| SELECT CAST(NULL AS interval month)
< error
| NOT_SUPPORTED 13 0:0
| INTERVAL MONTH type not supported
> exec
| SELECT CAST(NULL AS interval day)
< error
| NOT_SUPPORTED 13 0:0
| INTERVAL DAY type not supported
> exec
| SELECT CAST(NULL AS interval hour)
< error
| NOT_SUPPORTED 13 0:0
| INTERVAL HOUR type not supported
> exec
| SELECT CAST(NULL AS interval minute)
< error
| NOT_SUPPORTED 13 0:0
| INTERVAL MINUTE type not supported
> exec
| SELECT CAST(NULL AS interval second)
< error
| NOT_SUPPORTED 13 0:0
| INTERVAL SECOND type not supported
> exec
| SELECT CAST(NULL AS interval year to month)
< ok - -
> exec
| SELECT CAST(NULL AS interval day to hour)
< error
| NOT_SUPPORTED 13 0:0
| INTERVAL DAY TO HOUR type not supported
> exec
| SELECT CAST(NULL AS interval day to minute)
< error
| NOT_SUPPORTED 13 0:0
| INTERVAL DAY TO MINUTE type not supported
> exec
| SELECT CAST(NULL AS interval day to second)
< ok - -
> exec
| SELECT CAST(NULL AS interval hour to minute)
< error
| NOT_SUPPORTED 13 0:0
| INTERVAL HOUR TO MINUTE type not supported
> exec
| SELECT CAST(NULL AS interval hour to second)
< error
| NOT_SUPPORTED 13 0:0
| INTERVAL HOUR TO SECOND type not supported
> exec
| SELECT CAST(NULL AS interval minute to second)
< error
| NOT_SUPPORTED 13 0:0
| INTERVAL MINUTE TO SECOND type not supported
> exec
| SELECT CAST(NULL AS interval)
< error
| TYPE_MISMATCH 58 1:8
| line 1:8: Unknown type: interval
> exec
| SELECT CAST(NULL AS interval(5))
< error
| TYPE_MISMATCH 58 1:8
| line 1:8: Unknown type: interval(5)
> exec
| SELECT CAST(NULL AS interval array)
< error
| TYPE_MISMATCH 58 1:8
| line 1:8: Unknown type: array(interval)
> exec
| SELECT CAST(NULL AS ROW(x bigint, y double))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(bigint, double))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(bigint))
< ok - -
> exec
| SELECT CAST(NULL AS ROW("x y" bigint))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(a ARRAY(int)))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(a int ARRAY))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(x map(varchar, int)))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(a double precision))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(a row(b int)))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(1))
< error
| TYPE_MISMATCH 58 1:8
| line 1:8: Unknown type: ROW(1)
> exec
| SELECT CAST(NULL AS ROW(a b))
< error
| TYPE_MISMATCH 58 1:8
| line 1:8: Unknown type: ROW(a b)
> exec
| SELECT CAST(NULL AS ROW(a, b))
< error
| TYPE_MISMATCH 58 1:8
| line 1:8: Unknown type: ROW(a, b)
> exec
| SELECT CAST(NULL AS ROW(a bigint, b))
< error
| TYPE_MISMATCH 58 1:8
| line 1:8: Unknown type: ROW(a bigint, b)
> exec
| SELECT CAST(NULL AS row(x double precision))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(1, 2))
< error
| TYPE_MISMATCH 58 1:8
| line 1:8: Unknown type: ROW(1, 2)
> exec
| SELECT CAST(NULL AS ROW(bigint, 1))
< error
| TYPE_MISMATCH 58 1:8
| line 1:8: Unknown type: ROW(bigint, 1)
> exec
| SELECT CAST(NULL AS ROW(varchar(10), 5))
< error
| TYPE_MISMATCH 58 1:8
| line 1:8: Unknown type: ROW(varchar(10), 5)
> exec
| SELECT CAST(NULL AS ROW(5, 10, 15))
< error
| TYPE_MISMATCH 58 1:8
| line 1:8: Unknown type: ROW(5, 10, 15)
> exec
| SELECT CAST(NULL AS ROW(timestamp))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(interval day))
< error
| NOT_SUPPORTED 13 0:0
| INTERVAL DAY type not supported
> exec
| SELECT CAST(NULL AS ROW(double precision))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(array<int>))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(map<int, int>))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(time bigint))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(timestamp bigint))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(array bigint))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(a bigint array))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(bigint array))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(a timestamp(3) with time zone))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(a interval day to second))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(TIME WITHOUT TIME ZONE))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(TIMESTAMP WITH TIME ZONE))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(INTERVAL DAY TO SECOND))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(BIGINT ARRAY[5]))
< error
| GENERIC_INTERNAL_ERROR 65536 0:0
| Explicit array size not supported
> exec
| SELECT CAST(NULL AS ROW(DOUBLE PRECISION))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(a ARRAY<int>))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(a ARRAY(int)))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(x TIME WITHOUT TIME ZONE, y INTERVAL DAY TO SECOND))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(a int, b varchar(10), c TIMESTAMP WITH TIME ZONE, d INTERVAL DAY TO SECOND, e ARRAY<int>))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(a ARRAY<int>, b ARRAY(int), c bigint ARRAY[5]))
< error
| GENERIC_INTERNAL_ERROR 65536 0:0
| Explicit array size not supported
> exec
| SELECT CAST(NULL AS ARRAY<bigint>)
< ok - -
> exec
| SELECT CAST(NULL AS ARRAY(bigint))
< ok - -
> exec
| SELECT CAST(NULL AS ARRAY(ARRAY(bigint)))
< ok - -
> exec
| SELECT CAST(NULL AS ARRAY<ARRAY<bigint>>)
< ok - -
> exec
| SELECT CAST(NULL AS MAP<varchar, integer>)
< ok - -
> exec
| SELECT CAST(NULL AS MAP(varchar, integer))
< ok - -
> exec
| SELECT CAST(NULL AS MAP(varchar, ARRAY(bigint)))
< ok - -
> exec
| SELECT CAST(NULL AS MAP<varchar, ARRAY<bigint>>)
< ok - -
> exec
| SELECT CAST(NULL AS ARRAY(ROW(a bigint, b varchar)))
< ok - -
> exec
| SELECT CAST(NULL AS bigint ARRAY)
< ok - -
> exec
| SELECT CAST(NULL AS bigint ARRAY ARRAY)
< ok - -
> exec
| SELECT CAST(NULL AS bigint ARRAY[5])
< error
| GENERIC_INTERNAL_ERROR 65536 0:0
| Explicit array size not supported
> exec
| SELECT CAST(NULL AS decimal(10, 2) ARRAY)
< ok - -
> exec
| SELECT CAST(NULL AS INTERVAL DAY TO SECOND ARRAY)
< ok - -
> exec
| SELECT CAST(NULL AS TIMESTAMP(3) WITH TIME ZONE ARRAY)
< ok - -
> exec
| SELECT CAST(NULL AS MAP(varchar, ROW(a bigint, b ARRAY(double))))
< ok - -
> exec
| SELECT CAST(NULL AS ARRAY(MAP(varchar, bigint)))
< ok - -
> exec
| SELECT CAST(NULL AS varchar())
< error
| SYNTAX_ERROR 1 1:29
| line 1:29: mismatched input ')'. Expecting: <integer>, <type>
> exec
| SELECT CAST(NULL AS foo())
< error
| SYNTAX_ERROR 1 1:25
| line 1:25: mismatched input ')'. Expecting: <integer>, <type>
> exec
| SELECT CAST(NULL AS bigint())
< error
| SYNTAX_ERROR 1 1:28
| line 1:28: mismatched input ')'. Expecting: <integer>, <type>
> exec
| SELECT CAST(NULL AS ROW())
< error
| SYNTAX_ERROR 1 1:25
| line 1:25: mismatched input ')'. Expecting: <integer>, <type>
> exec
| SELECT CAST(NULL AS ROW(a bigint,))
< error
| SYNTAX_ERROR 1 1:34
| line 1:34: mismatched input ')'. Expecting: <identifier>, <type>
> exec
| SELECT CAST(NULL AS ROW(a b c))
< error
| SYNTAX_ERROR 1 1:29
| line 1:29: mismatched input 'c'. Expecting: '(', ')', ',', 'ARRAY'
> exec
| SELECT CAST(NULL AS ROW(a varchar(10) extra))
< error
| SYNTAX_ERROR 1 1:39
| line 1:39: mismatched input 'extra'. Expecting: ')', ',', 'ARRAY'
> exec
| SELECT CAST(NULL AS ROW(a bigint, 1))
< error
| SYNTAX_ERROR 1 1:35
| line 1:35: mismatched input '1'. Expecting: <identifier>, <type>
> exec
| SELECT CAST(NULL AS ROW(1, a bigint))
< error
| SYNTAX_ERROR 1 1:30
| line 1:30: mismatched input 'bigint'. Expecting: '(', ')', ',', 'ARRAY'
> exec
| SELECT CAST(NULL AS ROW(x bigint, y 5))
< error
| SYNTAX_ERROR 1 1:37
| line 1:37: mismatched input '5'. Expecting: '(', ')', ',', 'ARRAY', <type>
> exec
| SELECT CAST(NULL AS ROW(varchar(10) x, 5))
< error
| SYNTAX_ERROR 1 1:37
| line 1:37: mismatched input 'x'. Expecting: ')', ',', 'ARRAY'
> exec
| SELECT CAST(NULL AS interval month to year)
< error
| SYNTAX_ERROR 1 1:36
| line 1:36: mismatched input 'to'. Expecting: '(', ')', 'ARRAY'
> exec
| SELECT CAST(NULL AS interval year to day)
< error
| SYNTAX_ERROR 1 1:38
| line 1:38: mismatched input 'day'. Expecting: 'MONTH'
> exec
| SELECT CAST(NULL AS interval day to year)
< error
| SYNTAX_ERROR 1 1:37
| line 1:37: mismatched input 'year'. Expecting: 'HOUR', 'MINUTE', 'SECOND'
> exec
| SELECT CAST(NULL AS interval second to year)
< error
| SYNTAX_ERROR 1 1:37
| line 1:37: mismatched input 'to'. Expecting: '(', ')', 'ARRAY'
> exec
| SELECT CAST(NULL AS interval second to minute)
< error
| SYNTAX_ERROR 1 1:37
| line 1:37: mismatched input 'to'. Expecting: '(', ')', 'ARRAY'
> exec
| SELECT CAST(NULL AS interval to second)
< error
| SYNTAX_ERROR 1 1:30
| line 1:30: mismatched input 'to'. Expecting: '(', ')', 'ARRAY', 'DAY', 'HOUR', 'MINUTE', 'MONTH', 'SECOND', 'YEAR'
> exec
| SELECT CAST(NULL AS MAP<varchar>)
< error
| SYNTAX_ERROR 1 1:32
| line 1:32: mismatched input '>'. Expecting: '(', ',', 'ARRAY'
> exec
| SELECT CAST(NULL AS MAP<varchar,>)
< error
| SYNTAX_ERROR 1 1:33
| line 1:33: mismatched input '>'. Expecting: <type>
> exec
| SELECT CAST(NULL AS ARRAY<>)
< error
| SYNTAX_ERROR 1 1:26
| line 1:26: mismatched input '<>'. Expecting: '(', ')', '<', 'ARRAY'
> exec
| SELECT CAST(NULL AS ROW<bigint>)
< error
| SYNTAX_ERROR 1 1:24
| line 1:24: mismatched input '<'. Expecting: '(', ')', 'ARRAY'
> exec
| SELECT CAST(NULL AS timestamp with zone)
< error
| SYNTAX_ERROR 1 1:36
| line 1:36: mismatched input 'zone'. Expecting: 'TIME'
> exec
| SELECT CAST(NULL AS timestamp without zone)
< error
| SYNTAX_ERROR 1 1:39
| line 1:39: mismatched input 'zone'. Expecting: 'TIME'
> exec
| SELECT CAST(NULL AS double precision precision)
< error
| SYNTAX_ERROR 1 1:38
| line 1:38: mismatched input 'precision'. Expecting: ')', 'ARRAY'
> exec
| SELECT CAST(NULL AS bigint ARRAY[])
< error
| SYNTAX_ERROR 1 1:34
| line 1:34: mismatched input ']'. Expecting: <integer>
> exec
| SELECT CAST(NULL AS bigint ARRAY[abc])
< error
| SYNTAX_ERROR 1 1:34
| line 1:34: mismatched input 'abc'. Expecting: <integer>
> exec
| SELECT CAST(NULL AS bigint ARRAY[5)
< error
| SYNTAX_ERROR 1 1:35
| line 1:35: mismatched input ')'. Expecting: ']'
> exec
| SELECT CAST(NULL AS bigint ARRAY[5 6])
< error
| SYNTAX_ERROR 1 1:36
| line 1:36: mismatched input '6'. Expecting: ']'
