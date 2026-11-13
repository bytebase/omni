-- pin: default trinodb/trino:482 sha256:90b35b7c603eaa1f889bf03981a62b75f998ee6c0f851d9f4e341b49a57022b6 482

== TestDataType_RoundTripOracle default#1
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
| SELECT CAST(NULL AS DOUBLE PRECISION)
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
| SELECT CAST(NULL AS TIME)
< ok - -
> exec
| SELECT CAST(NULL AS TIME(3))
< ok - -
> exec
| SELECT CAST(NULL AS TIME WITH TIME ZONE)
< ok - -
> exec
| SELECT CAST(NULL AS TIME(6) WITH TIME ZONE)
< ok - -
> exec
| SELECT CAST(NULL AS TIME WITHOUT TIME ZONE)
< ok - -
> exec
| SELECT CAST(NULL AS TIMESTAMP)
< ok - -
> exec
| SELECT CAST(NULL AS TIMESTAMP(6))
< ok - -
> exec
| SELECT CAST(NULL AS TIMESTAMP WITH TIME ZONE)
< ok - -
> exec
| SELECT CAST(NULL AS TIMESTAMP(9) WITH TIME ZONE)
< ok - -
> exec
| SELECT CAST(NULL AS TIMESTAMP WITHOUT TIME ZONE)
< ok - -
> exec
| SELECT CAST(NULL AS INTERVAL YEAR)
< error
| NOT_SUPPORTED 13 0:0
| INTERVAL YEAR type not supported
> exec
| SELECT CAST(NULL AS INTERVAL MONTH)
< error
| NOT_SUPPORTED 13 0:0
| INTERVAL MONTH type not supported
> exec
| SELECT CAST(NULL AS INTERVAL DAY)
< error
| NOT_SUPPORTED 13 0:0
| INTERVAL DAY type not supported
> exec
| SELECT CAST(NULL AS INTERVAL HOUR)
< error
| NOT_SUPPORTED 13 0:0
| INTERVAL HOUR type not supported
> exec
| SELECT CAST(NULL AS INTERVAL MINUTE)
< error
| NOT_SUPPORTED 13 0:0
| INTERVAL MINUTE type not supported
> exec
| SELECT CAST(NULL AS INTERVAL SECOND)
< error
| NOT_SUPPORTED 13 0:0
| INTERVAL SECOND type not supported
> exec
| SELECT CAST(NULL AS INTERVAL YEAR TO MONTH)
< ok - -
> exec
| SELECT CAST(NULL AS INTERVAL DAY TO HOUR)
< error
| NOT_SUPPORTED 13 0:0
| INTERVAL DAY TO HOUR type not supported
> exec
| SELECT CAST(NULL AS INTERVAL DAY TO MINUTE)
< error
| NOT_SUPPORTED 13 0:0
| INTERVAL DAY TO MINUTE type not supported
> exec
| SELECT CAST(NULL AS INTERVAL DAY TO SECOND)
< ok - -
> exec
| SELECT CAST(NULL AS INTERVAL HOUR TO MINUTE)
< error
| NOT_SUPPORTED 13 0:0
| INTERVAL HOUR TO MINUTE type not supported
> exec
| SELECT CAST(NULL AS INTERVAL HOUR TO SECOND)
< error
| NOT_SUPPORTED 13 0:0
| INTERVAL HOUR TO SECOND type not supported
> exec
| SELECT CAST(NULL AS INTERVAL MINUTE TO SECOND)
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
| SELECT CAST(NULL AS ARRAY<interval>)
< error
| TYPE_MISMATCH 58 1:8
| line 1:8: Unknown type: ARRAY(interval)
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
| SELECT CAST(NULL AS ROW(a ARRAY<int>))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(x map(varchar, int)))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(a DOUBLE PRECISION))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(a ROW(b int)))
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
| SELECT CAST(NULL AS ROW(x DOUBLE PRECISION))
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
| SELECT CAST(NULL AS ROW(TIMESTAMP))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(INTERVAL DAY))
< error
| NOT_SUPPORTED 13 0:0
| INTERVAL DAY type not supported
> exec
| SELECT CAST(NULL AS ROW(DOUBLE PRECISION))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(ARRAY<int>))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(MAP<int, int>))
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
| SELECT CAST(NULL AS ROW(a ARRAY<bigint>))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(ARRAY<bigint>))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(a TIMESTAMP(3) WITH TIME ZONE))
< ok - -
> exec
| SELECT CAST(NULL AS ROW(a INTERVAL DAY TO SECOND))
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
| SELECT CAST(NULL AS ROW(ARRAY<BIGINT>))
< ok - -
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
| SELECT CAST(NULL AS ROW(a ARRAY<int>, b ARRAY(int), c ARRAY<bigint>))
< ok - -
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
| SELECT CAST(NULL AS ARRAY<bigint>)
< ok - -
> exec
| SELECT CAST(NULL AS ARRAY<ARRAY<bigint>>)
< ok - -
> exec
| SELECT CAST(NULL AS ARRAY<bigint>)
< ok - -
> exec
| SELECT CAST(NULL AS ARRAY<decimal(10, 2)>)
< ok - -
> exec
| SELECT CAST(NULL AS ARRAY<INTERVAL DAY TO SECOND>)
< ok - -
> exec
| SELECT CAST(NULL AS ARRAY<TIMESTAMP(3) WITH TIME ZONE>)
< ok - -
> exec
| SELECT CAST(NULL AS MAP(varchar, ROW(a bigint, b ARRAY(double))))
< ok - -
> exec
| SELECT CAST(NULL AS ARRAY(MAP(varchar, bigint)))
< ok - -
