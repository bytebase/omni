-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestMissingNamesMatchPG default#1
> exec
| SELECT 1/0;
| ALTER EXTENSION a ADD OPERATOR CLASS b USING c
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| ALTER EXTENSION a ADD OPERATOR FAMILY b USING c
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| ALTER EXTENSION a ADD OPERATOR CLASS USING a
< error
| 42601 syntax error at or near "USING"
| position: 50
> exec
| SELECT 1/0;
| ALTER EXTENSION a ADD OPERATOR FAMILY USING a
< error
| 42601 syntax error at or near "USING"
| position: 51
> exec
| SELECT 1/0;
| ALTER EXTENSION a ADD OPERATOR CLASS b c
< error
| 42601 syntax error at or near "c"
| position: 52
> exec
| SELECT 1/0;
| ALTER EXTENSION a ADD OPERATOR FAMILY b USING
< error
| 42601 syntax error at end of input
| position: 58
> exec
| SELECT 1/0;
| CREATE FUNCTION a(b c.d%TYPE) RETURNS int LANGUAGE sql RETURN 1
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE FUNCTION a(a.% TYPE)
< error
| 42601 syntax error at or near "%"
| position: 33
> exec
| SELECT 1/0;
| ALTER SYSTEM SET a TO TRUE
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| ALTER SYSTEM SET TO TRUE
< error
| 42601 syntax error at or near "TO"
| position: 30
