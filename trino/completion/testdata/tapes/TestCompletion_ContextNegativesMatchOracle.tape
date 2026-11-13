-- pin: default trinodb/trino:482 sha256:90b35b7c603eaa1f889bf03981a62b75f998ee6c0f851d9f4e341b49a57022b6 482

== TestCompletion_ContextNegativesMatchOracle default#1
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
| SELECT * FROM orders o JOIN customer c USING custkey
< error
| SYNTAX_ERROR 1 1:46
| line 1:46: mismatched input 'custkey'. Expecting: '('
> exec
| SELECT * FROM orders o JOIN customer c USING (custkey)
< ok - -
> exec
| UPDATE orders SET CASE
< error
| SYNTAX_ERROR 1 1:19
| line 1:19: mismatched input 'CASE'. Expecting: <identifier>
> exec
| UPDATE orders o SET totalprice = 1
< error
| SYNTAX_ERROR 1 1:15
| line 1:15: mismatched input 'o'. Expecting: '@', 'SET'
> exec
| DELETE FROM orders o WHERE orderkey = 1
< error
| SYNTAX_ERROR 1 1:20
| line 1:20: mismatched input 'o'. Expecting: '.', '@', 'WHERE', <EOF>
> exec
| MERGE INTO orders o USING customer c ON o.custkey = c.custkey WHEN MATCHED THEN UPDATE SET totalprice = c.custkey
< error
| NOT_SUPPORTED 13 0:0
| This connector does not support modifying table rows
