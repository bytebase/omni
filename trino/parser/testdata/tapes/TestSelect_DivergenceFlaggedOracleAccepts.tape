-- pin: default trinodb/trino:482 sha256:90b35b7c603eaa1f889bf03981a62b75f998ee6c0f851d9f4e341b49a57022b6 482

== TestSelect_DivergenceFlaggedOracleAccepts default#1
> exec
| SELECT * FROM JSON_TABLE('[1]', 'lax $' COLUMNS (a integer))
< ok - -
> exec
| WITH SESSION foo = 1 SELECT 1
< error
| INVALID_SESSION_PROPERTY 14 1:20
| line 1:20: Session property 'foo' does not exist
> exec
| SELECT 1 UNION CORRESPONDING SELECT 1
< error
| MISSING_COLUMN_NAME 55 1:10
| line 1:10: Anonymous columns are not allowed in set operations with CORRESPONDING
