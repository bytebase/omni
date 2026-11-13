-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestRoutineBodyStatementsMatchPG default#1
> exec
| SELECT 1/0;
| CREATE FUNCTION f() RETURNS int LANGUAGE sql BEGIN ATOMIC RETURN 1; END
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE FUNCTION f() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 1; RETURN 2; END
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE FUNCTION f() LANGUAGE sql BEGIN ATOMIC END
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE FUNCTION f() LANGUAGE sql BEGIN ATOMIC ; END
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE FUNCTION f() LANGUAGE sql BEGIN ATOMIC ; ; SELECT 1; END
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE FUNCTION f() RETURNS int LANGUAGE sql BEGIN ATOMIC RETURN 1 END
< error
| 42601 syntax error at or near "END"
| position: 80
> exec
| SELECT 1/0;
| CREATE FUNCTION f() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT 1; RETURN 2 END
< error
| 42601 syntax error at or near "END"
| position: 90
> exec
| SELECT 1/0;
| CREATE FUNCTION a() BEGIN ATOMIC; zzz END
< error
| 42601 syntax error at or near "zzz"
| position: 47
> exec
| SELECT 1/0;
| CREATE FUNCTION a() BEGIN ATOMIC zzz END
< error
| 42601 syntax error at or near "zzz"
| position: 46
> exec
| SELECT 1/0;
| CREATE FUNCTION a() BEGIN ATOMIC zzz; END
< error
| 42601 syntax error at or near "zzz"
| position: 46
> exec
| SELECT 1/0;
| CREATE FUNCTION a() BEGIN ATOMIC a; END
< error
| 42601 syntax error at or near "a"
| position: 46
> exec
| SELECT 1/0;
| CREATE FUNCTION a() BEGIN ATOMIC RETURN a zzz; END
< error
| 42601 syntax error at or near "zzz"
| position: 55
> exec
| SELECT 1/0;
| CREATE FUNCTION a() BEGIN ATOMIC RETURNING a; END
< error
| 42601 syntax error at or near "RETURNING"
| position: 46
