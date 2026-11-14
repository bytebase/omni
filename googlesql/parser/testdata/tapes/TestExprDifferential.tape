-- pin: default gcr.io/cloud-spanner-emulator/emulator:1.5.54 sha256:caf1bd24c081e005837b5977bae5a250e25cb4da9f25ec1abc91936ad67e4de2 v1.5.54

== TestExprDifferential default#1
> query
| SELECT 1
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT 1.5
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT NULL
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT TRUE
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT FALSE
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT 'hello'
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT 'a' 'b'
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Concatenation of subsequent string literals is not supported. Did you mean to use the || operator? [at 1:12] SELECT 'a' 'b'            ^"}
> query
| SELECT b'bytes'
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT 0x1F
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT DATE '2020-01-01'
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT TIMESTAMP '2020-01-01 00:00:00'
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT NUMERIC '1.5'
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT JSON '{}'
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT @p
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Incomplete query parameters p"}
> query
| SELECT @@sysvar
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Unrecognized name: @@sysvar [at 1:8] SELECT @@sysvar        ^"}
> query
| SELECT a FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:15] SELECT a FROM t               ^"}
> query
| SELECT t.a.b FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:19] SELECT t.a.b FROM t                   ^"}
> query
| SELECT 1 + 2 * 3
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT (1 + 2) * 3
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT 1 - 2 - 3
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT - - 1
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT ~ 5
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT 5 & 3 | 1
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT 1 << 2 + 3
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT 'a' || 'b' || 'c'
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT TRUE AND FALSE OR TRUE
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT NOT TRUE
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT NOT NOT TRUE
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT NOT 1 = 1
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT 1 = 2
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT 1 != 2
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT 1 <> 2
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT 1 < 2
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT 1 >= 2
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT 1 + 2 = 3
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT a IS NULL FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:23] SELECT a IS NULL FROM t                       ^"}
> query
| SELECT a IS NOT NULL FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:27] SELECT a IS NOT NULL FROM t                           ^"}
> query
| SELECT 1 IS TRUE
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"No matching signature for operator IS TRUE   Argument types: INT64   Signature: (BOOL) IS TRUE     Argument 1: Unable to coerce type INT64 to expected type BOOL [at 1:8] SELECT 1 IS TRUE        ^"}
> query
| SELECT 1 IS NOT FALSE
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"No matching signature for operator IS FALSE   Argument types: INT64   Signature: (BOOL) IS FALSE     Argument 1: Unable to coerce type INT64 to expected type BOOL [at 1:8] SELECT 1 IS NOT FALSE       …"}
> query
| SELECT 1 IS UNKNOWN
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Operand of IS UNKNOWN must be coercible to BOOL, but has type INT64 [at 1:8] SELECT 1 IS UNKNOWN        ^"}
> query
| SELECT a IS NOT DISTINCT FROM b FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"IS DISTINCT FROM is not supported [at 1:17] SELECT a IS NOT DISTINCT FROM b FROM t                 ^"}
> query
| SELECT a BETWEEN 1 AND 2 FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:31] SELECT a BETWEEN 1 AND 2 FROM t                               ^"}
> query
| SELECT a NOT BETWEEN 1 AND 2 FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:35] SELECT a NOT BETWEEN 1 AND 2 FROM t                                   ^"}
> query
| SELECT a IN (1, 2, 3) FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:28] SELECT a IN (1, 2, 3) FROM t                            ^"}
> query
| SELECT a IN (1) FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:22] SELECT a IN (1) FROM t                      ^"}
> query
| SELECT a NOT IN (1, 2) FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:29] SELECT a NOT IN (1, 2) FROM t                             ^"}
> query
| SELECT a IN UNNEST([1, 2]) FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:33] SELECT a IN UNNEST([1, 2]) FROM t                                 ^"}
> query
| SELECT a IN (SELECT x FROM t2) FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:37] SELECT a IN (SELECT x FROM t2) FROM t                                     ^"}
> query
| SELECT x LIKE '%a%' FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:26] SELECT x LIKE '%a%' FROM t                          ^"}
> query
| SELECT x NOT LIKE '%a%' FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:30] SELECT x NOT LIKE '%a%' FROM t                              ^"}
> query
| SELECT x LIKE ANY ('a', 'b') FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:35] SELECT x LIKE ANY ('a', 'b') FROM t                                   ^"}
> query
| SELECT * FROM t WHERE x = ANY (SELECT v FROM s)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:15] SELECT * FROM t WHERE x = ANY (SELECT v FROM s)               ^"}
> query
| SELECT * FROM t WHERE x > ALL (SELECT v FROM s)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:15] SELECT * FROM t WHERE x \u003e ALL (SELECT v FROM s)               ^"}
> query
| SELECT * FROM t WHERE x < SOME (SELECT v FROM s)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:15] SELECT * FROM t WHERE x \u003c SOME (SELECT v FROM s)               ^"}
> query
| SELECT * FROM t WHERE x >= ANY (SELECT v FROM s)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:15] SELECT * FROM t WHERE x \u003e= ANY (SELECT v FROM s)               ^"}
> query
| SELECT * FROM t WHERE x <= ALL (SELECT v FROM s)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:15] SELECT * FROM t WHERE x \u003c= ALL (SELECT v FROM s)               ^"}
> query
| SELECT * FROM t WHERE x != ANY (SELECT v FROM s)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:15] SELECT * FROM t WHERE x != ANY (SELECT v FROM s)               ^"}
> query
| SELECT * FROM t WHERE x <> ALL (SELECT v FROM s)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:15] SELECT * FROM t WHERE x \u003c\u003e ALL (SELECT v FROM s)               ^"}
> query
| SELECT * FROM t WHERE x = ANY (1, 2, 3)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:15] SELECT * FROM t WHERE x = ANY (1, 2, 3)               ^"}
> query
| SELECT * FROM t WHERE x > ALL (1, 2)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:15] SELECT * FROM t WHERE x \u003e ALL (1, 2)               ^"}
> query
| SELECT * FROM t WHERE x = ANY UNNEST([1, 2, 3])
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:15] SELECT * FROM t WHERE x = ANY UNNEST([1, 2, 3])               ^"}
> query
| SELECT * FROM t WHERE x = ANY @{a=1} (SELECT v FROM s)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:15] SELECT * FROM t WHERE x = ANY @{a=1} (SELECT v FROM s)               ^"}
> query
| SELECT CASE WHEN a THEN 1 ELSE 2 END FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:43] SELECT CASE WHEN a THEN 1 ELSE 2 END FROM t                                           ^"}
> query
| SELECT CASE a WHEN 1 THEN 2 END FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:38] SELECT CASE a WHEN 1 THEN 2 END FROM t                                      ^"}
> query
| SELECT CASE WHEN 1 THEN 2 WHEN 3 THEN 4 ELSE 5 END
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"No matching signature for operator CASE; all WHEN argument types must be BOOL but found: INT64; actual argument types (WHEN THEN) ELSE: (INT64 INT64) (INT64 INT64) INT64 [at 1:8] SELECT CASE WHEN 1 TH…"}
> query
| SELECT CAST(x AS INT64) FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:30] SELECT CAST(x AS INT64) FROM t                              ^"}
> query
| SELECT SAFE_CAST(x AS STRING) FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:36] SELECT SAFE_CAST(x AS STRING) FROM t                                    ^"}
> query
| SELECT CAST(x AS STRING FORMAT 'fmt') FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:44] SELECT CAST(x AS STRING FORMAT 'fmt') FROM t                                            ^"}
> query
| SELECT CAST(x AS STRING FORMAT 'f' AT TIME ZONE 'UTC') FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:61] SELECT CAST(x AS STRING FORMAT 'f' AT TIME ZONE 'UTC') FROM t                                                             ^"}
> query
| SELECT EXTRACT(YEAR FROM d) FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:34] SELECT EXTRACT(YEAR FROM d) FROM t                                  ^"}
> query
| SELECT EXTRACT(HOUR FROM ts AT TIME ZONE 'UTC') FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:54] SELECT EXTRACT(HOUR FROM ts AT TIME ZONE 'UTC') FROM t                                                      ^"}
> query
| SELECT f(x, y) FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:21] SELECT f(x, y) FROM t                     ^"}
> query
| SELECT mypkg.myfunc(x) FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:29] SELECT mypkg.myfunc(x) FROM t                             ^"}
> query
| SELECT COUNT(*) FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:22] SELECT COUNT(*) FROM t                      ^"}
> query
| SELECT COUNT(DISTINCT x) FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:31] SELECT COUNT(DISTINCT x) FROM t                               ^"}
> query
| SELECT IF(a, 1, 2) FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:25] SELECT IF(a, 1, 2) FROM t                         ^"}
> query
| SELECT GROUPING(x) FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:25] SELECT GROUPING(x) FROM t                         ^"}
> query
| SELECT ARRAY_AGG(x IGNORE NULLS) FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:39] SELECT ARRAY_AGG(x IGNORE NULLS) FROM t                                       ^"}
> query
| SELECT ARRAY_AGG(x ORDER BY y) FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:37] SELECT ARRAY_AGG(x ORDER BY y) FROM t                                     ^"}
> query
| SELECT ARRAY_AGG(x LIMIT 10) FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:35] SELECT ARRAY_AGG(x LIMIT 10) FROM t                                   ^"}
> query
| SELECT func(a => 1, b => 2)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Function not found: func [at 1:8] SELECT func(a =\u003e 1, b =\u003e 2)        ^"}
> query
| SELECT ARRAY_TRANSFORM([1], e -> e + 1)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT REDUCE([1], (a, b) -> a + b)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Function not found: REDUCE [at 1:8] SELECT REDUCE([1], (a, b) -\u003e a + b)        ^"}
> query
| SELECT SUM(x) OVER (PARTITION BY a ORDER BY b) FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:53] SELECT SUM(x) OVER (PARTITION BY a ORDER BY b) FROM t                                                     ^"}
> query
| SELECT SUM(x) OVER w FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:27] SELECT SUM(x) OVER w FROM t                           ^"}
> query
| SELECT ROW_NUMBER() OVER (ORDER BY x ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:93] ...BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) FROM t                                                      ^"}
> query
| SELECT COUNT(*) OVER () FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:30] SELECT COUNT(*) OVER () FROM t                              ^"}
> query
| SELECT [1, 2, 3]
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT ARRAY[1, 2]
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT ARRAY<INT64>[1, 2]
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT []
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT ARRAY<INT64>[]
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT STRUCT(1 AS x, 2 AS y).x
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT STRUCT<x INT64, y STRING>(1, 'a').x
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT STRUCT().field1
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Field name field1 does not exist in STRUCT\u003c\u003e [at 1:17] SELECT STRUCT().field1                 ^"}
> query
| SELECT (1, 2, 3).field1
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Field name field1 does not exist in STRUCT\u003cINT64, INT64, INT64\u003e [at 1:18] SELECT (1, 2, 3).field1                  ^"}
> query
| SELECT a[0] FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:18] SELECT a[0] FROM t                  ^"}
> query
| SELECT a[OFFSET(0)] FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:26] SELECT a[OFFSET(0)] FROM t                          ^"}
> query
| SELECT a[ORDINAL(1)] FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:27] SELECT a[ORDINAL(1)] FROM t                           ^"}
> query
| SELECT a[SAFE_OFFSET(0)] FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:31] SELECT a[SAFE_OFFSET(0)] FROM t                               ^"}
> query
| SELECT f(x).y FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:20] SELECT f(x).y FROM t                    ^"}
> query
| SELECT a[0].b FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:20] SELECT a[0].b FROM t                    ^"}
> query
| SELECT INTERVAL 5 DAY
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT INTERVAL x YEAR TO MONTH FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:38] SELECT INTERVAL x YEAR TO MONTH FROM t                                      ^"}
> query
| SELECT EXISTS(SELECT 1)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT ARRAY(SELECT 1)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT (SELECT 1) + 2
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT 1 +
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Unexpected end of statement [at 1:11] SELECT 1 +           ^"}
> query
| SELECT 1 2
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Expected end of input but got integer literal \"2\" [at 1:10] SELECT 1 2          ^"}
> query
| SELECT a BETWEEN 1 FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Unexpected keyword FROM [at 1:20] SELECT a BETWEEN 1 FROM t                    ^"}
> query
| SELECT a BETWEEN 1 OR 2 FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Expression in BETWEEN must be parenthesized [at 1:18] SELECT a BETWEEN 1 OR 2 FROM t                  ^"}
> query
| SELECT a = b = c FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Unexpected \"=\" [at 1:14] SELECT a = b = c FROM t              ^"}
> query
| SELECT 1 < 2 < 3
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Unexpected \"\u003c\" [at 1:14] SELECT 1 \u003c 2 \u003c 3              ^"}
> query
| SELECT a = b IS NULL FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Unexpected keyword IS [at 1:14] SELECT a = b IS NULL FROM t              ^"}
> query
| SELECT a IN (1) IN (2) FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Expression to the left of IN must be parenthesized [at 1:17] SELECT a IN (1) IN (2) FROM t                 ^"}
> query
| SELECT 1 IS NULL IS NULL
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Expression to the left of IS must be parenthesized [at 1:18] SELECT 1 IS NULL IS NULL                  ^"}
> query
| SELECT 1 LIKE 'a' LIKE 'b'
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Unexpected keyword LIKE [at 1:19] SELECT 1 LIKE 'a' LIKE 'b'                   ^"}
> query
| SELECT * FROM t WHERE x = ANY ()
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Unexpected \")\" [at 1:32] SELECT * FROM t WHERE x = ANY ()                                ^"}
> query
| SELECT * FROM t WHERE x = ANY ANY (SELECT v FROM s)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Unexpected keyword ANY [at 1:31] SELECT * FROM t WHERE x = ANY ANY (SELECT v FROM s)                               ^"}
> query
| SELECT * FROM t WHERE x = ANY SELECT v FROM s
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Unexpected keyword SELECT [at 1:31] SELECT * FROM t WHERE x = ANY SELECT v FROM s                               ^"}
> query
| SELECT * FROM t WHERE x IS ANY (SELECT v FROM s)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Expected keyword FALSE or keyword NULL or keyword TRUE or keyword UNKNOWN but got keyword ANY [at 1:28] SELECT * FROM t WHERE x IS ANY (SELECT v FROM s)                            ^"}
> query
| SELECT * FROM t WHERE x = ANY (SELECT v FROM s) = y
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Expression to the left of comparison must be parenthesized [at 1:49] SELECT * FROM t WHERE x = ANY (SELECT v FROM s) = y                                                 ^"}
> query
| SELECT 1 BETWEEN 0 AND 2 BETWEEN 0 AND 1
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Unexpected keyword BETWEEN [at 1:26] SELECT 1 BETWEEN 0 AND 2 BETWEEN 0 AND 1                          ^"}
> query
| SELECT CASE END
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Unexpected keyword END [at 1:13] SELECT CASE END             ^"}
> query
| SELECT CAST(x INT64) FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Expected keyword AS but got identifier \"INT64\" [at 1:15] SELECT CAST(x INT64) FROM t               ^"}
> query
| SELECT CAST(x AS) FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Unexpected \")\" [at 1:17] SELECT CAST(x AS) FROM t                 ^"}
> query
| SELECT CAST(CAST AS INT64)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Expected \"(\" but got keyword AS [at 1:18] SELECT CAST(CAST AS INT64)                  ^"}
> query
| SELECT EXTRACT(x y) FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Unexpected identifier \"y\" [at 1:18] SELECT EXTRACT(x y) FROM t                  ^"}
> query
| SELECT a IN () FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Unexpected \")\" [at 1:14] SELECT a IN () FROM t              ^"}
> query
| SELECT a..b FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Unexpected \".\" [at 1:10] SELECT a..b FROM t          ^"}
> query
| SELECT * 5
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Expected end of input but got integer literal \"5\" [at 1:10] SELECT * 5          ^"}
> query
| SELECT EXTRACT(HOUR FROM ts @ TIME ZONE 'UTC') FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Expected \")\" or keyword AT but got \"@\" [at 1:29] SELECT EXTRACT(HOUR FROM ts @ TIME ZONE 'UTC') FROM t                             ^"}
