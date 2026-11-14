-- pin: default trinodb/trino:482 sha256:90b35b7c603eaa1f889bf03981a62b75f998ee6c0f851d9f4e341b49a57022b6 482

== TestEmptyColumnTableRejectedByTrino default#1
> exec
| CREATE TABLE IF NOT EXISTS "default"."empty_table" (
|
| )
< error
| SYNTAX_ERROR 1 3:1
| line 3:1: mismatched input ')'. Expecting: 'LIKE', <identifier>
