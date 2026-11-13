-- pin: default trinodb/trino:482 sha256:90b35b7c603eaa1f889bf03981a62b75f998ee6c0f851d9f4e341b49a57022b6 482

== TestFoundation_SplitRoundTripOracle default#1
> exec
| SELECT 1
< ok - -
> exec
|  SELECT 2
< ok - -
> exec
| SELECT 1
< ok - -
> exec
| SELECT 'a;b' AS c
< ok - -
> exec
|  SELECT 2
< ok - -
> exec
| SELECT 1 AS "a;b"
< ok - -
> exec
|  SELECT 2
< ok - -
> exec
| SELECT 1 -- a;b
|
< ok - -
> exec
|  SELECT 2
< ok - -
> exec
| SELECT 1 /* a;b */
< ok - -
> exec
|  SELECT 2
< ok - -
> exec
| CREATE SCHEMA s
< ok - -
> exec
|  DROP SCHEMA s
< ok - -
> exec
| CREATE FUNCTION test.default.f(x integer) RETURNS integer BEGIN RETURN x * 2; END
< error
| CATALOG_NOT_FOUND 44 0:0
| Catalog 'test' not found
> exec
| CREATE FUNCTION test.default.g(x integer) RETURNS integer BEGIN IF x > 0 THEN RETURN 1; END IF; RETURN 0; END
< error
| CATALOG_NOT_FOUND 44 0:0
| Catalog 'test' not found
> exec
| WITH FUNCTION f(x integer) RETURNS integer RETURN x * 2 SELECT f(1)
< ok - -
> exec
|  SELECT 2
< ok - -
> exec
| CREATE FUNCTION test.default.g(x integer) RETURNS integer BEGIN IF x > 0 THEN RETURN 1; END IF; RETURN 0; END
< error
| CATALOG_NOT_FOUND 44 0:0
| Catalog 'test' not found
> exec
|  SELECT 2
< ok - -
