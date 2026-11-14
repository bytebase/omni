-- pin: default trinodb/trino:482 sha256:90b35b7c603eaa1f889bf03981a62b75f998ee6c0f851d9f4e341b49a57022b6 482

== TestCompletion_ContextsAreValidSyntax default#1
> exec
| DROP TABLE IF EXISTS memory.default.orders
< ok - -
> exec
| DROP TABLE IF EXISTS memory.default.customer
< ok - -
> exec
| CREATE TABLE memory.default.orders (orderkey bigint, custkey bigint, totalprice double)
< ok - -
> exec
| CREATE TABLE memory.default.customer (custkey bigint, name varchar)
< ok - -
> exec
| SELECT * FROM orders
< ok - -
> exec
| SELECT * FROM orders o JOIN customer ON o.custkey = customer.custkey
< ok - -
> exec
| SELECT orderkey FROM orders
< ok - -
> exec
| SELECT * FROM customer WHERE custkey
< error
| TYPE_MISMATCH 58 1:30
| line 1:30: WHERE clause must evaluate to a boolean: actual type bigint
> exec
| SELECT * FROM orders o WHERE o.totalprice
< error
| TYPE_MISMATCH 58 1:30
| line 1:30: WHERE clause must evaluate to a boolean: actual type double
> exec
| SELECT * FROM memory.default
< error
| SCHEMA_NOT_FOUND 45 1:15
| line 1:15: Schema 'memory' does not exist
> exec
| SELECT 1
< ok - -
