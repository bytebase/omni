-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestS3KnownBetterWhitelist default#1
> exec
| CREATE FUNCTION f6() RETURNS int BEGIN ATOMIC SELECT t4.end FROM t4; END;
< error
| 42P01 relation "t4" does not exist
| position: 66
> exec
| CREATE FUNCTION f7() RETURNS int BEGIN ATOMIC SELECT 1 AS end; END;
< ok 0 -
