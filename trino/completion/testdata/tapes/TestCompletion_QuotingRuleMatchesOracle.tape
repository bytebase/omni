-- pin: default trinodb/trino:482 sha256:90b35b7c603eaa1f889bf03981a62b75f998ee6c0f851d9f4e341b49a57022b6 482

== TestCompletion_QuotingRuleMatchesOracle default#1
> exec
| SELECT 1 AS select
< error
| SYNTAX_ERROR 1 1:13
| line 1:13: mismatched input 'select'. Expecting: <identifier>
> exec
| SELECT 1 AS "select"
< ok - -
> exec
| SELECT TABLE_NAME FROM information_schema.columns
< ok - -
> exec
| SELECT "TABLE_NAME" FROM information_schema.columns
< ok - -
