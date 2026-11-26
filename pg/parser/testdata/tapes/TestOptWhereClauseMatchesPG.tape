-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestOptWhereClauseMatchesPG default#1
> exec
| SELECT 1/0;
| ALTER TABLE t ADD EXCLUDE (a WITH =) WHERE (a > 0)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| ALTER TABLE t ADD EXCLUDE (a WITH =) WHERE a > 0
< error
| 42601 syntax error at or near "a"
| position: 56
> exec
| SELECT 1/0;
| ALTER TABLE t ADD EXCLUDE (a WITH =) WHERE (a) OR b
< error
| 42601 syntax error at or near "OR"
| position: 60
> exec
| SELECT 1/0;
| CREATE TABLE t (a int, EXCLUDE USING gist (a WITH &&) WHERE (a IS NOT NULL) DEFERRABLE)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE TABLE t (a int, EXCLUDE USING gist (a WITH &&) WHERE a IS NOT NULL)
< error
| 42601 syntax error at or near "a"
| position: 73
> exec
| SELECT 1/0;
| ALTER PUBLICATION p ADD TABLE t WHERE (a > 0)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| ALTER PUBLICATION p ADD TABLE t WHERE a > 0
< error
| 42601 syntax error at or near "a"
| position: 51
